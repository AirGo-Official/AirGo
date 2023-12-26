package service

import (
	"github.com/ppoonk/AirGo/global"
	"github.com/ppoonk/AirGo/model"
	"strconv"
	"strings"
	"time"
)

// 查询节点流量
func GetNodeTraffic(params *model.FieldParamsReq) (*model.NodesWithTotal, error) {
	var nodesWithTotal model.NodesWithTotal
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
	err = global.DB.Model(&model.Node{}).Count(&nodesWithTotal.Total).Where(dataSql).Preload("TrafficLogs", global.DB.Where("created_at > ? and created_at < ?", startTime, endTime)).Preload("Access").Find(&nodesWithTotal.NodeList).Error
	if err != nil {
		return nil, err
	}
	for i1, _ := range nodesWithTotal.NodeList {
		//处理流量记录
		for _, v := range nodesWithTotal.NodeList[i1].TrafficLogs {
			nodesWithTotal.NodeList[i1].TotalUp = nodesWithTotal.NodeList[i1].TotalUp + v.U
			nodesWithTotal.NodeList[i1].TotalDown = nodesWithTotal.NodeList[i1].TotalDown + v.D
		}
		nodesWithTotal.NodeList[i1].TrafficLogs = []model.TrafficLog{} //清空traffic
		//处理关联的access
		nodesWithTotal.NodeList[i1].AccessIds = []int64{} //防止出现null
		for _, v := range nodesWithTotal.NodeList[i1].Access {
			nodesWithTotal.NodeList[i1].AccessIds = append(nodesWithTotal.NodeList[i1].AccessIds, v.ID)
		}
		nodesWithTotal.NodeList[i1].Access = []model.Access{}
	}
	return &nodesWithTotal, err
}

// 获取 node status，用于探针
func GetNodesStatus() *[]model.NodeStatus {
	var nodesIds []model.Node
	global.DB.Model(&model.Node{}).Select("id", "remarks", "traffic_rate").Where("enabled = ? AND enable_transfer = ?", true, false).Order("node_order").Find(&nodesIds)
	var nodestatusArr []model.NodeStatus
	for _, v := range nodesIds {
		var nodeStatus = model.NodeStatus{}
		vStatus, ok := global.LocalCache.Get(strconv.FormatInt(v.ID, 10) + global.NodeStatus)
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
func UpdateNode(node *model.Node) error {
	//查询关联access
	global.DB.Model(&model.Access{}).Where("id in ?", node.AccessIds).Find(&node.Access)
	//更新关联
	global.DB.Model(&node).Association("Access").Replace(&node.Access)
	//更新节点
	err := global.DB.Save(&node).Error
	//更新节点绑定的中转
	if !node.EnableTransfer { //当前更新的节点是一个落地直连节点
		var nodeArr []model.Node
		global.DB.Debug().Where("transfer_node_id = ?", node.ID).Find(&nodeArr)
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
			global.DB.Save(&nodeArr)
		}
	}
	return err
}

// 删除节点
func DeleteNode(node *model.Node) error {
	var funcs = []func() error{
		func() error { //删除商品关联的节点
			return global.DB.Where("node_id = ?", node.ID).Delete(&model.GoodsAndNodes{}).Error
		},
		func() error { //删除节点关联的访问控制
			return global.DB.Model(&model.Node{ID: node.ID}).Association("Access").Replace(nil)
		},
		func() error { //删除节点关联的流量统计信息
			return global.DB.Where("node_id = ?", node.ID).Delete(&model.TrafficLog{}).Error
		},
		func() error { //删除节点
			return global.DB.Delete(&node).Error
		},
	}
	for _, v := range funcs {
		err := v()
		if err != nil {
			return err
		}
	}
	return nil
}

// 删除节点,临时代码，处理之前版本删除节点遗留的数据库垃圾数据
func DeleteNodeTemp() error {
	var ids []int64
	global.DB.Model(&model.Node{}).Select("id").Find(&ids)

	var funcs = []func() error{
		func() error { //删除商品关联的节点
			return global.DB.Where("node_id NOT IN ?", ids).Delete(&model.GoodsAndNodes{}).Error
		},
		func() error { //删除节点关联的访问控制
			return global.DB.Where("node_id NOT IN ?", ids).Delete(&model.NodeAndAccess{}).Error
		},
		func() error { //删除节点关联的流量统计信息
			return global.DB.Where("node_id NOT IN ?", ids).Delete(&model.TrafficLog{}).Error
		},
	}
	for _, v := range funcs {
		err := v()
		if err != nil {
			return err
		}
	}
	return nil
}
