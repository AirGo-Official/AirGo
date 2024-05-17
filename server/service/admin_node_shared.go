package service

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/ppoonk/AirGo/constant"
	"github.com/ppoonk/AirGo/global"
	"github.com/ppoonk/AirGo/model"
	"github.com/ppoonk/AirGo/utils/encrypt_plugin"
	"github.com/ppoonk/AirGo/utils/net_plugin"
	"net/url"
	"strconv"
	"strings"
	"time"
)

func ParseVMessLink(link string) *model.Node {
	node := new(model.Node)
	node.Enabled = true
	node.Protocol = constant.NODE_PROTOCOL_VMESS
	if strings.ToLower(link[:8]) == "vmess://" {
		link = link[8:]
	} else {
		return nil
	}
	if len(link) == 0 {
		return nil
	}
	jsonStr := encrypt_plugin.SubBase64Decode(link)
	if jsonStr == "" {
		return nil
	}
	var mapResult map[string]interface{}
	err := json.Unmarshal([]byte(jsonStr), &mapResult)
	if err != nil {
		return nil
	}
	if version, ok := mapResult["v"]; ok {
		node.V = fmt.Sprintf("%v", version)
	}
	if ps, ok := mapResult["ps"]; ok {
		node.Remarks = fmt.Sprintf("%v", ps) //别名
	} else {
		return nil
	}
	if addr, ok := mapResult["add"]; ok {
		node.Address = fmt.Sprintf("%v", addr) //地址
	} else {
		return nil
	}
	if scy, ok := mapResult["scy"]; ok {
		node.Scy = fmt.Sprintf("%v", scy) //加密方式 auto,none,chacha20-poly1305,aes-128-gcm,zero
	} else {
		node.Scy = "auto"
	}
	if port, ok := mapResult["port"]; ok {
		value, err := strconv.ParseInt(fmt.Sprintf("%v", port), 10, 64)
		if err == nil {
			node.Port = value //端口
		} else {
			return nil
		}
	} else {
		return nil
	}

	if id, ok := mapResult["id"]; ok {
		node.UUID = fmt.Sprintf("%v", id) //uuid
	} else {
		return nil
	}
	if aid, ok := mapResult["aid"]; ok {
		if value, err := strconv.ParseInt(fmt.Sprintf("%v", aid), 10, 64); err == nil {
			node.Aid = value //额外id
		} else {
			return nil
		}
	} else {
		return nil
	}
	if net, ok := mapResult["net"]; ok {
		node.Network = fmt.Sprintf("%v", net) //传输协议
	} else {
		return nil
	}
	if type1, ok := mapResult["type"]; ok {
		node.Type = fmt.Sprintf("%v", type1)
	} else {
		return nil
	}

	//获取混淆
	if host, ok := mapResult["host"]; ok {
		node.Host = fmt.Sprintf("%v", host)
	} else {
		return nil
	}

	if path, ok := mapResult["path"]; ok {
		node.Path = fmt.Sprintf("%v", path)
	} else {
		return nil
	}
	if tls, ok := mapResult["tls"]; ok {
		node.Security = fmt.Sprintf("%v", tls)
	} else {
		return nil
	}
	if sni, ok := mapResult["sni"]; ok {
		node.Sni = fmt.Sprintf("%v", sni)
	}
	if alpn, ok := mapResult["alpn"]; ok {
		node.Alpn = fmt.Sprintf("%v", alpn)
	}
	return node
}

func ParseVLessLink(link string) *model.Node {
	u, err := url.Parse(link)
	if err != nil {
		return nil
	}
	if u.User == nil || u.Scheme != constant.NODE_PROTOCOL_VLESS {
		return nil
	}
	node := new(model.Node)
	node.Enabled = true
	node.Protocol = constant.NODE_PROTOCOL_VLESS

	//remarks
	node.Remarks = u.Fragment
	if node.Remarks == "" {
		node.Remarks = u.Host
	}
	//address
	node.Address = u.Hostname()
	//port
	node.Port, err = strconv.ParseInt(u.Port(), 10, 64)
	if err != nil {
		return nil
	}
	//uuid
	node.UUID = u.User.Username()

	//解析参数
	urlQuery := u.Query()
	if urlQuery.Get("flow") != "" {
		node.VlessFlow = urlQuery.Get("flow")
	}
	if urlQuery.Get("encryption") != "" {
		node.VlessEncryption = urlQuery.Get("encryption")
	}
	if urlQuery.Get("type") != "" {
		node.Network = urlQuery.Get("type")
	}
	if urlQuery.Get("headerType") != "" {
		node.Type = urlQuery.Get("headerType")
	}
	if urlQuery.Get("host") != "" {
		node.Host = urlQuery.Get("host")
	}
	if urlQuery.Get("path") != "" {
		node.Path = urlQuery.Get("path")
	}
	if urlQuery.Get("security") != "" {
		node.Security = urlQuery.Get("security")
	}
	if urlQuery.Get("sni") != "" {
		node.Sni = urlQuery.Get("sni")
	}
	if urlQuery.Get("fp") != "" {
		node.Fingerprint = urlQuery.Get("fp")
	}
	if urlQuery.Get("alpn") != "" {
		node.Alpn = urlQuery.Get("alpn")
	}
	if urlQuery.Get("pbk") != "" {
		node.PublicKey = urlQuery.Get("pbk")
	}
	if urlQuery.Get("sid") != "" {
		node.ShortId = urlQuery.Get("sid")
	}
	if urlQuery.Get("allowInsecure") != "" {
		node.AllowInsecure = true
	}
	return node
}

func ParseTrojanLink(link string) *model.Node {
	u, err := url.Parse(link)
	if err != nil {
		return nil
	}
	if u.User == nil || u.Scheme != constant.NODE_PROTOCOL_TROJAN {
		return nil
	}
	node := new(model.Node)
	node.Enabled = true
	node.Protocol = constant.NODE_PROTOCOL_TROJAN
	//remarks
	node.Remarks = u.Fragment
	if node.Remarks == "" {
		node.Remarks = u.Host
	}
	//address
	node.Address = u.Hostname()
	//port
	node.Port, err = strconv.ParseInt(u.Port(), 10, 64)
	if err != nil {
		return nil
	}
	//uuid
	node.UUID = u.User.Username()

	//解析参数
	urlQuery := u.Query()
	if urlQuery.Get("network") != "" {
		node.Network = urlQuery.Get("network")
	}
	if urlQuery.Get("type") != "" {
		node.Type = urlQuery.Get("type")
	}
	//获取混淆
	if urlQuery.Get("host") != "" {
		node.Host = urlQuery.Get("host")
	} else {
		return nil
	}
	if urlQuery.Get("path") != "" {
		node.Path = urlQuery.Get("path")
	}
	if urlQuery.Get("tls") != "" {
		node.Security = urlQuery.Get("tls")
	}
	if urlQuery.Get("sni") != "" {
		node.Sni = urlQuery.Get("sni")
	}
	if urlQuery.Get("alpn") != "" {
		node.Alpn = urlQuery.Get("alpn")
	}
	if urlQuery.Get("allowInsecure") != "" {
		node.AllowInsecure = true
	}
	return node
}

func ParseSSLink(link string) *model.Node {
	ss, err := url.Parse(link)
	var node model.Node
	if err != nil {
		global.Logrus.Error(err.Error())
		return nil
	}
	node.Protocol = constant.NODE_PROTOCOL_SHADOWSOCKS
	node.Remarks = ss.Fragment
	node.Address = ss.Hostname()
	node.Port, err = strconv.ParseInt(ss.Port(), 10, 64)

	p, _ := SubBase64Decode(ss.User.String())
	arr := strings.SplitN(p, ":", 2)
	node.Scy = arr[0]
	node.UUID = arr[1] //Passwd存到uuid字段

	return &node
}

func ParseHy2Link(link string) *model.Node {
	u, err := url.Parse(link)
	if err != nil {
		return nil
	}
	if u.User == nil || u.Scheme != "hy2" {
		return nil
	}
	node := new(model.Node)
	node.Enabled = true
	node.Protocol = constant.NODE_PROTOCOL_HYSTERIA2
	//remarks
	node.Remarks = u.Fragment
	if node.Remarks == "" {
		node.Remarks = u.Host
	}
	//address
	node.Address = u.Hostname()
	//port
	node.Port, err = strconv.ParseInt(u.Port(), 10, 64)
	if err != nil {
		return nil
	}
	//解析参数
	urlQuery := u.Query()
	//uuid
	node.UUID = u.User.Username()
	if urlQuery.Get("sni") != "" {
		node.Sni = urlQuery.Get("sni")
	}
	return node
}

func (n *AdminNode) ParseSubUrl(urlStr string) *[]model.Node {
	//去掉前后空格
	urlStr = strings.TrimSpace(urlStr)
	//订阅url
	if !strings.HasPrefix(urlStr, constant.NODE_PROTOCOL_VMESS) && !strings.HasPrefix(urlStr, constant.NODE_PROTOCOL_VLESS) && !strings.HasPrefix(urlStr, constant.NODE_PROTOCOL_TROJAN) && !strings.HasPrefix(urlStr, "hy2") {
		if _, err := url.ParseRequestURI(urlStr); err == nil {
			rsp, err := net_plugin.ClientWithDNS("223.6.6.6", 5*time.Second).Get(urlStr)
			if err != nil {
				return nil
			}
			defer rsp.Body.Close()
			subLink := net_plugin.ReadDate(rsp)
			if len(subLink) == 0 {
				return nil
			}
			urlStr = subLink
		}
	}
	// base64编码
	if urlStrBase64Decode, err := SubBase64Decode(urlStr); err == nil {
		urlStr = urlStrBase64Decode
	}
	list := strings.Fields(urlStr) //节点url数组
	var Nodes []model.Node
	for _, v := range list {
		data := ParseOne(v)
		if data == nil {
			continue
		}
		Nodes = append(Nodes, *data)
	}
	return &Nodes
}

// 解析一条节点,vmess vless trojan hysteria2
func ParseOne(link string) *model.Node {
	u, err := url.Parse(link)
	if err != nil {
		return nil
	}
	switch u.Scheme {
	case constant.NODE_PROTOCOL_VMESS:
		if obj := ParseVMessLink(link); obj != nil {
			return obj
		}
	case constant.NODE_PROTOCOL_VLESS:
		if obj := ParseVLessLink(link); obj != nil {
			return obj
		}
	case constant.NODE_PROTOCOL_TROJAN:
		if obj := ParseTrojanLink(link); obj != nil {
			return obj
		}
	case "ss":
		if obj := ParseSSLink(link); obj != nil {
			return obj
		}
	case "hy2":
		if obj := ParseHy2Link(link); obj != nil {
			return obj
		}
	}
	return nil
}

// 对节点base64格式进行解析
func SubBase64Decode(str string) (string, error) {
	i := len(str) % 4
	switch i {
	case 1:
		str = str[:len(str)-1]
	case 2:
		str += "=="
	case 3:
		str += "="
	}
	//str = strings.Split(str, "//")[1]
	var data []byte
	var err error
	if strings.Contains(str, "-") || strings.Contains(str, "_") {
		data, err = base64.URLEncoding.DecodeString(str)

	} else {
		data, err = base64.StdEncoding.DecodeString(str)
	}
	if err != nil {
		fmt.Println(err)
	}
	return string(data), err
}
