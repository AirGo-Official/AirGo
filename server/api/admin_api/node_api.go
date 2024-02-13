package admin_api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/ppoonk/AirGo/constant"
	"github.com/ppoonk/AirGo/global"
	"github.com/ppoonk/AirGo/model"
	"github.com/ppoonk/AirGo/service/admin_logic"
	"github.com/ppoonk/AirGo/service/common_logic"
	"github.com/ppoonk/AirGo/utils/encrypt_plugin"
	"github.com/ppoonk/AirGo/utils/response"
	"time"
)

var nodeService admin_logic.Node

// 获取全部节点
func GetAllNode(ctx *gin.Context) {
	nodeArr, total, err := common_logic.CommonSqlFind[model.Node, string, []model.Node]("")
	if err != nil {
		global.Logrus.Error(err.Error())
		response.Fail("GetAllNode error:"+err.Error(), nil, ctx)
		return
	}
	response.OK("GetAllNode success", &model.CommonDataResp{total, nodeArr}, ctx)
}

// 新建节点
func NewNode(ctx *gin.Context) {
	var node model.Node
	err := ctx.ShouldBind(&node)
	if err != nil {
		global.Logrus.Error(err.Error())
		response.Fail(constant.ERROR_REQUEST_PARAMETER_PARSING_ERROR+err.Error(), nil, ctx)
		return
	}
	//fmt.Println("新建节点")
	//service.Show(node)
	n, _, _ := common_logic.CommonSqlFirst[model.Node, string, model.Node](fmt.Sprintf("remarks = '%s'", node.Remarks))
	if n.Remarks != "" {
		response.Fail("Node name is duplicate", nil, ctx)
		return
	}
	//根据节点类型，修改一些默认参数
	switch node.NodeType {
	case constant.NODE_TYPE_VMESS:
	case constant.NODE_TYPE_VLESS:
	case constant.NODE_TYPE_TROJAN:
	case constant.NODE_TYPE_HYSTERIA:
	case constant.NODE_TYPE_SHADOWSOCKS:
		node.ServerKey = encrypt_plugin.RandomString(32)
	case constant.NODE_TYPE_TRANSFER:
		//查询中转绑定节点
		n, _, err = common_logic.CommonSqlFirst[model.Node, string, model.Node](fmt.Sprintf("id = %d", node.TransferNodeID))
		if err != nil {
			global.Logrus.Error(err.Error())
			response.Fail("NewNode error:"+err.Error(), nil, ctx)
			return
		}
		//fmt.Println("查询中转绑定节点 n:", n)
		n.ID = 0
		n.CreatedAt, n.UpdatedAt = time.Now(), time.Now()
		n.Remarks = node.Remarks
		n.EnableTransfer = true
		n.TransferNodeID = node.TransferNodeID
		n.TransferAddress = node.TransferAddress
		n.TransferPort = node.TransferPort
		node = n
	}
	err = common_logic.CommonSqlCreate[model.Node](node)
	if err != nil {
		global.Logrus.Error(err.Error())
		response.Fail("NewNode error:"+err.Error(), nil, ctx)
		return
	}
	response.OK("NewNode success", nil, ctx)
}

// 删除节点
func DeleteNode(ctx *gin.Context) {
	var node model.Node
	err := ctx.ShouldBind(&node)
	if err != nil {
		global.Logrus.Error(err.Error())
		response.Fail(constant.ERROR_REQUEST_PARAMETER_PARSING_ERROR+err.Error(), nil, ctx)
		return
	}
	err = nodeService.DeleteNode(&node)
	if err != nil {
		global.Logrus.Error(err.Error())
		response.Fail("DeleteNode error:"+err.Error(), nil, ctx)
		return
	}
	response.OK("DeleteNode success", nil, ctx)
}

// 更新节点
func UpdateNode(ctx *gin.Context) {
	var node model.Node
	err := ctx.ShouldBind(&node)
	if err != nil {
		global.Logrus.Error(err.Error())
		response.Fail(constant.ERROR_REQUEST_PARAMETER_PARSING_ERROR+err.Error(), nil, ctx)
		return
	}
	err = nodeService.UpdateNode(&node)
	if err != nil {
		global.Logrus.Error(err.Error())
		response.Fail("UpdateNode error:"+err.Error(), nil, ctx)
		return
	}
	response.OK("UpdateNode success", nil, ctx)

}

// 查询节点流量
func GetNodeListWithTraffic(ctx *gin.Context) {
	var params model.QueryParams
	err := ctx.ShouldBind(&params)
	if err != nil {
		global.Logrus.Error(err.Error())
		response.Fail(constant.ERROR_REQUEST_PARAMETER_PARSING_ERROR+err.Error(), nil, ctx)
		return
	}
	res, err := nodeService.GetNodeListWithTraffic(&params)
	if err != nil {
		global.Logrus.Error(err.Error())
		response.Fail("GetNodeListWithTraffic error:"+err.Error(), nil, ctx)
		return
	}
	response.OK("GetNodeListWithTraffic success", res, ctx)
}

// 节点排序
func NodeSort(ctx *gin.Context) {
	var nodeArr []model.Node
	err := ctx.ShouldBind(&nodeArr)
	if err != nil {
		global.Logrus.Error(err.Error())
		response.Fail(constant.ERROR_REQUEST_PARAMETER_PARSING_ERROR+err.Error(), nil, ctx)
		return
	}
	err = common_logic.CommonSqlUpdateMultiLine[[]model.Node](nodeArr, "id", []string{"node_order"})
	if err != nil {
		global.Logrus.Error(err.Error())
		response.Fail("NodeSort error:"+err.Error(), nil, ctx)
		return
	}
	response.OK("NodeSort success", nil, ctx)
}

// 新增共享节点
func NewNodeShared(ctx *gin.Context) {
	var url model.NodeSharedReq
	err := ctx.ShouldBind(&url)
	if err != nil {
		global.Logrus.Error(err.Error())
		response.Fail(constant.ERROR_REQUEST_PARAMETER_PARSING_ERROR+err.Error(), nil, ctx)
		return
	}
	nodeArr := nodeService.ParseSubUrl(url.Url)
	if nodeArr != nil {
		for _, v := range *nodeArr {
			n, _, _ := common_logic.CommonSqlFind[model.NodeShared, model.NodeShared, model.NodeShared](model.NodeShared{
				Remarks: v.Remarks,
			})
			if n.Remarks != "" {
				continue
			}
			err = common_logic.CommonSqlCreate[[]model.NodeShared](*nodeArr)
			if err != nil {
				global.Logrus.Error(err.Error())
				response.Fail("NewNodeShared error:"+err.Error(), nil, ctx)
				return
			}
		}
		response.OK("NewNodeShared success", nil, ctx)
	}
}

// 获取共享节点列表
func GetNodeSharedList(ctx *gin.Context) {
	nodeArr, total, err := common_logic.CommonSqlFind[model.NodeShared, string, []model.NodeShared]("")
	if err != nil {
		global.Logrus.Error(err.Error())
		response.Fail("GetNodeSharedList"+err.Error(), nil, ctx)
		return
	}
	response.OK("GetNodeSharedList success", &model.CommonDataResp{total, nodeArr}, ctx)

}

// 删除共享节点
func DeleteNodeShared(ctx *gin.Context) {
	var node model.NodeShared
	err := ctx.ShouldBind(&node)
	if err != nil {
		global.Logrus.Error(err.Error())
		response.Fail(constant.ERROR_REQUEST_PARAMETER_PARSING_ERROR+err.Error(), nil, ctx)
		return
	}
	err = common_logic.CommonSqlDelete[model.Node, model.NodeShared](node)
	if err != nil {
		global.Logrus.Error(err.Error())
		response.Fail("DeleteNodeShared error:"+err.Error(), nil, ctx)
		return
	}
	response.OK("DeleteNodeShared success", nil, ctx)
}

// reality x25519
func Createx25519(ctx *gin.Context) {
	str := encrypt_plugin.RandomString(43)
	pub, pri, err := encrypt_plugin.ExecuteX25519(str)
	if err != nil {
		global.Logrus.Error(err.Error())
		return
	}
	response.OK("Createx25519 success", model.AGREALITYx25519{PublicKey: pub, PrivateKey: pri}, ctx)
}
