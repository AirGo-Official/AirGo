package public_api

import (
	"encoding/base64"
	"fmt"
	"github.com/ppoonk/AirGo/api"
	"github.com/ppoonk/AirGo/constant"
	"github.com/ppoonk/AirGo/global"
	"github.com/ppoonk/AirGo/model"
	"github.com/ppoonk/AirGo/service"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// AGGetNodeInfo
// @Tags [public api] node
// @Summary 获取节点配置信息
// @Produce json
// @Param id query int64 true "节点ID"
// @Param key query string true "节点密钥"
// @Success 200 {object} model.Node "成功"
// @Failure 400  "请求错误"
// @Failure 304  "数据和上次一致"
// @Router /api/public/airgo/node/getNodeInfo [get]
func AGGetNodeInfo(ctx *gin.Context) {
	//验证key
	if global.Server.Subscribe.TEK != ctx.Query("key") {
		ctx.AbortWithStatus(400)
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
		ctx.AbortWithStatus(400)
		return
	}
	//处理ss节点加密
	if node.Protocol == "shadowsocks" {
		node.ServerKey = service.AdminNodeSvc.GetShadowsocksServerKey(node)
	}
	//etag
	api.EtagHandler(node, ctx)
}

// AGReportNodeStatus
// @Tags [public api] node
// @Summary 上报节点状态
// @Produce json
// @Param key query string true "节点密钥"
// @Param data body model.AGNodeStatus true "参数"
// @Success 200 {object} string "成功"
// @Failure 400  "请求错误"
// @Failure 304  "数据和上次一致"
// @Router /api/public/airgo/node/AGReportNodeStatus [post]
func AGReportNodeStatus(ctx *gin.Context) {
	//验证key
	if global.Server.Subscribe.TEK != ctx.Query("key") {
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
	cacheStatus, ok := global.LocalCache.Get(fmt.Sprintf("%s%d",
		constant.CACHE_NODE_STATUS_BY_NODEID, AGNodeStatus.ID))
	if ok {
		oldStatus := cacheStatus.(model.NodeStatus)
		oldStatus.Status = true
		oldStatus.CPU = AGNodeStatus.CPU
		oldStatus.Mem = AGNodeStatus.Mem
		oldStatus.Disk = AGNodeStatus.Disk
		//oldStatus.Uptime=AGNodeStatus.Uptime
		global.LocalCache.Set(fmt.Sprintf("%s%d",
			constant.CACHE_NODE_STATUS_BY_NODEID, AGNodeStatus.ID),
			oldStatus,
			constant.CAHCE_NODE_STATUS_TIMEOUT*time.Minute)
	} else {
		var status model.NodeStatus
		status.Status = true
		status.ID = AGNodeStatus.ID
		status.CPU = AGNodeStatus.CPU
		status.Mem = AGNodeStatus.Mem
		status.Disk = AGNodeStatus.Disk
		global.LocalCache.Set(fmt.Sprintf("%s%d",
			constant.CACHE_NODE_STATUS_BY_NODEID, AGNodeStatus.ID),
			status,
			constant.CAHCE_NODE_STATUS_TIMEOUT*time.Minute)
	}
	ctx.String(200, "success")
}

// AGGetUserlist
// @Tags [public api] node
// @Summary 获取用户列表
// @Produce json
// @Param id query int64 true "节点ID"
// @Param key query string true "节点密钥"
// @Success 200 {object} string "成功"
// @Failure 400  "请求错误"
// @Failure 304  "数据和上次一致"
// @Router /api/public/airgo/user/AGGetUserlist [get]
func AGGetUserlist(ctx *gin.Context) {
	//验证key
	if global.Server.Subscribe.TEK != ctx.Query("key") {
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
	goods, err := service.AdminShopSvc.FindGoodsByNodeID(nodeIDInt)
	if err != nil {
		ctx.AbortWithStatus(400)
		return
	}
	//goods属于哪些用户
	var goodsArr []int64
	for _, v := range goods {
		goodsArr = append(goodsArr, v.ID)
	}
	var users []model.AGUserInfo //返回给节点服务器的数据，其中的 customer_server id 对应 Xrayr 或 v2bx 中的 uid; 处理上报流量时也要注意对应关系
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

// AGReportUserTraffic
// @Tags [public api] node
// @Summary 上报用户流量
// @Produce json
// @Param key query string true "节点密钥"
// @Param data body model.AGUserTraffic true "参数"
// @Success 200 {object} string "成功"
// @Failure 400  "请求错误"
// @Failure 304  "数据和上次一致"
// @Router /api/public/airgo/user/AGReportUserTraffic [post]
func AGReportUserTraffic(ctx *gin.Context) {
	//验证key
	if global.Server.Subscribe.TEK != ctx.Query("key") {
		return
	}
	var AGUserTraffic model.AGUserTraffic //Xrayr 或 v2bx 中的 uid 对应 customer_server id
	err := ctx.ShouldBind(&AGUserTraffic)
	if err != nil {
		global.Logrus.Error("error", err.Error())
		ctx.AbortWithStatus(400)
		return
	}
	//fmt.Println("用户流量统计", AGUserTraffic)
	//查询节点倍率
	node, err := service.AdminNodeSvc.FirstNode(&model.Node{ID: AGUserTraffic.ID})
	if err != nil {
		global.Logrus.Error("error", err.Error())
		ctx.AbortWithStatus(400)
		return
	}
	if node.TrafficRate < 0 {
		node.TrafficRate = 1
	}
	// 处理流量统计
	var customerServerIDs []int64
	var customerServiceArr []model.CustomerService
	var trafficLog = model.NodeTrafficLog{
		NodeID: node.ID,
	}
	userTrafficLogMap := make(map[int64]model.UserTrafficLog)
	for _, v := range AGUserTraffic.UserTraffic {
		//每个用户流量
		customerServerIDs = append(customerServerIDs, v.UID)
		//需要更新的用户订阅信息（*倍率）
		customerServiceArr = append(customerServiceArr, model.CustomerService{
			ID:       v.UID,
			UsedUp:   int64(float64(v.Upload) * node.TrafficRate),
			UsedDown: int64(float64(v.Download) * node.TrafficRate),
		})
		//需要插入的用户流量统计（*倍率）
		userTrafficLogMap[v.UID] = model.UserTrafficLog{
			SubUserID: v.UID,
			UserName:  v.Email,
			U:         int64(float64(v.Upload) * node.TrafficRate),
			D:         int64(float64(v.Download) * node.TrafficRate),
		}
		//该节点总流量（无需倍率）
		trafficLog.D = trafficLog.U + v.Upload
		trafficLog.U = trafficLog.D + v.Download

	}
	// 处理节点状态
	_ = global.Queue.Publish(constant.NODE_BACKEND_TASK, &service.NodeBackendServiceMessage{
		Title: constant.NODE_BACKEND_TASK_TITLE_NODE_STATUS,
		Data: &service.NodeStatusMessage{
			CustomerServerIDs: customerServerIDs,
			NodeTrafficLog:    &trafficLog,
		},
	})
	//插入节点流量统计
	_ = global.Queue.Publish(constant.NODE_BACKEND_TASK, &service.NodeBackendServiceMessage{
		Title: constant.NODE_BACKEND_TASK_TITLE_NODE_TRAFFIC,
		Data: &service.NodeTrafficMessage{
			NodeTrafficLog: &trafficLog,
			AGUserTraffic:  &AGUserTraffic,
		},
	})
	//插入用户流量统计
	_ = global.Queue.Publish(constant.NODE_BACKEND_TASK, &service.NodeBackendServiceMessage{
		Title: constant.NODE_BACKEND_TASK_TITLE_UPDATE_CUSTOMER_TRAFFICLOG,
		Data: &service.UpdateCustomerTrafficLogMessage{
			CustomerServerIDs: customerServerIDs,
			UserTrafficLogMap: userTrafficLogMap,
		},
	})
	//更新用户已用流量信息
	_ = global.Queue.Publish(constant.NODE_BACKEND_TASK, &service.NodeBackendServiceMessage{
		Title: constant.NODE_BACKEND_TASK_TITLE_UPDATE_CUSTOMER_TRAFFICUSED,
		Data: &service.UpdateCustomerTrafficUsedMessage{
			CustomerServerIDs:   customerServerIDs,
			CustomerServiceList: &customerServiceArr,
		},
	})
	ctx.String(200, "success")

}

// AGReportNodeOnlineUsers
// @Tags [public api] node
// @Summary 上报在线用户
// @Produce json
// @Param key query string true "节点密钥"
// @Param data body model.AGOnlineUser true "参数"
// @Success 200 {object} string "成功"
// @Failure 400  "请求错误"
// @Failure 304  "数据和上次一致"
// @Router /api/public/airgo/user/AGReportNodeOnlineUsers [post]
func AGReportNodeOnlineUsers(ctx *gin.Context) {
	//验证key
	if global.Server.Subscribe.TEK != ctx.Query("key") {
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
	//TODO 未用到
}
