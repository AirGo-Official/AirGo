package service

import (
	"AirGo/global"
	"AirGo/model"
	"strconv"
	"time"
)

// 查询节点流量
func GetNodeTraffic(params model.PaginationParams) model.NodesWithTotal {
	var nodesWithTotal model.NodesWithTotal
	var startTime, endTime time.Time
	//时间格式转换
	if len(params.Date) == 2 {
		startTime, _ = time.ParseInLocation("2006-01-02 15:04:05", params.Date[0], time.Local)
		endTime, _ = time.ParseInLocation("2006-01-02 15:04:05", params.Date[1], time.Local)
	} else {
		//默认前1个月数据
		endTime = time.Now().Local()
		startTime = endTime.AddDate(0, 0, -30)
	}
	if params.Search != "" {
		err := global.DB.Model(&model.Node{}).Count(&nodesWithTotal.Total).Where("remarks LIKE ?", "%"+params.Search+"%").Limit(int(params.PageSize)).Offset((int(params.PageNum)-1)*int(params.PageSize)).Preload("TrafficLogs", global.DB.Where("created_at > ? and created_at < ?", startTime, endTime)).Order("node_order").Find(&nodesWithTotal.NodeList).Error
		if err != nil {
			global.Logrus.Error("查询节点流量error:", err.Error())
			return model.NodesWithTotal{}
		}
	} else {
		err := global.DB.Model(&model.Node{}).Count(&nodesWithTotal.Total).Limit(int(params.PageSize)).Offset((int(params.PageNum)-1)*int(params.PageSize)).Preload("TrafficLogs", global.DB.Where("created_at > ? and created_at < ?", startTime, endTime)).Order("node_order").Find(&nodesWithTotal.NodeList).Error
		if err != nil {
			global.Logrus.Error("查询节点流量error:", err.Error())
			return model.NodesWithTotal{}
		}
	}
	for i1, _ := range nodesWithTotal.NodeList {
		for _, v := range nodesWithTotal.NodeList[i1].TrafficLogs {
			nodesWithTotal.NodeList[i1].TotalUp = nodesWithTotal.NodeList[i1].TotalUp + v.U
			nodesWithTotal.NodeList[i1].TotalDown = nodesWithTotal.NodeList[i1].TotalDown + v.D
		}
		//nodeArr[i1].TrafficLogs=[]model.TrafficLog{} //清空traffic
	}
	return nodesWithTotal
}

// 获取 node status，用于探针
func GetNodesStatus() *[]model.NodeStatus {
	var nodesIds []model.Node
	global.DB.Model(&model.Node{}).Select("id", "remarks", "traffic_rate").Order("node_order").Find(&nodesIds)
	var nodestatusArr []model.NodeStatus
	for _, v := range nodesIds {
		var nodeStatus = model.NodeStatus{}
		vStatus, ok := global.LocalCache.Get(strconv.FormatInt(v.ID, 10) + "status")
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
