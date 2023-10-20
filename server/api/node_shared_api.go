package api

import (
	"AirGo/global"
	"AirGo/model"
	"AirGo/service"
	"AirGo/utils/response"
	"github.com/gin-gonic/gin"
)

// 新增节点
func NewNodeShared(ctx *gin.Context) {
	var url model.NodeSharedReq
	err := ctx.ShouldBind(&url)
	if err != nil {
		global.Logrus.Error("NewNodeShared", err.Error())
		response.Fail("NewNodeShared"+err.Error(), nil, ctx)
		return
	}
	nodeArr := service.ParseUrl(url.Url)
	if nodeArr != nil {
		err = service.CommonSqlCreate[[]model.NodeShared](*nodeArr)
		if err != nil {
			global.Logrus.Error("NewNodeShared", err.Error())
			response.Fail("NewNodeShared"+err.Error(), nil, ctx)
			return
		}
		response.OK("新增节点成功", nil, ctx)
	}
}

// 获取节点列表
func GetNodeSharedList(ctx *gin.Context) {
	nodeArr, _, err := service.CommonSqlFind[model.NodeShared, string, []model.Node]("")
	if err != nil {
		global.Logrus.Error("GetNodeSharedList", err.Error())
		response.Fail("GetNodeSharedList"+err.Error(), nil, ctx)
		return
	}
	response.OK("获取节点列表成功", nodeArr, ctx)

}

// 删除节点
func DeleteNodeShared(ctx *gin.Context) {
	var node model.NodeShared
	err := ctx.ShouldBind(&node)
	if err != nil {
		global.Logrus.Error("DeleteNodeShared", err.Error())
		response.Fail("DeleteNodeShared"+err.Error(), nil, ctx)
		return
	}
	err = service.CommonSqlDelete[model.Node, model.NodeShared](model.Node{}, node)
	if err != nil {
		global.Logrus.Error("DeleteNodeShared", err.Error())
		response.Fail("DeleteNodeShared"+err.Error(), nil, ctx)
		return
	}
	response.OK("删除节点成功", nil, ctx)
}
