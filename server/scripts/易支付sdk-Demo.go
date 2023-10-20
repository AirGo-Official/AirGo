package main

import (
	"bytes"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"time"
)

func main() {
	//根据商户号、商户秘钥 初始化API
	api := NewZFAPI("商户号xxxxx", "商户秘钥xxxxxx")

	//设置异步通知地址(付款成功后渠道主动通知)
	api.SetNotifyURL("http://xxxxxxxxxxxx")

	//设置同步跳转地址(付款成功后前段跳转展示)
	api.SetReturnURL("http://xxxxxxxxxxxx")

	//发起支付请求
	if payInfo, err := api.StartWechatOrder("0.1", `customData`); err != nil {
		fmt.Printf("发起支付遇到错误：%s\n", err.Error())
	} else {
		fmt.Printf("渠道唯一ID:%s 支付页URL:%s\n", payInfo.ID, payInfo.PayURL)
	}
}

// ZFAPI ZFAPI
type ZFAPI struct {
	mSecretKey  string
	mMerchantID string
	mNotifyURL  string
	mReturnURL  string
}

// NewZFAPI NewZFAPI
func NewZFAPI(merchantID, secretKey string) *ZFAPI {
	return &ZFAPI{
		mMerchantID: merchantID,
		mSecretKey:  secretKey,
	}
}

// SetNotifyURL SetNotifyURL
func (z *ZFAPI) SetNotifyURL(url string) {
	z.mNotifyURL = url
}

// SetReturnURL SetReturnURL
func (z *ZFAPI) SetReturnURL(url string) {
	z.mReturnURL = url
}

// PayOrderInfo PayOrderInfo
type PayOrderInfo struct {
	ID        string `json:"id"`     //渠道唯一ID
	PayURL    string `json:"payUrl"` //支付页URL
	Amount    string `json:"-"`      //订单金额
	OrderNo   string `json:"-"`      //订单商户号
	AttchData string `json:"-"`      //订单附加透传参数
}

// StartWechatOrder 发起微信付款(金额，附加透传参数)
func (z *ZFAPI) StartWechatOrder(amount, attchData string) (info PayOrderInfo, err error) {
	//生成自己的订单号
	now := time.Now()
	orderNo := fmt.Sprintf("%04d%02d%02d%02d02%d02%02d%04d", now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), now.Second(), 1000+(rand.Int31()%8999))
	return z.StartOrder(orderNo, "wechat", amount, attchData)
}

// StartAlipayOrder 发起支付宝付款(金额，附加透传参数)
func (z *ZFAPI) StartAlipayOrder(amount, attchData string) (info PayOrderInfo, err error) {
	//生成自己的订单号
	now := time.Now()
	orderNo := fmt.Sprintf("%04d%02d%02d%02d02%d02%02d%04d", now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), now.Second(), 1000+(rand.Int31()%8999))
	return z.StartOrder(orderNo, "alipay", amount, attchData)
}

// StartOrder zhifu.fm-API(商户订单号，支付类型，金额，附加透传参数)
func (z *ZFAPI) StartOrder(orderNo, payType, amount, attchData string) (info PayOrderInfo, err error) {

	//计算签名
	var buff bytes.Buffer
	buff.WriteString(z.mMerchantID)
	buff.WriteString(orderNo)
	buff.WriteString(amount)
	buff.WriteString(z.mNotifyURL)
	buff.WriteString(z.mSecretKey)
	sign := fmt.Sprintf("%x", md5.Sum(buff.Bytes()))

	//构造请求参数
	buff.Reset()
	buff.WriteString("sign=")
	buff.WriteString(sign)
	buff.WriteString("&amount=")
	buff.WriteString(amount)
	buff.WriteString("&orderNo=")
	buff.WriteString(orderNo)
	buff.WriteString("&payType=")
	buff.WriteString(payType)
	buff.WriteString("&merchantNum=")
	buff.WriteString(z.mMerchantID)
	buff.WriteString("&notifyUrl=")
	buff.WriteString(z.mNotifyURL)
	buff.WriteString("&attch=")
	buff.WriteString(attchData)
	buff.WriteString("&returnUrl=")
	buff.WriteString(z.mReturnURL)

	//调用渠道接口
	var resp *http.Response
	if resp, err = http.Post("http://zfapi.nnt.ltd/api/startOrder", "application/x-www-form-urlencoded", &buff); err != nil {
		return
	}

	//读取所有响应数据
	var data []byte
	if data, err = ioutil.ReadAll(resp.Body); err != nil {
		return
	}

	type Result struct {
		Code int          `json:"code"`
		Msg  string       `json:"msg"`
		Data PayOrderInfo `json:"data"`
	}

	//解析渠道返回的JSON
	var result Result
	if err = json.Unmarshal(data, &result); err != nil {
		return
	}

	//检查应答状态码
	if result.Code != 200 {
		err = fmt.Errorf("code:%d msg:%s", result.Code, result.Msg)
		return
	}

	//返回支付页相关信息
	info = result.Data
	info.Amount = amount
	info.OrderNo = orderNo
	info.AttchData = attchData
	return
}
