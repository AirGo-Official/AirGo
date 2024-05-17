package admin_api

import (
	"github.com/gin-gonic/gin"
	"github.com/ppoonk/AirGo/constant"
	"github.com/ppoonk/AirGo/global"
	"github.com/ppoonk/AirGo/model"
	"github.com/ppoonk/AirGo/service"
	"github.com/ppoonk/AirGo/utils/encrypt_plugin"
	"github.com/ppoonk/AirGo/utils/response"
)

// NewNode
// @Tags [admin api] node
// @Summary 新建节点
// @Produce json
// @Param Authorization header string false "Bearer 用户token"
// @Param data body model.Node true "参数"
// @Success 200 {object} response.ResponseStruct "请求成功；正常：业务代码 code=0；错误：业务代码code=1"
// @Router /api/admin/node/newNode [post]
func NewNode(ctx *gin.Context) {
	var node model.Node
	err := ctx.ShouldBind(&node)
	if err != nil {
		global.Logrus.Error(err.Error())
		response.Fail(constant.ERROR_REQUEST_PARAMETER_PARSING_ERROR+err.Error(), nil, ctx)
		return
	}
	err = service.AdminNodeSvc.NewNode(&node)
	if err != nil {
		global.Logrus.Error(err.Error())
		response.Fail("NewNode error:"+err.Error(), nil, ctx)
		return
	}
	response.OK("NewNode success", nil, ctx)
}

// DeleteNode
// @Tags [admin api] node
// @Summary 删除节点
// @Produce json
// @Param Authorization header string false "Bearer 用户token"
// @Param data body model.Node true "参数"
// @Success 200 {object} response.ResponseStruct "请求成功；正常：业务代码 code=0；错误：业务代码code=1"
// @Router /api/admin/node/deleteNode [delete]
func DeleteNode(ctx *gin.Context) {
	var node model.Node
	err := ctx.ShouldBind(&node)
	if err != nil {
		global.Logrus.Error(err.Error())
		response.Fail(constant.ERROR_REQUEST_PARAMETER_PARSING_ERROR+err.Error(), nil, ctx)
		return
	}
	err = service.AdminNodeSvc.DeleteNode(&node)
	if err != nil {
		global.Logrus.Error(err.Error())
		response.Fail("DeleteNode error:"+err.Error(), nil, ctx)
		return
	}
	response.OK("DeleteNode success", nil, ctx)
}

// UpdateNode
// @Tags [admin api] node
// @Summary 更新节点
// @Produce json
// @Param Authorization header string false "Bearer 用户token"
// @Param data body model.Node true "参数"
// @Success 200 {object} response.ResponseStruct "请求成功；正常：业务代码 code=0；错误：业务代码code=1"
// @Router /api/admin/node/updateNode [post]
func UpdateNode(ctx *gin.Context) {
	var node model.Node
	err := ctx.ShouldBind(&node)
	if err != nil {
		global.Logrus.Error(err.Error())
		response.Fail(constant.ERROR_REQUEST_PARAMETER_PARSING_ERROR+err.Error(), nil, ctx)
		return
	}
	err = service.AdminNodeSvc.UpdateNode(&node)
	if err != nil {
		global.Logrus.Error(err.Error())
		response.Fail("UpdateNode error:"+err.Error(), nil, ctx)
		return
	}
	response.OK("UpdateNode success", nil, ctx)

}

// GetNodeList
// @Tags [admin api] node
// @Summary 获取节点列表
// @Produce json
// @Param Authorization header string false "Bearer 用户token"
// @Param data body model.QueryParams true "参数"
// @Success 200 {object} response.ResponseStruct "请求成功；正常：业务代码 code=0；错误：业务代码code=1"
// @Router /api/admin/node/getNodeList [post]
func GetNodeList(ctx *gin.Context) {
	var params model.QueryParams
	err := ctx.ShouldBind(&params)
	if err != nil {
		global.Logrus.Error(err.Error())
		response.Fail(constant.ERROR_REQUEST_PARAMETER_PARSING_ERROR+err.Error(), nil, ctx)
		return
	}
	res, err := service.AdminNodeSvc.GetNodeList(&params)
	if err != nil {
		global.Logrus.Error(err.Error())
		response.Fail("GetNodeList error:"+err.Error(), nil, ctx)
		return
	}
	response.OK("GetNodeList success", res, ctx)
}

// GetNodeListWithTraffic
// @Tags [admin api] node
// @Summary 获取节点列表，带流量信息
// @Produce json
// @Param Authorization header string false "Bearer 用户token"
// @Param data body model.QueryParams true "参数"
// @Success 200 {object} response.ResponseStruct "请求成功；正常：业务代码 code=0；错误：业务代码code=1"
// @Router /api/admin/node/getNodeListWithTraffic [post]
func GetNodeListWithTraffic(ctx *gin.Context) {
	var params model.QueryParams
	err := ctx.ShouldBind(&params)
	if err != nil {
		global.Logrus.Error(err.Error())
		response.Fail(constant.ERROR_REQUEST_PARAMETER_PARSING_ERROR+err.Error(), nil, ctx)
		return
	}
	params.TableName = "node"
	res, err := service.AdminNodeSvc.GetNodeListWithTraffic(&params)
	if err != nil {
		global.Logrus.Error(err.Error())
		response.Fail("GetNodeListWithTraffic error:"+err.Error(), nil, ctx)
		return
	}
	response.OK("GetNodeListWithTraffic success", res, ctx)
}

// NodeSort
// @Tags [admin api] node
// @Summary 节点排序
// @Produce json
// @Param Authorization header string false "Bearer 用户token"
// @Param data body []model.Node true "参数"
// @Success 200 {object} response.ResponseStruct "请求成功；正常：业务代码 code=0；错误：业务代码code=1"
// @Router /api/admin/node/nodeSort [post]
func NodeSort(ctx *gin.Context) {
	var nodeArr []model.Node
	err := ctx.ShouldBind(&nodeArr)
	if err != nil {
		global.Logrus.Error(err.Error())
		response.Fail(constant.ERROR_REQUEST_PARAMETER_PARSING_ERROR+err.Error(), nil, ctx)
		return
	}
	err = service.CommonSqlUpdateMultiLine[[]model.Node](nodeArr, "id", []string{"node_order"})
	if err != nil {
		global.Logrus.Error(err.Error())
		response.Fail("NodeSort error:"+err.Error(), nil, ctx)
		return
	}
	response.OK("NodeSort success", nil, ctx)
}

// ParseUrl
// @Tags [admin api] node
// @Summary 解析
// @Produce json
// @Param Authorization header string false "Bearer 用户token"
// @Param data body model.NodeSharedReq true "参数"
// @Success 200 {object} response.ResponseStruct "请求成功；正常：业务代码 code=0；错误：业务代码code=1"
// @Router /api/admin/node/parseUrl [post]
func ParseUrl(ctx *gin.Context) {
	var url model.NodeSharedReq
	err := ctx.ShouldBind(&url)
	if err != nil {
		global.Logrus.Error(err.Error())
		response.Fail(constant.ERROR_REQUEST_PARAMETER_PARSING_ERROR+err.Error(), nil, ctx)
		return
	}
	nodeArr := service.AdminNodeSvc.ParseSubUrl(url.Url)
	response.OK("NewNodeShared success", nodeArr, ctx)
}

// NewNodeShared
// @Tags [admin api] node
// @Summary 新增共享节点
// @Produce json
// @Param Authorization header string false "Bearer 用户token"
// @Param data body []model.Node true "参数"
// @Success 200 {object} response.ResponseStruct "请求成功；正常：业务代码 code=0；错误：业务代码code=1"
// @Router /api/admin/node/newNodeShared [post]
func NewNodeShared(ctx *gin.Context) {
	var nodes []model.Node
	err := ctx.ShouldBind(&nodes)
	if err != nil {
		global.Logrus.Error(err.Error())
		response.Fail(constant.ERROR_REQUEST_PARAMETER_PARSING_ERROR+err.Error(), nil, ctx)
		return
	}
	for _, v := range nodes {
		_ = service.AdminNodeSvc.NewNode(&v)

	}
	response.OK("NewNodeShared success", nil, ctx)
}

// Createx25519
// @Tags [admin api] node
// @Summary reality x25519
// @Produce json
// @Param Authorization header string false "Bearer 用户token"
// @Success 200 {object} response.ResponseStruct "请求成功；正常：业务代码 code=0；错误：业务代码code=1"
// @Router /api/admin/node/createx25519 [get]
func Createx25519(ctx *gin.Context) {
	str := encrypt_plugin.RandomString(43)
	pub, pri, err := encrypt_plugin.ExecuteX25519(str)
	if err != nil {
		global.Logrus.Error(err.Error())
		return
	}
	response.OK("Createx25519 success", model.AGREALITYx25519{PublicKey: pub, PrivateKey: pri}, ctx)
}

// GetNodeServerStatus
// @Tags [admin api] node
// @Summary 获取节点服务器状态
// @Produce json
// @Param Authorization header string false "Bearer 用户token"
// @Success 200 {object} response.ResponseStruct "请求成功；正常：业务代码 code=0；错误：业务代码code=1"
// @Router /api/admin/node/getNodeServerStatus [get]
func GetNodeServerStatus(ctx *gin.Context) {
	list := service.AdminNodeSvc.GetNodesStatus()
	response.OK("GetNodeServerStatus success", list, ctx)

}
