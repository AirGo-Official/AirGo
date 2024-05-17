package service

import (
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/ppoonk/AirGo/constant"
	"github.com/ppoonk/AirGo/global"
	"github.com/ppoonk/AirGo/model"
	"github.com/ppoonk/AirGo/utils/encrypt_plugin"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"strconv"
	"strings"
	"time"
)

type AdminNode struct {
}

var AdminNodeSvc *AdminNode

// 新建节点
func (n *AdminNode) NewNode(nodeParams *model.Node) error {
	return global.DB.Transaction(func(tx *gorm.DB) error {
		//检查remarks冲突，避免clash等客户端更新订阅时会报错
		tempNode, err := n.FirstNode(&model.Node{Remarks: nodeParams.Remarks})
		if err == nil { //已存在,重新命名，防止clash客户端等无法解析同名节点
			nodeParams.Remarks = fmt.Sprintf("%s-%s", tempNode.Remarks, encrypt_plugin.RandomString(4))
		}
		switch nodeParams.NodeType {
		case constant.NODE_TYPE_NORMAL:
			if nodeParams.Protocol == constant.NODE_PROTOCOL_SHADOWSOCKS && nodeParams.ServerKey == "" {
				nodeParams.ServerKey = encrypt_plugin.RandomString(32)
			}
			if nodeParams.Protocol == constant.NODE_PROTOCOL_HYSTERIA2 && nodeParams.HyObfs == "Salamander" && nodeParams.HyObfsPassword == "" {
				nodeParams.HyObfsPassword = encrypt_plugin.RandomString(32)
			}

		case constant.NODE_TYPE_TRANSFER:
			//如果该节点是中转节点，则把父节点(对接xrayr等的正常节点)的参数拷贝给该节点，减少更新订阅时查询次数
			//步骤  1、先把父节点基础参数修改为前端传过来的 nodeParams  2、父节点全部参数赋值给 nodeParams 3、由 nodeParams 数据创建节点
			parentNode, err := n.FirstNode(&model.Node{ID: nodeParams.TransferNodeID})
			if err != nil {
				return err
			}
			if parentNode.NodeType != constant.NODE_TYPE_NORMAL {
				return errors.New("Parent node is not 'normal node'")
			}
			parentNode.ID = 0
			parentNode.NodeType = constant.NODE_TYPE_TRANSFER //更换类型为中转
			parentNode.CreatedAt, parentNode.UpdatedAt = time.Now(), time.Now()
			parentNode.Remarks = nodeParams.Remarks
			parentNode.TransferNodeID = nodeParams.TransferNodeID
			parentNode.TransferAddress = nodeParams.TransferAddress
			parentNode.TransferPort = nodeParams.TransferPort
			nodeParams = parentNode
		case constant.NODE_TYPE_SHARED:

		default:
			return errors.New(constant.ERROR_INVALID_NODE_TYPE)
		}

		//矫正一些参数
		nodeParams.ID = 0           //防止前端意外传过来id，导致数据库无法创建
		nodeParams.Enabled = true   //默认启用节点
		nodeParams.NodeOrder = 9999 //默认将排序放到最低下

		//创建
		return tx.Create(&nodeParams).Error
	})

}

// 查询节点
func (n *AdminNode) FirstNode(nodeParams *model.Node) (*model.Node, error) {
	var node model.Node
	err := global.DB.Where(&nodeParams).First(&node).Error
	return &node, err
}

// 查询节点列表
func (n *AdminNode) GetNodeList(params *model.QueryParams) (*model.CommonDataResp, error) {
	var nodeList []model.Node
	var total int64
	_, dataSql := CommonSqlFindSqlHandler(params)
	dataSql = dataSql[strings.Index(dataSql, "WHERE ")+6:] //去掉`WHERE `
	if dataSql == "" {
		dataSql = "id > 0" //当前端什么参数没有传时，默认添加一个参数
	}
	err := global.DB.Model(&model.Node{}).
		Count(&total).Where(dataSql).
		Preload("Access").
		Find(&nodeList).
		Error
	if err != nil {
		return nil, err
	}
	return &model.CommonDataResp{total, nodeList}, err
}

// 更新节点
func (n *AdminNode) UpdateNode(node *model.Node) error {
	return global.DB.Transaction(func(tx *gorm.DB) error {
		//更新节点
		err := tx.Model(&node).Association("Access").Replace(&node.Access)
		if err != nil {
			return err
		}
		//关联的 Access 已在上面更新，另外两个关联不需更新
		err = tx.Omit("Access", "Goods", "TrafficLogs").Save(&node).Error
		if err != nil {
			return err
		}
		if node.NodeType == constant.NODE_TYPE_NORMAL { //该节点是正常节点，则更新该节点节点绑定的中转
			var nodeArr []model.Node
			err = tx.Where(&model.Node{TransferNodeID: node.ID}).Find(&nodeArr).Error
			if err != nil {
				return err
			}
			if len(nodeArr) > 0 {
				for k, _ := range nodeArr { //遍历中转节点
					temp := *node
					temp.NodeType = constant.NODE_TYPE_TRANSFER
					temp.NodeOrder = nodeArr[k].NodeOrder
					temp.ID = nodeArr[k].ID
					temp.CreatedAt, temp.UpdatedAt = nodeArr[k].CreatedAt, nodeArr[k].UpdatedAt
					temp.Remarks = nodeArr[k].Remarks
					temp.TransferNodeID = nodeArr[k].TransferNodeID
					temp.TransferAddress = nodeArr[k].TransferAddress
					temp.TransferPort = nodeArr[k].TransferPort
					nodeArr[k] = temp
				}
				return tx.Save(&nodeArr).Error
			}
		}
		return nil
	})
}

// 删除节点
func (n *AdminNode) DeleteNode(node *model.Node) error {
	return global.DB.Transaction(func(tx *gorm.DB) error {
		err := tx.Where("node_id = ?", node.ID).Delete(&model.NodeTrafficLog{}).Error
		if err != nil {
			return err
		}
		return tx.Select(clause.Associations).Delete(&node).Error
	})
}

// 更新节点流量记录
func (n *AdminNode) UpdateNodeTraffic(trafficLog *model.NodeTrafficLog, AGUserTraffic *model.AGUserTraffic) {
	//查询当天的数据
	now := time.Now()
	zeroTime := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	var nodeTraffic model.NodeTrafficLog
	err := global.DB.Where("node_id = ? AND created_at > ?", AGUserTraffic.ID, zeroTime).Last(&nodeTraffic).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			global.DB.Create(&trafficLog)
		}
	} else {
		nodeTraffic.U += trafficLog.U
		nodeTraffic.D += trafficLog.D
		global.DB.Save(&nodeTraffic)
	}
}

// 清理节点流量记录
func (n *AdminNode) ClearNodeTraffic() error {
	y, m, _ := time.Now().Date()
	startTime := time.Date(y, m-2, 1, 0, 0, 0, 0, time.Local) //清除2个月之前的数据
	return global.DB.Transaction(func(tx *gorm.DB) error {
		return tx.Where("created_at < ?", startTime).Delete(&model.NodeTrafficLog{}).Error
	})
}

// 查询节点列表 with流量
func (n *AdminNode) GetNodeListWithTraffic(params *model.QueryParams) (*model.CommonDataResp, error) {
	var nodeList []model.Node
	var total int64
	var startTime, endTime time.Time
	//时间格式转换
	startTime, err := time.ParseInLocation("2006-01-02 15:04:05", params.FieldParamsList[0].ConditionValue, time.Local)
	if err != nil {
		return nil, err
	}
	endTime, _ = time.ParseInLocation("2006-01-02 15:04:05", params.FieldParamsList[1].ConditionValue, time.Local)
	if err != nil {
		return nil, err
	}
	//注意：	params.FieldParamsList 数组前两项传时间，第三个开始传查询参数
	params.FieldParamsList = append([]model.FieldParamsItem{}, params.FieldParamsList[2:]...)
	_, dataSql := CommonSqlFindSqlHandler(params)
	dataSql = dataSql[strings.Index(dataSql, "WHERE ")+6:] //去掉`WHERE `
	if dataSql == "" {
		dataSql = "id > 0" //当前端什么参数没有传时，默认添加一个参数
	}
	err = global.DB.
		Model(&model.Node{}).
		Count(&total).
		Where(dataSql).
		//Preload("TrafficLogs", global.DB.Where("created_at > ? and created_at < ?", startTime, endTime)).
		Preload("TrafficLogs", "created_at > ? and created_at < ?", startTime, endTime).
		Find(&nodeList).
		Error
	if err != nil {
		return nil, err
	}
	for i1, _ := range nodeList {
		//处理流量记录
		for _, v := range nodeList[i1].TrafficLogs {
			nodeList[i1].TotalUp = nodeList[i1].TotalUp + v.U
			nodeList[i1].TotalDown = nodeList[i1].TotalDown + v.D
		}
		nodeList[i1].TrafficLogs = []model.NodeTrafficLog{} //清空traffic
	}
	//fmt.Println("nodeList:")
	//Show(nodeList)
	return &model.CommonDataResp{total, nodeList}, err
}

// 更新node status
func (n *AdminNode) UpdateNodeStatus(customerServerIDs []int64, trafficLog *model.NodeTrafficLog) {
	var duration float64 = 60 //默认60秒间隔
	cacheStatus, ok := global.LocalCache.Get(fmt.Sprintf("%s%d", constant.CACHE_NODE_STATUS_BY_NODEID, trafficLog.NodeID))
	if ok {
		oldStatus := cacheStatus.(model.NodeStatus)
		oldStatus.Status = true
		oldStatus.UserAmount = int64(len(customerServerIDs))
		now := time.Now()
		duration = now.Sub(oldStatus.LastTime).Seconds()
		oldStatus.D, _ = strconv.ParseFloat(fmt.Sprintf("%.2f", float64(trafficLog.D)/duration), 64) //Byte per second
		oldStatus.U, _ = strconv.ParseFloat(fmt.Sprintf("%.2f", float64(trafficLog.U)/duration), 64)
		oldStatus.LastTime = now
		global.LocalCache.Set(fmt.Sprintf("%s%d",
			constant.CACHE_NODE_STATUS_BY_NODEID, trafficLog.NodeID),
			oldStatus,
			constant.CAHCE_NODE_STATUS_TIMEOUT*time.Minute)
	} else {
		var nodeStatus model.NodeStatus
		nodeStatus.Status = true
		nodeStatus.ID = trafficLog.NodeID
		nodeStatus.UserAmount = int64(len(customerServerIDs))
		nodeStatus.D, _ = strconv.ParseFloat(fmt.Sprintf("%.2f", float64(trafficLog.D)/duration), 64) //Byte per second
		nodeStatus.U, _ = strconv.ParseFloat(fmt.Sprintf("%.2f", float64(trafficLog.U)/duration), 64)
		nodeStatus.LastTime = time.Now()
		global.LocalCache.Set(fmt.Sprintf("%s%d", constant.CACHE_NODE_STATUS_BY_NODEID,
			trafficLog.NodeID),
			nodeStatus,
			constant.CAHCE_NODE_STATUS_TIMEOUT*time.Minute)
	}
}

// 获取 node status
func (n *AdminNode) GetNodesStatus() *[]model.NodeStatus {
	var nodesArr []model.Node
	global.DB.
		Model(&model.Node{}).
		Select("id", "remarks", "traffic_rate").
		Where("enabled = ? AND node_type = ?", true, constant.NODE_TYPE_NORMAL).
		Order("node_order").
		Find(&nodesArr)
	var nodestatusArr []model.NodeStatus
	for _, v := range nodesArr {
		var nodeStatus = model.NodeStatus{}
		vStatus, ok := global.LocalCache.Get(fmt.Sprintf("%s%d", constant.CACHE_NODE_STATUS_BY_NODEID, v.ID))
		if !ok { //cache过期，离线了
			nodeStatus.ID = v.ID
			nodeStatus.Name = v.Remarks
			nodeStatus.TrafficRate = v.TrafficRate
			nodeStatus.Status = false
			nodeStatus.D = 0
			nodeStatus.U = 0
			nodestatusArr = append(nodestatusArr, nodeStatus)
		} else {
			nodeStatus = vStatus.(model.NodeStatus)
			nodeStatus.Name = v.Remarks
			nodeStatus.TrafficRate = v.TrafficRate
			nodestatusArr = append(nodestatusArr, nodeStatus)
		}
	}
	return &nodestatusArr
}

func (n *AdminNode) GetShadowsocksServerKey(node model.Node) string {
	switch node.Scy {
	case "2022-blake3-aes-128-gcm":
		return base64.StdEncoding.EncodeToString([]byte(node.ServerKey[:16]))
	case "2022-blake3-aes-256-gcm", "2022-blake3-chacha20-poly1305":
		return base64.StdEncoding.EncodeToString([]byte(node.ServerKey))
	default:
		return node.UUID
	}
}
