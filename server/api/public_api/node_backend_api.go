package public_api

import (
	"encoding/base64"
	"fmt"
	"github.com/ppoonk/AirGo/api"
	"github.com/ppoonk/AirGo/constant"
	"github.com/ppoonk/AirGo/global"
	"github.com/ppoonk/AirGo/model"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

func AGGetNodeInfo(ctx *gin.Context) {
	//验证key
	if global.Server.Website.TEK != ctx.Query("key") {
		return
	}
	id := ctx.Query("id")
	nodeIDInt, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		ctx.AbortWithStatus(400)
		return
	}
	var node model.Node
	err = global.DB.Model(&model.Node{}).Where(&model.Node{ID: nodeIDInt}).Preload("Access").First(&node).Error
	if err != nil {
		global.Logrus.Error("AGGetNodeInfo error,id="+id, err.Error())
		return
	}
	//处理探针
	global.GoroutinePool.Submit(func() {
		//取消离线节点的通知状态
		global.LocalCache.Delete(fmt.Sprintf("%s%d", constant.CACHE_NODE_STATUS_IS_NOTIFIED_BY_NODEID, node.ID))

		cacheStatus, ok := global.LocalCache.Get(constant.CACHE_NODE_STATUS_BY_NODEID + id)
		if ok && cacheStatus != nil {
			oldStatus := cacheStatus.(model.NodeStatus)
			oldStatus.Status = true
			global.LocalCache.Set(constant.CACHE_NODE_STATUS_BY_NODEID+id, oldStatus, 2*time.Minute) //2分钟后过期

		} else {
			var status model.NodeStatus
			status.Status = true
			status.ID, _ = strconv.ParseInt(id, 64, 10)
			global.LocalCache.Set(constant.CACHE_NODE_STATUS_BY_NODEID+id, status, 2*time.Minute) //2分钟后过期
		}
	})
	//处理ss节点加密
	if node.Protocol == "shadowsocks" {
		switch node.Scy {
		case "2022-blake3-aes-128-gcm":
			node.ServerKey = base64.StdEncoding.EncodeToString([]byte(node.ServerKey[:16]))
		default:
			node.ServerKey = base64.StdEncoding.EncodeToString([]byte(node.ServerKey))
		}
	}
	//etag
	api.EtagHandler(node, ctx)
}

func AGReportNodeStatus(ctx *gin.Context) {
	//验证key
	if global.Server.Website.TEK != ctx.Query("key") {
		return
	}
	var AGNodeStatus model.AGNodeStatus
	err := ctx.ShouldBind(&AGNodeStatus)
	if err != nil {
		global.Logrus.Error("error", err.Error())
		ctx.AbortWithStatus(400)
		return
	}
	//处理探针
	global.GoroutinePool.Submit(func() {
		cacheStatus, ok := global.LocalCache.Get(fmt.Sprintf("%s%d", constant.CACHE_NODE_STATUS_BY_NODEID, AGNodeStatus.ID))
		if ok && cacheStatus != nil {
			oldStatus := cacheStatus.(model.NodeStatus)
			oldStatus.Status = true
			oldStatus.CPU = AGNodeStatus.CPU
			oldStatus.Mem = AGNodeStatus.Mem
			oldStatus.Disk = AGNodeStatus.Disk
			//oldStatus.Uptime=AGNodeStatus.Uptime
			global.LocalCache.Set(fmt.Sprintf("%s%d", constant.CACHE_NODE_STATUS_BY_NODEID, AGNodeStatus.ID), oldStatus, 2*time.Minute) //2分钟后过期
		} else {
			var status model.NodeStatus
			status.Status = true
			status.ID = AGNodeStatus.ID
			status.CPU = AGNodeStatus.CPU
			status.Mem = AGNodeStatus.Mem
			status.Disk = AGNodeStatus.Disk
			global.LocalCache.Set(fmt.Sprintf("%s%d", constant.CACHE_NODE_STATUS_BY_NODEID, AGNodeStatus.ID), status, 2*time.Minute) //2分钟后过期
		}
	})

	ctx.String(200, "success")
}

func AGGetUserlist(ctx *gin.Context) {
	//验证key
	if global.Server.Website.TEK != ctx.Query("key") {
		return
	}
	id := ctx.Query("id")
	nodeIDInt, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		ctx.AbortWithStatus(400)
		return
	}
	var node model.Node
	err = global.DB.Model(&model.Node{}).Where(&model.Node{ID: nodeIDInt}).First(&node).Error
	if err != nil {
		ctx.AbortWithStatus(400)
		return
	}
	//节点属于哪些goods
	goods, err := shopService.FindGoodsByNodeID(nodeIDInt)
	if err != nil {
		ctx.AbortWithStatus(400)
		return
	}
	//goods属于哪些用户
	var goodsArr []int64
	for _, v := range goods {
		goodsArr = append(goodsArr, v.ID)
	}
	var users []model.AGUserInfo
	err = global.DB.
		Model(&model.CustomerService{}).
		Where("goods_id in (?) and sub_status = ?", goodsArr, true).
		Select("id, sub_uuid AS uuid, user_name, node_connector, node_speed_limit").
		Find(&users).Error
	if err != nil {
		global.Logrus.Error("error,id="+id, err.Error())
		ctx.AbortWithStatus(400)
		return
	}
	//处理ss加密
	switch node.Protocol {
	case constant.NODE_PROTOCOL_SHADOWSOCKS:
		switch strings.HasPrefix(node.Scy, "2022") {
		case true:
			for k, _ := range users {
				p := users[k].UUID.String()
				if node.Scy == "2022-blake3-aes-128-gcm" {
					p = p[:16]
				}
				p = base64.StdEncoding.EncodeToString([]byte(p))
				users[k].Passwd = p
			}
		default:
			for k, _ := range users {
				users[k].Passwd = users[k].UUID.String()
			}
		}
	default:
	}
	//fmt.Println("users:", users)
	api.EtagHandler(users, ctx)
}

func ssEncryptionHandler(node model.Node, user *model.AGUserInfo) {
	switch node.Protocol {
	case constant.NODE_PROTOCOL_SHADOWSOCKS:
		if strings.HasPrefix(node.Scy, "2022") {
			//
			p := user.UUID.String()
			if node.Scy == "2022-blake3-aes-128-gcm" {
				p = p[:16]
			}
			p = base64.StdEncoding.EncodeToString([]byte(p))
			user.Passwd = p

		} else {
			user.Passwd = user.UUID.String()
		}
	default:

	}
}

func AGReportUserTraffic(ctx *gin.Context) {
	//验证key
	if global.Server.Website.TEK != ctx.Query("key") {
		return
	}
	var AGUserTraffic model.AGUserTraffic
	err := ctx.ShouldBind(&AGUserTraffic)
	if err != nil {
		global.Logrus.Error("error", err.Error())
		ctx.AbortWithStatus(400)
		return
	}
	//fmt.Println("上报用户流量：", AGUserTraffic)
	//查询节点倍率
	node, err := nodeService.FirstNode(&model.Node{ID: AGUserTraffic.ID})
	if err != nil {
		global.Logrus.Error("error", err.Error())
		ctx.AbortWithStatus(400)
		return
	}
	if node.TrafficRate < 0 {
		node.TrafficRate = 1
	}
	// 处理流量统计
	var userIds []int64
	var customerServiceArr []model.CustomerService
	var trafficLog = model.NodeTrafficLog{
		NodeID: node.ID,
	}
	//var userTrafficLog []model.UserTrafficLog
	userTrafficLogMap := make(map[int64]model.UserTrafficLog)
	for _, v := range AGUserTraffic.UserTraffic {
		//每个用户流量
		userIds = append(userIds, v.UID)
		//需要更新的用户订阅信息
		customerServiceArr = append(customerServiceArr, model.CustomerService{
			ID:       v.UID,
			UsedUp:   int64(float64(v.Upload) * node.TrafficRate),
			UsedDown: int64(float64(v.Download) * node.TrafficRate),
		})
		//需要插入的用户流量统计
		userTrafficLogMap[v.UID] = model.UserTrafficLog{
			SubUserID: v.UID,
			UserName:  v.Email,
			U:         int64(float64(v.Upload) * node.TrafficRate),
			D:         int64(float64(v.Download) * node.TrafficRate),
		}
		//该节点总流量
		trafficLog.D = trafficLog.U + v.Upload
		trafficLog.U = trafficLog.D + v.Download

	}
	// 处理探针
	global.GoroutinePool.Submit(func() {
		nodeService.UpdateNodeStatus(userIds, &trafficLog)
	})
	//插入节点流量统计
	global.GoroutinePool.Submit(func() {
		nodeService.UpdateNodeTraffic(&trafficLog, &AGUserTraffic)
	})
	//插入用户流量统计
	global.GoroutinePool.Submit(func() {
		admin_customerService.UpdateCustomerServiceTrafficLog(userTrafficLogMap, userIds)
	})
	//更新用户已用流量信息
	global.GoroutinePool.Submit(func() {
		admin_customerService.UpdateCustomerServiceTrafficUsed(&customerServiceArr, userIds)
	})
	ctx.String(200, "success")

}

func AGReportNodeOnlineUsers(ctx *gin.Context) {
	//验证key
	if global.Server.Website.TEK != ctx.Query("key") {
		return
	}
	var AGOnlineUser model.AGOnlineUser
	err := ctx.ShouldBind(&AGOnlineUser)
	if err != nil {
		global.Logrus.Error("error", err.Error())
		ctx.AbortWithStatus(400)
		return
	}
	ctx.String(200, "success")
}
