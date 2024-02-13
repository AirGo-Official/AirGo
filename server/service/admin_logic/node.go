package admin_logic

import (
	"errors"
	"fmt"
	"github.com/ppoonk/AirGo/constant"
	"github.com/ppoonk/AirGo/global"
	"github.com/ppoonk/AirGo/model"
	"github.com/ppoonk/AirGo/service/common_logic"
	"gorm.io/gorm"
	"strconv"
	"strings"
	"time"
)

type Node struct {
}

// 更新节点流量记录
func (n *Node) UpdateNodeTraffic(trafficLog *model.NodeTrafficLog, AGUserTraffic *model.AGUserTraffic) {
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
		nodeTraffic.U = trafficLog.U + trafficLog.U
		nodeTraffic.D = trafficLog.D + trafficLog.D
		global.DB.Save(&nodeTraffic)
	}
}

// 清理节点流量记录
func (n *Node) ClearNodeTraffic() error {
	y, m, _ := time.Now().Date()
	startTime := time.Date(y, m-2, 1, 0, 0, 0, 0, time.Local) //清除2个月之前的数据
	return global.DB.Transaction(func(tx *gorm.DB) error {
		return tx.Where("created_at < ?", startTime).Delete(&model.UserTrafficLog{}).Error
	})
}

// 查询节点
func (n *Node) FirstNode(nodeParams *model.Node) (*model.Node, error) {
	var node model.Node
	err := global.DB.Where(&nodeParams).First(&node).Error
	return &node, err
}

// 查询节点流量
func (n *Node) GetNodeListWithTraffic(params *model.QueryParams) (*model.CommonDataResp, error) {
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
	_, dataSql := common_logic.CommonSqlFindSqlHandler(params)
	dataSql = dataSql[strings.Index(dataSql, "WHERE ")+6:] //去掉`WHERE `
	if dataSql == "" {
		dataSql = "id > 0" //当前端什么参数没有传时，默认添加一个参数
	}
	err = global.DB.Model(&model.Node{}).
		Count(&total).Where(dataSql).
		//Preload("TrafficLogs", global.DB.Where("created_at > ? and created_at < ?", startTime, endTime)).
		Preload("TrafficLogs", "created_at > ? and created_at < ?", startTime, endTime).
		Preload("Access").
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
	return &model.CommonDataResp{total, nodeList}, err
}

// 更新node status
func (n *Node) UpdateNodeStatus(userIds []int64, trafficLog *model.NodeTrafficLog) {
	var duration float64 = 60 //默认60秒间隔
	cacheStatus, ok := global.LocalCache.Get(fmt.Sprintf("%s%d", constant.CACHE_NODE_STATUS_BY_NODEID, trafficLog.ID))
	if ok && cacheStatus != nil {
		oldStatus := cacheStatus.(model.NodeStatus)
		oldStatus.Status = true
		oldStatus.UserAmount = int64(len(userIds))
		now := time.Now()
		duration = now.Sub(oldStatus.LastTime).Seconds()
		oldStatus.D, _ = strconv.ParseFloat(fmt.Sprintf("%.2f", float64(trafficLog.D)/duration), 64) //Byte per second
		oldStatus.U, _ = strconv.ParseFloat(fmt.Sprintf("%.2f", float64(trafficLog.U)/duration), 64)
		oldStatus.LastTime = now
		global.LocalCache.Set(fmt.Sprintf("%s%d", constant.CACHE_NODE_STATUS_BY_NODEID, trafficLog.ID), oldStatus, 2*time.Minute)
	} else {
		var nodeStatus model.NodeStatus
		nodeStatus.Status = true
		nodeStatus.ID = trafficLog.NodeID
		nodeStatus.UserAmount = int64(len(userIds))
		nodeStatus.D, _ = strconv.ParseFloat(fmt.Sprintf("%.2f", float64(trafficLog.D)/duration), 64) //Byte per second
		nodeStatus.U, _ = strconv.ParseFloat(fmt.Sprintf("%.2f", float64(trafficLog.U)/duration), 64)
		nodeStatus.LastTime = time.Now()
		global.LocalCache.Set(fmt.Sprintf("%s%d", constant.CACHE_NODE_STATUS_BY_NODEID, trafficLog.ID), nodeStatus, 2*time.Minute)
	}
}

// 获取 node status
func (n *Node) GetNodesStatus() *[]model.NodeStatus {
	var nodesIds []model.Node
	global.DB.Model(&model.Node{}).Select("id", "remarks", "traffic_rate").Where("enabled = ? AND enable_transfer = ?", true, false).Order("node_order").Find(&nodesIds)
	var nodestatusArr []model.NodeStatus
	for _, v := range nodesIds {
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

// 更新节点
func (n *Node) UpdateNode(node *model.Node) error {
	return global.DB.Transaction(func(tx *gorm.DB) error {
		//查询关联access
		var accessIDs []int64
		for _, v := range node.Access {
			accessIDs = append(accessIDs, v.ID)
		}
		err := tx.Model(&model.Access{}).Where("id in ?", accessIDs).Find(&node.Access).Error
		if err != nil {
			return err
		}
		//更新关联
		err = tx.Model(&node).Association("Access").Replace(&node.Access)
		if err != nil {
			return err
		}
		//更新节点
		err = tx.Save(&node).Error
		if err != nil {
			return err
		}
		//更新节点绑定的中转
		if !node.EnableTransfer { //当前更新的节点是一个落地直连节点
			var nodeArr []model.Node
			err = tx.Where(&model.Node{TransferNodeID: node.ID}).Find(&nodeArr).Error
			if err != nil {
				return err
			}
			if len(nodeArr) > 0 {
				for k, v := range nodeArr { //遍历中转节点
					temp := *node
					temp.ID = v.ID
					temp.CreatedAt, temp.UpdatedAt = v.CreatedAt, v.UpdatedAt
					temp.Remarks = v.Remarks
					temp.EnableTransfer = true
					temp.TransferNodeID = v.TransferNodeID
					temp.TransferAddress = v.TransferAddress
					temp.TransferPort = v.TransferPort
					nodeArr[k] = temp
				}
				return tx.Save(&nodeArr).Error
			}
		}
		return nil
	})
}

// 删除节点
func (n *Node) DeleteNode(node *model.Node) error {
	return global.DB.Transaction(func(tx *gorm.DB) error {
		err := tx.Model(&model.Node{ID: node.ID}).Association("Goods").Replace(nil)
		if err != nil {
			return err
		}
		err = tx.Model(&model.Node{ID: node.ID}).Association("Access").Replace(nil)
		if err != nil {
			return err
		}
		err = tx.Where("node_id = ?", node.ID).Delete(&model.NodeTrafficLog{}).Error
		if err != nil {
			return err
		}
		return tx.Delete(&node).Error
	})
}
