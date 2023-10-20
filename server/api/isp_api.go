package api

import (
	"AirGo/global"
	"AirGo/model"
	"AirGo/service"
	"AirGo/utils/encrypt_plugin"
	"AirGo/utils/isp_plugin"
	"AirGo/utils/jwt_plugin"
	"AirGo/utils/other_plugin"
	"AirGo/utils/response"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// 获取监控
func GetMonitorByUserID(ctx *gin.Context) {
	uIDInt, ok := other_plugin.GetUserIDFromGinContext(ctx)
	if !ok {
		response.Fail("获取信息,uID参数错误", nil, ctx)
		return
	}
	isp, _, err := service.CommonSqlFind[model.ISP, string, model.ISP]("user_id = " + fmt.Sprintf("%d", uIDInt))

	if err == gorm.ErrRecordNotFound {
		//创建新的
		var ispNew = model.ISP{
			UserID: uIDInt,
			UnicomConfig: model.UnicomConfig{
				APPID: encrypt_plugin.RandomString(160),
			},
			TelecomConfig: model.TelecomConfig{
				DeviceUid: encrypt_plugin.RandomString2(16),
			},
		}
		service.CommonSqlCreate[model.ISP](ispNew)
		isp, _, _ = service.CommonSqlFind[model.ISP, string, model.ISP]("user_id = " + fmt.Sprintf("%d", uIDInt))
	}
	response.OK("获取成功", isp, ctx)

}

// 发送验证码
func SendCode(ctx *gin.Context) {
	var isp model.ISP
	err := ctx.ShouldBind(&isp)
	if err != nil {
		global.Logrus.Error("运营商参数错误:", err)
		response.Fail("运营商参数错误", nil, ctx)
		return
	}
	//处理mobile
	mb, _ := encrypt_plugin.RSAEnCrypt(isp.Mobile, isp_plugin.UnicomPublicKey)
	isp.UnicomConfig.UnicomMobile = mb
	//
	var resp string
	switch isp.ISPType {
	case "unicom":
		resp, err = isp_plugin.UnicomCode(&isp)
	case "telecom":
		//fmt.Println("telecom")
		resp, err = isp_plugin.TelecomCode(&isp)
	}
	if err != nil {
		global.Logrus.Error("发送验证码错误:", err)
		response.Fail("发送验证码错误", nil, ctx)
		return
	}
	//判断是否为空
	respMap := make(map[string]interface{})
	err = json.Unmarshal([]byte(resp), &respMap)
	if err != nil {
		global.Logrus.Error("resp解析错误:", err)
		response.Fail("resp解析错误", nil, ctx)
		return
	}
	switch isp.ISPType {
	case "unicom":
		if respMap["rsp_code"].(string) != "0000" {
			msg := respMap["rsp_desc"].(string)
			if msg == "" {
				msg = "发送验证码失败"
			}
			response.Fail(msg, nil, ctx)
			return
		}

	case "telecom":
		if respMap["responseData"].(map[string]interface{})["resultCode"].(string) != "0000" {
			msg := respMap["responseData"].(map[string]interface{})["resultDesc"].(string)
			if msg == "" {
				msg = "发送验证码失败"
			}
			response.Fail(msg, nil, ctx)
			return
		}

	}
	response.OK("验证码发送成功", string(json.RawMessage(resp)), ctx)

}

// 登录运营商
func ISPLogin(ctx *gin.Context) {
	uIDInt, ok := other_plugin.GetUserIDFromGinContext(ctx)
	if !ok {
		response.Fail("获取信息,uID参数错误", nil, ctx)
		return
	}

	var isp model.ISP
	err := ctx.ShouldBind(&isp)
	if err != nil {
		global.Logrus.Error("运营商参数错误:", err)
		response.Fail("运营商参数错误", nil, ctx)
		return
	}
	//fmt.Println("登录：", isp)
	if isp.ISPType == "loginAgain" {
		//清空手机号信息，重新登录
		isp1, _, _ := service.CommonSqlFind[model.ISP, string, model.ISP]("user_id = " + fmt.Sprintf("%d", uIDInt))
		isp1.UnicomConfig.Cookie = ""
		isp1.UnicomConfig.Password = ""
		isp1.UnicomConfig.UnicomMobile = ""
		isp1.Status = false
		isp1.TelecomConfig.TelecomPassword = ""
		isp1.TelecomConfig.PhoneNum = ""
		isp1.TelecomConfig.TelecomToken = ""
		service.CommonSqlSave[model.ISP](isp1)
		response.OK("获取成功", isp1, ctx)
		return
	}

	var resp, cookie string
	switch isp.ISPType {
	case "unicom":
		mb, _ := encrypt_plugin.RSAEnCrypt(isp.Mobile, isp_plugin.UnicomPublicKey)
		pw, _ := encrypt_plugin.RSAEnCrypt(isp.UnicomConfig.Password, isp_plugin.UnicomPublicKey)
		isp.UnicomConfig.UnicomMobile = mb
		isp.UnicomConfig.Password = pw
		resp, cookie, err = isp_plugin.UnicomCodeLogin(isp.UnicomConfig.Password, isp.UnicomConfig.UnicomMobile, isp.UnicomConfig.APPID)
		if err != nil {
			global.Logrus.Error("尝试登录错误:", err)
			response.Fail("尝试登录错误", nil, ctx)
			return
		}
		if resp == "" || cookie == "" {
			global.Logrus.Error("尝试登录错误,resp,cookie为空")
			response.Fail("尝试登录错误", nil, ctx)
			return
		}
		//判断响应
		respMap := make(map[string]interface{})
		err = json.Unmarshal([]byte(resp), &respMap)

		if err != nil {
			global.Logrus.Error("resp解析错误:", err)
			response.Fail("resp解析错误", nil, ctx)
			return
		}
		if respMap["code"] != "0" {
			response.Fail(respMap["dsc"].(string), nil, ctx)
			return
		}
		//处理cookie，保存isp信息
		isp.UserID = uIDInt
		isp.UnicomConfig.Cookie = cookie
		isp.Status = true
		service.CommonSqlSave[model.ISP](isp)
		//response.OK("登录成功", resp, ctx)
		response.OK("登录成功", string(json.RawMessage(resp)), ctx)
	case "telecom":
		resp, err = isp_plugin.TelecomLogin(&isp)
		//判断响应
		respMap := make(map[string]interface{})
		err = json.Unmarshal([]byte(resp), &respMap)
		if respMap["responseData"].(map[string]interface{})["resultCode"].(string) != "0000" {
			response.Fail(respMap["responseData"].(map[string]interface{})["resultDesc"].(string), nil, ctx)
			return
		}
		//提取cookie
		cookie = respMap["responseData"].(map[string]interface{})["data"].(map[string]interface{})["loginSuccessResult"].(map[string]interface{})["token"].(string)
		//处理cookie，保存isp信息
		isp.UserID = uIDInt
		isp.TelecomConfig.TelecomToken = cookie
		isp.Status = true
		service.CommonSqlSave[model.ISP](isp)
		response.OK("登录成功", string(json.RawMessage(resp)), ctx)
	}

}

// 套餐查询
func QueryPackage(ctx *gin.Context) {
	token := ctx.Query("id")
	claims, err := jwt_plugin.ParseTokenHs256(token, global.Server.JWT.SigningKey)
	if err != nil {
		response.Fail(err.Error(), nil, ctx)
		return
	}
	//log.Println("token解析后 claims.ID：", claims.ID)
	//设置user id
	uID := claims.UserID
	//查询monitor
	isp, _, err := service.CommonSqlFind[model.ISP, string, model.ISP]("user_id = " + fmt.Sprintf("%d", uID))
	if err != nil {
		ctx.JSON(200, gin.H{
			"packageName": "查询流量失败，请重新登录",
		})
		return
	}
	//查询套餐
	var resp string
	switch isp.ISPType {
	case "unicom":
		resp, err = isp_plugin.UnicomQueryTraffic(&isp)
	case "telecom":
		resp, err = isp_plugin.TelecomQueryPackage(&isp)
	}
	if err != nil {
		//修改monitor状态
		isp.Status = false
		service.CommonSqlSave[model.ISP](isp)
		ctx.JSON(200, gin.H{
			"packageName": "查询流量失败，请重新登录",
			"mobile":      err.Error(),
		})
		return
	}
	ctx.String(200, resp) //将响应原样返回，避免json编码出现斜杠

}
