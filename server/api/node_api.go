package api

import (
	"AirGo/global"
	"AirGo/model"
	"AirGo/service"
	"AirGo/utils/encrypt_plugin"
	"AirGo/utils/response"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
)

func Show(data any) {
	b, _ := json.Marshal(data)
	fmt.Println(string(b))
}

// 获取全部节点
func GetAllNode(ctx *gin.Context) {
	nodeArr, _, err := service.CommonSqlFind[model.Node, string, []model.Node]("")
	if err != nil {
		global.Logrus.Error("获取全部节点错误:", err)
		response.Fail("获取全部节点错误", nil, ctx)
		return
	}
	response.OK("获取全部节点成功", nodeArr, ctx)
}

// 新建节点
func NewNode(ctx *gin.Context) {
	var node model.Node
	err := ctx.ShouldBind(&node)
	if err != nil {
		global.Logrus.Error("新建节点参数错误:", err)
		response.Fail("新建节点参数错误", nil, ctx)
		return
	}
	node.ServerKey = encrypt_plugin.RandomString(32)
	//Show(node)
	err = service.CommonSqlCreate[model.Node](node)
	if err != nil {
		global.Logrus.Error("新建节点错误:", err)
		response.Fail("新建节点错误", nil, ctx)
		return
	}
	response.OK("新建节点成功", nil, ctx)
}

// 删除节点
func DeleteNode(ctx *gin.Context) {
	var node model.Node
	err := ctx.ShouldBind(&node)
	if err != nil {
		global.Logrus.Error("删除节点参数错误:", err)
		response.Fail("删除节点参数错误", nil, ctx)
		return
	}
	err = service.CommonSqlDelete[model.Node, model.Node](model.Node{}, node)
	if err != nil {
		global.Logrus.Error("删除节点错误:", err)
		response.Fail("删除节点错误", nil, ctx)
		return
	}
	response.OK("删除节点成功", nil, ctx)
}

// 更新节点
func UpdateNode(ctx *gin.Context) {
	var node model.Node
	err := ctx.ShouldBind(&node)
	if err != nil {
		global.Logrus.Error("更新节点参数错误:", err)
		response.Fail("更新节点参数错误", nil, ctx)
		return
	}
	err = service.CommonSqlSave[model.Node](node)
	if err != nil {
		global.Logrus.Error("更新节点错误:", err)
		response.Fail("更新节点错误", nil, ctx)
		return
	}
	response.OK("更新节点成功", nil, ctx)

}

// 查询节点流量
func GetNodeTraffic(ctx *gin.Context) {
	var trafficParams model.PaginationParams
	err := ctx.ShouldBind(&trafficParams)
	if err != nil {
		global.Logrus.Error("查询节点错误:", err)
		response.Fail("查询节点流量参数错误"+err.Error(), nil, ctx)
		return
	}
	res := service.GetNodeTraffic(trafficParams)
	response.OK("查询节点成功", res, ctx)
}

// 节点排序
func NodeSort(ctx *gin.Context) {
	var nodeArr []model.Node
	err := ctx.ShouldBind(&nodeArr)
	if err != nil {
		global.Logrus.Error("节点排序参数错误:", err)
		response.Fail("节点排序参数错误:"+err.Error(), nil, ctx)
		return
	}
	err = service.CommonSqlUpdateMultiLine[[]model.Node](nodeArr, "id", []string{"node_order"})
	if err != nil {
		global.Logrus.Error("节点排序错误:", err)
		response.Fail("节点排序错误:"+err.Error(), nil, ctx)
		return
	}
	response.OK("节点排序成功", nil, ctx)
}
