package isp_plugin

import (
	"AirGo/model"
	"AirGo/utils/net_plugin"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

const (
	UnicomCodeUrl         = "https://m.client.10010.com/mobileService/sendRadomNum.htm"
	UnicomLoginUrl        = "https://m.client.10010.com/mobileService/radomLogin.htm"
	UnicomQueryTrafficUrl = "https://m.client.10010.com/servicequerybusiness/operationservice/queryOcsPackageFlowLeftContentRevisedInJune"
	Version               = "iphone_c@10.5"
	UnicomPublicKey       = "MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQDc+CZK9bBA9IU+gZUOc6FUGu7yO9WpTNB0PzmgFBh96Mg1WrovD1oqZ+eIF4LjvxKXGOdI79JRdve9NPhQo07+uqGQgE4imwNnRx7PFtCRryiIEcUoavuNtuRVoBAm6qdB0SrctgaqGfLgKvZHOnwTjyNqjBUxzMeQlEC2czEMSwIDAQAB"
)

func UnicomCode(isp *model.ISP) (string, error) {

	client := net_plugin.ClientWithDNS("114.114.114.114", 10*time.Second)

	formValues := url.Values{}
	formValues.Set("version", Version)
	formValues.Set("mobile", isp.UnicomConfig.UnicomMobile)
	formDataStr := formValues.Encode()
	formDataBytes := []byte(formDataStr)
	formBytesReader := bytes.NewReader(formDataBytes)

	req, err := http.NewRequest("POST", UnicomCodeUrl, formBytesReader)
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	//fmt.Println("请求参数：", req)
	resp, err := client.Do(req)
	defer resp.Body.Close()
	if err != nil {
		return "", err
	}
	out := net_plugin.ReadDate(resp)
	return out, err

}

func UnicomCodeLogin(password, mobile, appId string) (string, string, error) {
	client := net_plugin.ClientWithDNS("114.114.114.114", 10*time.Second)
	formValues := url.Values{}
	formValues.Set("version", Version)
	formValues.Set("mobile", mobile)
	formValues.Set("appId", appId)
	formValues.Set("password", password)
	formDataStr := formValues.Encode()
	formDataBytes := []byte(formDataStr)
	formBytesReader := bytes.NewReader(formDataBytes)

	req, _ := http.NewRequest("POST", UnicomLoginUrl, formBytesReader)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := client.Do(req)
	defer resp.Body.Close()
	if err != nil {
		return "", "", err
	}
	out := net_plugin.ReadDate(resp)
	var cookieArr []string
	cookies := resp.Cookies()
	for _, c := range cookies {
		cookieArr = append(cookieArr, c.Name+"="+c.Value)
	}
	cookie := strings.Join(cookieArr, ";")
	return out, cookie, err
}

func UnicomQueryTraffic(isp *model.ISP) (string, error) {
	client := net_plugin.ClientWithDNS("114.114.114.114", 10*time.Second)

	req, err := http.NewRequest("POST", UnicomQueryTrafficUrl, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("cookie", isp.UnicomConfig.Cookie)

	resp, err := client.Do(req)
	defer resp.Body.Close()
	if err != nil {
		return "", err
	}
	respData := net_plugin.ReadDate(resp)
	//fmt.Println("respData", respData)
	return UnicomQueryTrafficHandler(respData, isp.Mobile)

}

// 处理联通流量查询resp
// 如读取key为data里面嵌套的result：r["data"].(map[string]interface{})["result"]
// 如读取key为result的数组第四个数据：r["data"].(map[string]interface{})["result"].([]interface{})[3]
func UnicomQueryTrafficHandler(resp, mobile string) (string, error) {
	//判断999999，999998
	if resp == "999999" || resp == "999998" {
		fmt.Println("查询流量失败，重新登录")
		return "", errors.New("查询流量失败，重新登录")
	}
	//判断响应
	respMap := make(map[string]interface{})
	err := json.Unmarshal([]byte(resp), &respMap)
	if err != nil {
		fmt.Println("resp解析错误:", err)
		return "", err
	}
	//判断res.code
	if respMap["code"] != "0000" {
		desc, ok := respMap["desc"]
		if !ok {
			desc = "查询流量失败"
		}
		fmt.Println(desc)
		return "", errors.New(desc.(string))
	}
	//查询日期
	date := time.Now().Format("2006-01-02 15:04:05")
	//套餐名
	//fmt.Println("packageName", respMap["packageName"])
	packageName := respMap["packageName"].(string)
	//手机号
	newMobile := mobile[0:3] + "****" + mobile[7:]
	//已用流量
	sum := respMap["summary"].(map[string]interface{})["sum"].(string)
	freeFlow := respMap["summary"].(map[string]interface{})["freeFlow"].(string)
	//联通公免流量
	var mlResourcesStr string
	mlResourcesDetails, ok := respMap["MlResources"].([]interface{})[0].(map[string]interface{})["details"].([]interface{})
	if ok {
		if len(mlResourcesDetails) > 0 {
			for _, v := range mlResourcesDetails {
				item := v.(map[string]interface{})["feePolicyName"].(string) + " [已用:" + v.(map[string]interface{})["use"].(string) + "MB]\n"
				item += "--------------------------------------------------------------------------------\n"
				mlResourcesStr += item
			}
		}
	}
	//fmt.Println("联通公免流量ok")
	//非共享
	var unsharedList string
	unshared, ok := respMap["unshared"].([]interface{})
	if ok {
		unsharedDetails, ok := unshared[0].(map[string]interface{})["details"].([]interface{})
		if ok {
			if len(unsharedDetails) > 0 {
				for _, v := range unsharedDetails {
					name := v.(map[string]interface{})["feePolicyName"].(string)
					use := v.(map[string]interface{})["use"].(string)
					remain := v.(map[string]interface{})["remain"].(string)
					total := v.(map[string]interface{})["total"].(string)
					var item string
					totalFloat64, _ := strconv.ParseFloat(total, 64)
					switch totalFloat64 {
					case 0:
						item = name + "\n" + "[已用:" + use + "MB]" + "\n" + "--------------------------------------------------------------------------------\n"
					default:
						item = name + "\n" + "[已用:" + use + "MB]" + "[剩余:" + remain + "MB]" + "[总:" + total + "MB]\n" + "--------------------------------------------------------------------------------\n"

					}
					unsharedList += item
					//fmt.Println(item)
				}
			}
		}
	}
	//fmt.Println("非共享ok")
	//共享
	var sharedList string
	shared, ok := respMap["resources"].([]interface{})[0].(map[string]interface{})["details"].([]interface{})
	if ok {
		if len(shared) > 0 {

			for _, v := range shared {
				name := v.(map[string]interface{})["feePolicyName"].(string)
				use := v.(map[string]interface{})["use"].(string)
				remain := v.(map[string]interface{})["remain"].(string)
				total := v.(map[string]interface{})["total"].(string)
				var item string
				totalFloat64, _ := strconv.ParseFloat(total, 64)
				switch totalFloat64 {
				case 0:
					item = name + "\n" + "[已用:" + use + "MB]" + "\n"
				default:
					item = name + "\n" + "[已用:" + use + "MB]" + "[剩余:" + remain + "MB]" + "[总:" + total + "MB]\n"

				}
				//处理副卡
				viceCardlist, ok := v.(map[string]interface{})["viceCardlist"].([]interface{})
				if ok {
					for _, v1 := range viceCardlist {
						item = item + "[" + v1.(map[string]interface{})["usernumber"].(string) + "已用:" + v1.(map[string]interface{})["use"].(string) + "MB]\n"
					}
				}
				item += "--------------------------------------------------------------------------------\n"
				sharedList += item
				//fmt.Println(item)

			}

		}
	}

	//fmt.Println("out:", out)

	var result = map[string]interface{}{
		"ispType":     "unicom",
		"packageName": packageName,
		"mobile":      newMobile,
		"date":        date,
		"sum":         sum,
		"freeFlow":    freeFlow,
		"mlResources": mlResourcesStr,
		"flow":        unsharedList + sharedList,
		//"unshared":    unsharedList,
		//"shared":      sharedList,
	}
	resultJson, err := json.Marshal(result)
	//fmt.Println(string(resultJson))
	return string(resultJson), err

}
