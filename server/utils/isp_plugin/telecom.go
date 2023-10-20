package isp_plugin

import (
	"AirGo/model"
	"AirGo/utils/net_plugin"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"
)

const (
	TelecomCodeUrl   = "https://appgologin.189.cn:9031/login/client/getLoginRandomCode"
	TelecomLoginUrl  = "https://appgologin.189.cn:9031/login/client/userLoginNormal"
	TelecomQueryUrl  = "https://appfuwu.189.cn:9021/query/userFluxPackage"
	TelecomPublicKey = "MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQDBkLT15ThVgz6/NOl6s8GNPofdWzWbCkWnkaAm7O2LjkM1H7dMvzkiqdxU02jamGRHLX/ZNMCXHnPcW/sDhiFCBN18qFvy8g6VYb9QtroI09e176s+ZCtiv7hbin2cCTj99iUpnEloZm19lwHyo69u5UMiPMpq0/XKBO8lYhN/gwIDAQAB"
)

// 发送登录验证码
func TelecomCode(isp *model.ISP) (string, error) {
	client := net_plugin.ClientWithDNS("114.114.114.114", 10*time.Second)
	timestamp := time.Now().Format("20060102150405")

	var jsonMap = map[string]interface{}{
		"headerInfos": map[string]interface{}{
			"code":           "getLoginRandomCode",
			"timestamp":      timestamp,
			"clientType":     "#10.4.1#channel38#HUAWEI MNA-AL00#",
			"shopId":         "20002",
			"source":         "110003",
			"sourcePassword": "Sid98s",
			"token":          "",
			"userLoginName":  "",
		},
		"content": map[string]interface{}{
			"attach": "test",
			"fieldData": map[string]interface{}{
				"payType":        "",
				"phoneNum":       isp.TelecomConfig.PhoneNum,
				"validationCode": "",
				"imsi":           "",
				"salesProdId":    "",
				"key":            "",
				"scene":          "55",
			},
		},
	}
	jsonValues, err := json.Marshal(jsonMap)
	if err != nil {
		return "", err
	}

	formBytesReader := bytes.NewReader(jsonValues)

	req, err := http.NewRequest("POST", TelecomCodeUrl, formBytesReader)
	if err != nil {
		return "", err
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")
	req.Header.Set("Connection", "Keep-Alive")
	req.Header.Set("Accept-Encoding", "gzip")
	//fmt.Println("请求参数：", req)
	resp, err := client.Do(req)
	defer resp.Body.Close()
	if err != nil {
		return "", err
	}
	out := net_plugin.ReadDate(resp)
	return out, err
}

// 电信登录
func TelecomLogin(isp *model.ISP) (string, error) {
	client := net_plugin.ClientWithDNS("114.114.114.114", 10*time.Second)
	var jsonMap = map[string]interface{}{
		"headerInfos": map[string]interface{}{
			"code":           "userLoginNormal",
			"timestamp":      isp.TelecomConfig.Timestamp,
			"broadAccount":   "",
			"broadToken":     "",
			"clientType":     "#9.6.1#channel50#iPhone 14 Pro#",
			"shopId":         "20002",
			"source":         "110003",
			"sourcePassword": "Sid98s",
			"token":          "",
			"userLoginName":  isp.Mobile,
		},
		"content": map[string]interface{}{
			"attach": "test",
			"fieldData": map[string]interface{}{
				"loginType":                  "4",
				"accountType":                "",
				"loginAuthCipherAsymmertric": isp.TelecomConfig.LoginAuthCipherAsymmertric,
				"deviceUid":                  isp.TelecomConfig.DeviceUid,
				"phoneNum":                   isp.TelecomConfig.PhoneNum,
				"isChinatelecom":             "0",
				"systemVersion":              "15.4.0",
				"authentication":             isp.TelecomConfig.TelecomPassword,
			},
		},
	}

	jsonValues, err := json.Marshal(jsonMap)
	if err != nil {
		return "", err
	}
	//fmt.Println("电信登录json：", string(jsonValues))
	formBytesReader := bytes.NewReader(jsonValues)

	req, err := http.NewRequest("POST", TelecomLoginUrl, formBytesReader)
	if err != nil {
		return "", err
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")
	req.Header.Set("Connection", "Keep-Alive")
	req.Header.Set("Accept-Encoding", "gzip")
	//fmt.Println("请求参数：", req)
	resp, err := client.Do(req)
	defer resp.Body.Close()
	if err != nil {
		return "", err
	}
	out := net_plugin.ReadDate(resp)
	return out, err

}

// 电信查询套餐详情
func TelecomQueryPackage(isp *model.ISP) (string, error) {
	client := net_plugin.ClientWithDNS("114.114.114.114", 10*time.Second)
	timestamp := time.Now().Format("20060102150405")
	var jsonMap = map[string]interface{}{
		"content": map[string]interface{}{
			"fieldData": map[string]interface{}{
				"queryFlag":  "0",
				"accessAuth": "1",
				"account":    isp.TelecomConfig.PhoneNum,
			},
			"attach": "test",
		},
		"headerInfos": map[string]interface{}{
			"clientType":     "#9.6.1#channel50#iPhone 14 Pro#",
			"timestamp":      timestamp,
			"code":           "userFluxPackage",
			"shopId":         "20002",
			"source":         "110003",
			"sourcePassword": "Sid98s",
			"token":          isp.TelecomConfig.TelecomToken,
			"userLoginName":  isp.Mobile,
		},
	}
	jsonValues, err := json.Marshal(jsonMap)
	if err != nil {
		return "", err
	}
	//fmt.Println("电信查询流量：", string(jsonValues))
	formBytesReader := bytes.NewReader(jsonValues)
	req, err := http.NewRequest("POST", TelecomQueryUrl, formBytesReader)
	if err != nil {
		return "", err
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")
	req.Header.Set("Connection", "Keep-Alive")
	req.Header.Set("Accept-Encoding", "gzip")

	resp, err := client.Do(req)
	defer resp.Body.Close()
	if err != nil {
		return "", err
	}
	respData := net_plugin.ReadDate(resp)

	//fmt.Println("respData", respData)
	return TelecomQueryTrafficHandler(respData, isp.Mobile)
}
func TelecomQueryTrafficHandler(resp, mobile string) (string, error) {
	//判断响应

	respMap := make(map[string]interface{})
	err := json.Unmarshal([]byte(resp), &respMap)
	if err != nil {

		fmt.Println("resp解析错误:", err)
		return "", err
	}
	//判断res.code
	if respMap["responseData"].(map[string]interface{})["resultCode"].(string) != "0000" {
		desc, ok := respMap["responseData"].(map[string]interface{})["resultDesc"].(string)
		if !ok {
			desc = "查询流量失败"
		}
		fmt.Println(desc)
		return "", errors.New(desc)
	}

	//手机号
	newMobile := mobile[0:3] + "****" + mobile[7:]
	//查询日期
	date := time.Now().Format("2006-01-02 15:04:05")
	//套餐包
	packageName := respMap["responseData"].(map[string]interface{})["data"].(map[string]interface{})["mainProductOFFInfo"].(map[string]interface{})["productOFFName"].(string)
	//流量包列表
	ratableResourcePackages := respMap["responseData"].(map[string]interface{})["data"].(map[string]interface{})["productOFFRatable"].(map[string]interface{})["ratableResourcePackages"].([]interface{})

	//国内通用流量,判断是否为空!!!!!!!!!!!!!!!!!
	var generic string
	genericUse := ratableResourcePackages[0].(map[string]interface{})["leftStructure"].(map[string]interface{})["title"].(string) + "---" + ratableResourcePackages[0].(map[string]interface{})["leftStructure"].(map[string]interface{})["num"].(string) + ratableResourcePackages[0].(map[string]interface{})["leftStructure"].(map[string]interface{})["unit"].(string)
	genericRemainRightStructure, ok := ratableResourcePackages[0].(map[string]interface{})["rightStructure"].(map[string]interface{})
	if ok {
		specialRemain := genericRemainRightStructure["title"].(string) + "---" + genericRemainRightStructure["num"].(string) + genericRemainRightStructure["unit"].(string)
		generic = "【" + genericUse + "】---【" + specialRemain + "】"
	} else {
		generic = "【" + genericUse + "】"
	}

	//国内通用流量详情
	var genericInfo string
	productInfos := ratableResourcePackages[0].(map[string]interface{})["productInfos"].([]interface{})
	for _, v := range productInfos {
		//判断不限量
		var name, use, remain, total string
		switch v.(map[string]interface{})["isInfiniteAmount"] {
		case "0": //限量
			name = v.(map[string]interface{})["title"].(string)
			use = v.(map[string]interface{})["leftTitle"].(string) + v.(map[string]interface{})["leftHighlight"].(string)
			remain = v.(map[string]interface{})["rightTitle"].(string) + v.(map[string]interface{})["rightHighlight"].(string)
			total = v.(map[string]interface{})["rightCommon"].(string)
			item := name + "\n" + "[" + use + "]" + "[" + remain + "]" + "[" + total + "]\n" + "--------------------------------------------------------------------------------\n"
			genericInfo += item

		case "1": //不限量
			name = v.(map[string]interface{})["title"].(string) + "【不限量】"
			use = v.(map[string]interface{})["infiniteTitle"].(string) + v.(map[string]interface{})["infiniteValue"].(string) + v.(map[string]interface{})["infiniteUnit"].(string)
			item := name + "\n" + "[" + use + "]" + "\n" + "--------------------------------------------------------------------------------\n"
			genericInfo += item
		}

	}

	//专用流量,判断是否为空!!!!!!!!!!!!!!!!!
	var special, specialUse string
	//specialUse := ratableResourcePackages[1].(map[string]interface{})["leftStructure"].(map[string]interface{})["num"].(string) + ratableResourcePackages[0].(map[string]interface{})["leftStructure"].(map[string]interface{})["unit"].(string) + "[" + ratableResourcePackages[0].(map[string]interface{})["leftStructure"].(map[string]interface{})["title"].(string) + "]\n"

	specialUseLeftStructure, ok := ratableResourcePackages[1].(map[string]interface{})["leftStructure"].(map[string]interface{})
	if ok {
		specialUse = specialUseLeftStructure["title"].(string) + specialUseLeftStructure["num"].(string) + specialUseLeftStructure["unit"].(string)
	}
	//
	//specialRemainRightStructure, ok := ratableResourcePackages[1].(map[string]interface{})["rightStructure"].(map[string]interface{})
	//if ok {
	//	specialRemain = specialRemainRightStructure["title"].(string) + specialRemainRightStructure["num"].(string) + specialRemainRightStructure["unit"].(string)
	//}
	special = "【" + specialUse + "】"

	//专用流量详情
	var specialInfo string
	specialProductInfos, ok := ratableResourcePackages[1].(map[string]interface{})["productInfos"].([]interface{})
	if ok {
		for _, v := range specialProductInfos {
			var name, use, remain, total string

			switch v.(map[string]interface{})["isInfiniteAmount"] {
			case "0": //限量
				name = v.(map[string]interface{})["title"].(string)
				use = v.(map[string]interface{})["leftTitle"].(string) + v.(map[string]interface{})["leftHighlight"].(string)
				remain = v.(map[string]interface{})["rightTitle"].(string) + v.(map[string]interface{})["rightHighlight"].(string)
				total = v.(map[string]interface{})["rightCommon"].(string)
				item := name + "\n" + "[" + use + "]" + "[" + remain + "]" + "[" + total + "]\n" + "--------------------------------------------------------------------------------\n"
				specialInfo += item
			case "1": //不限量
				name = v.(map[string]interface{})["title"].(string) + "【不限量】"
				use = v.(map[string]interface{})["infiniteTitle"].(string) + v.(map[string]interface{})["infiniteValue"].(string) + v.(map[string]interface{})["infiniteUnit"].(string)
				item := name + "\n" + "[" + use + "]" + "\n" + "--------------------------------------------------------------------------------\n"
				specialInfo += item

			}

		}
	}

	var result = map[string]interface{}{
		"ispType":     "telecom",
		"packageName": packageName,
		"mobile":      newMobile,
		"date":        date,
		"generic":     generic,
		"genericInfo": genericInfo,
		"special":     special,
		"specialInfo": specialInfo,
	}
	resultJson, err := json.Marshal(result)
	//fmt.Println(string(resultJson))
	return string(resultJson), err
}
