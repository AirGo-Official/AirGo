package service

import (
	"AirGo/global"
	"AirGo/model"
	"AirGo/utils/encrypt_plugin"
	"encoding/base64"
	"encoding/json"
	"fmt"
	uuid "github.com/satori/go.uuid"
	"gopkg.in/yaml.v2"
	"gorm.io/gorm"
	"net/url"
	"strconv"
	"strings"
)

func Show(data any) {
	b, _ := json.Marshal(data)
	fmt.Println(string(b))

}

// 获取订阅
func GetUserSub(url string, subType string) string {
	//查找用户
	var u model.User
	err := global.DB.Where("subscribe_url = ? and sub_status = 1 and d + u < t", url).First(&u).Error
	if err != nil {
		return ""
	}
	//根据goodsID 查找具体的节点
	var goods model.Goods
	err = global.DB.Where("id = ?", u.SubscribeInfo.GoodsID).Preload("Nodes", func(db *gorm.DB) *gorm.DB { return db.Order("node_order") }).Find(&goods).Error
	// 计算剩余天数，流量
	expiredTime := u.SubscribeInfo.ExpiredAt.Format("2006-01-02")
	expiredBd1 := (float64(u.SubscribeInfo.T - u.SubscribeInfo.U - u.SubscribeInfo.D)) / 1024 / 1024 / 1024
	expiredBd2 := strconv.FormatFloat(expiredBd1, 'f', 2, 64)

	var firstNode = model.Node{
		Remarks:  "到期时间:" + expiredTime,
		Address:  global.Server.Subscribe.SubName,
		Port:     6666,
		Aid:      0,
		Network:  "ws",
		Enabled:  true,
		NodeType: "vmess",
	}
	var secondNode = model.Node{
		Remarks:  "剩余流量:" + expiredBd2 + "GB",
		Address:  global.Server.Subscribe.SubName,
		Port:     6666,
		Aid:      0,
		Network:  "ws",
		Enabled:  true,
		NodeType: "vmess",
	}
	//插入计算剩余天数，流量
	goods.Nodes = append(goods.Nodes, model.Node{}, model.Node{})
	copy(goods.Nodes[2:], goods.Nodes[0:])
	goods.Nodes[0] = firstNode
	goods.Nodes[1] = secondNode
	//再插入共享的节点
	nodeList, _, err := CommonSqlFind[model.NodeShared, string, []model.Node]("")
	if err == nil {
		for _, v := range nodeList {
			v.Enabled = true //设为启用
			v.IsSharedNode = true
			goods.Nodes = append(goods.Nodes, v)
		}
	}
	//fmt.Println(goods.Nodes)
	//根据subType生成不同客户端订阅
	switch subType {
	case "v2ray":
		return V2rayNGSubscribe(&goods.Nodes, u)
	case "clash":
		return ClashSubscribe(&goods.Nodes, u)
	default:
		return V2rayNGSubscribe(&goods.Nodes, u)
	}
}

// v2rayNG 订阅
func V2rayNGSubscribe(nodes *[]model.Node, user model.User) string {
	var subArr []string
	for _, v := range *nodes {
		//剔除禁用节点
		if !v.Enabled {
			continue
		}
		switch v.NodeType {
		case "vmess":
			if v.IsSharedNode {
				uuid, _ := uuid.FromString(v.UUID)
				user = model.User{UUID: uuid, SubscribeInfo: model.SubscribeInfo{Host: v.Host}}
			}
			if res := V2rayNGVmess(v, user); res != "" {
				subArr = append(subArr, res)
			}
			continue
		case "vless", "trojan":
			if v.IsSharedNode {
				uuid, _ := uuid.FromString(v.UUID)
				user = model.User{UUID: uuid, SubscribeInfo: model.SubscribeInfo{Host: v.Host}}
			}
			if res := V2rayNGVlessTrojan(v, user); res != "" {
				subArr = append(subArr, res)
			}
			continue
		case "shadowsocks":
			if v.IsSharedNode {
				if res := V2rayNGShadowsocksShared(v); res != "" {
					subArr = append(subArr, res)
				}

			} else {
				if res := V2rayNGShadowsocks(v, user); res != "" {
					subArr = append(subArr, res)
				}
			}
			continue
		}
	}
	return base64.StdEncoding.EncodeToString([]byte(strings.Join(subArr, "\r\n")))
}

// v2rayNG vmess
func V2rayNGVmess(node model.Node, user model.User) string {
	var vmess = &model.Vmess{
		V:            "2",
		Name:         node.Remarks,
		Address:      node.Address,
		Port:         fmt.Sprintf("%d", node.Port),
		Uuid:         user.UUID.String(),
		Aid:          fmt.Sprintf("%d", node.Aid),
		Net:          node.Network,
		Disguisetype: "",
		Host:         user.SubscribeInfo.Host,
		Path:         node.Path,
		Tls:          "",
		Alpn:         "",
		Fp:           "",
		Sni:          "",
	}
	switch node.Network {
	case "ws":
		vmess.Disguisetype = node.Type
		vmess.Path = node.Path
	case "tcp":
		vmess.Disguisetype = node.Type
		vmess.Path = node.Path
	case "grpc":
		vmess.Disguisetype = node.GrpcMode
		vmess.Path = node.ServiceName
	}
	switch node.Security {
	case "tls":
		vmess.Tls = node.Security
		vmess.Sni = node.Sni
		vmess.Fp = node.Fingerprint
	}
	vmessMarshal, err := json.Marshal(vmess)
	if err != nil {
		return ""
	}
	vmessStr := base64.StdEncoding.EncodeToString([]byte(vmessMarshal))
	return "vmess://" + vmessStr
}

// v2rayNG vless trojan
func V2rayNGVlessTrojan(node model.Node, user model.User) string {
	var vlessUrl url.URL
	switch node.NodeType {
	case "vless":
		vlessUrl.Scheme = "vless"
	case "trojan":
		vlessUrl.Scheme = "trojan"
	}
	vlessUrl.User = url.UserPassword(user.UUID.String(), "")
	vlessUrl.Host = node.Address + ":" + strconv.FormatInt(node.Port, 10)
	values := url.Values{}
	switch vlessUrl.Scheme {
	case "vless":
		values.Add("encryption", "none")
		values.Add("type", node.Network)
		values.Add("flow", node.VlessFlow)
		switch node.Network {
		case "ws":
			values.Add("host", user.SubscribeInfo.Host)
			values.Add("path", node.Path)
		case "tcp":
			values.Add("headerType", node.Type)
			values.Add("host", user.SubscribeInfo.Host)
			values.Add("path", node.Path)
		case "kcp", "quic":
			values.Add("headerType", node.Type)
		case "grpc":
			values.Add("mode", node.GrpcMode)
			values.Add("serviceName", node.ServiceName)
		}
		switch node.Security {
		case "tls":
			values.Add("security", node.Security)
			values.Add("alpn", node.Alpn)
			values.Add("fp", node.Fingerprint)
			values.Add("sni", node.Sni)
		case "reality":
			values.Add("security", node.Security)
			values.Add("pbk", node.PublicKey)
			values.Add("fp", node.Fingerprint)
			values.Add("spx", node.SpiderX)
			values.Add("sni", node.Sni)
			values.Add("sid", node.ShortId)
		}
	case "trojan":
	}
	vlessUrl.RawQuery = values.Encode()
	vlessUrl.Fragment = node.Remarks
	return strings.ReplaceAll(vlessUrl.String(), ":@", "@")

}

// v2rayNG ss
func V2rayNGShadowsocks(n model.Node, user model.User) string {
	var ss url.URL
	ss.Scheme = "ss"
	switch strings.HasPrefix(n.Scy, "2022") {
	case true:
		p1 := n.Scy
		p2 := n.ServerKey
		if p2 == "" {
			p2 = encrypt_plugin.RandomString(32)
		}
		p3 := user.UUID.String()
		if p1 == "2022-blake3-aes-128-gcm" {
			p2 = p2[:16]
			p3 = p3[0:16]
		}
		p2 = base64.StdEncoding.EncodeToString([]byte(p2))
		p3 = base64.StdEncoding.EncodeToString([]byte(p3))
		p := base64.StdEncoding.EncodeToString([]byte(p1 + ":" + p2 + ":" + p3))
		ss.User = url.UserPassword(p, "")
	default:
		p1 := n.Scy
		p2 := user.UUID.String()
		p := base64.StdEncoding.EncodeToString([]byte(p1 + ":" + p2))
		ss.User = url.UserPassword(p, "")
	}

	ss.Host = n.Address + ":" + fmt.Sprintf("%d", n.Port)
	ss.Fragment = n.Remarks
	return strings.ReplaceAll(ss.String(), ":@", "@")

}

func V2rayNGShadowsocksShared(n model.Node) string {
	var ss url.URL
	ss.Scheme = "ss"
	ss.User = url.UserPassword(base64.StdEncoding.EncodeToString([]byte(n.Scy+":"+n.UUID)), "")
	ss.Host = n.Address + ":" + fmt.Sprintf("%d", n.Port)
	ss.Fragment = n.Remarks
	return strings.ReplaceAll(ss.String(), ":@", "@")
}

// clash 订阅
func ClashSubscribe(nodes *[]model.Node, user model.User) string {
	var proxiesArr []model.ClashProxy
	//所有节点名称数组
	var nameArr []string
	for _, v := range *nodes {
		//剔除禁用节点
		if !v.Enabled {
			continue
		}
		nameArr = append(nameArr, v.Remarks)
		var proxy model.ClashProxy
		if v.IsSharedNode { //共享节点
			proxy = ClashVmessVlessTrojanShared(v)
		} else {
			proxy = ClashVmessVlessTrojan(v, user)
		}

		proxiesArr = append(proxiesArr, proxy)
	}
	var proxyGroup1 = model.ClashProxyGroup{
		Name:    global.Server.Subscribe.SubName,
		Type:    "select",
		Proxies: nameArr,
	}
	var proxyGroup2 = model.ClashProxyGroup{
		Name:     "自动选择",
		Type:     "url-test",
		Proxies:  nameArr,
		Url:      "http://www.apple.com/library/test/success.html",
		Interval: 86400,
	}
	var proxyGroup3 = model.ClashProxyGroup{
		Name:     "故障转移",
		Type:     "fallback",
		Proxies:  nameArr,
		Url:      "http://www.apple.com/library/test/success.html",
		Interval: 7200,
	}

	var clashYaml model.ClashYaml
	clashYaml.Port = 7890
	clashYaml.SocksPort = 7891
	clashYaml.RedirPort = 7892
	clashYaml.AllowLan = false
	clashYaml.Mode = "rule"
	clashYaml.LogLevel = "silent"
	clashYaml.ExternalController = "0.0.0.0:9090"
	clashYaml.Secret = ""
	clashYaml.Proxies = proxiesArr
	clashYaml.ProxyGroups = append(clashYaml.ProxyGroups, proxyGroup1, proxyGroup2, proxyGroup3)
	clashYaml.Rules = []string{
		"DOMAIN-SUFFIX,local,DIRECT",
		"IP-CIDR,127.0.0.0/8,DIRECT",
		"IP-CIDR,172.16.0.0/12,DIRECT",
		"IP-CIDR,192.168.0.0/16,DIRECT",
		"IP-CIDR,10.0.0.0/8,DIRECT",
		"IP-CIDR,17.0.0.0/8,DIRECT",
		"IP-CIDR,100.64.0.0/10,DIRECT",
		"IP-CIDR,224.0.0.0/4,DIRECT",
		"IP-CIDR6,fe80::/10,DIRECT",
		"MATCH," + global.Server.Subscribe.SubName,
	}
	res, err := yaml.Marshal(clashYaml)
	if err != nil {
		global.Logrus.Error("yaml.Marshal error:", err)
		return ""
	}
	return string(res)

}

// Clash vmess vless trojan
func ClashVmessVlessTrojan(n model.Node, user model.User) model.ClashProxy {
	var proxy model.ClashProxy
	switch n.NodeType {
	case "vmess":
		proxy.Type = "vmess"
		proxy.Uuid = user.UUID.String()
		proxy.Alterid = fmt.Sprintf("%d", n.Aid)
		proxy.Cipher = "auto"
	case "vless":
		proxy.Type = "vless"
		proxy.Uuid = user.UUID.String()
		proxy.Flow = n.VlessFlow
	case "trojan":
		proxy.Type = "trojan"
		proxy.Uuid = user.UUID.String()
		proxy.Sni = n.Sni
	case "shadowsocks":
		proxy.Type = "ss"
		proxy.Cipher = n.Scy
		switch strings.HasPrefix(n.Scy, "2022") {
		case true:
			p1 := n.Scy
			p2 := n.ServerKey
			if p2 == "" {
				p2 = encrypt_plugin.RandomString(32)
			}
			p3 := user.UUID.String()
			if p1 == "2022-blake3-aes-128-gcm" {
				p2 = p2[:16]
				p3 = p3[0:16]
			}
			p2 = base64.StdEncoding.EncodeToString([]byte(p2))
			p3 = base64.StdEncoding.EncodeToString([]byte(p3))
			p := base64.StdEncoding.EncodeToString([]byte(p1 + ":" + p2 + ":" + p3))
			proxy.Password = p
		default:
			p1 := n.Scy
			p2 := user.UUID.String()
			p := base64.StdEncoding.EncodeToString([]byte(p1 + ":" + p2))
			proxy.Password = p
		}
	}
	if n.EnableTransfer {
		proxy.Server = n.TransferAddress
		proxy.Port = int(n.TransferPort)
	} else {
		proxy.Server = n.Address
		proxy.Port = int(n.Port)
	}
	proxy.Name = n.Remarks
	proxy.Udp = true
	proxy.Network = n.Network
	proxy.SkipCertVerify = n.AllowInsecure

	switch proxy.Network {
	case "ws":
		proxy.WsOpts.Path = n.Path
		proxy.WsOpts.Headers = make(map[string]string, 1)
		proxy.WsOpts.Headers["Host"] = user.SubscribeInfo.Host
	case "grpc":
		proxy.GrpcOpts.GrpcServiceName = "grpc"
	case "tcp":
	case "h2":
		proxy.H2Opts.Path = n.Path
		proxy.H2Opts.Host = append(proxy.H2Opts.Host, user.SubscribeInfo.Host)
	}
	switch n.Security {
	case "tls":
		proxy.Tls = true
		proxy.Servername = n.Sni
		proxy.ClientFingerprint = n.Fingerprint
		proxy.Alpn = append(proxy.Alpn, n.Alpn)
		proxy.ClientFingerprint = n.Fingerprint
	case "reality":
		proxy.Tls = true
		proxy.Servername = n.Sni
		proxy.RealityOpts.PublicKey = n.PublicKey
		proxy.RealityOpts.ShortID = n.ShortId
		proxy.ClientFingerprint = n.Fingerprint
		proxy.Alpn = append(proxy.Alpn, n.Alpn)
	}
	return proxy
}

func ClashVmessVlessTrojanShared(n model.Node) model.ClashProxy {
	var proxy model.ClashProxy
	switch n.NodeType {
	case "vmess":
		proxy.Type = "vmess"
		proxy.Uuid = n.UUID
		proxy.Alterid = fmt.Sprintf("%d", n.Aid)
		proxy.Cipher = "auto"
	case "vless":
		proxy.Type = "vless"
		proxy.Uuid = n.UUID
		proxy.Flow = n.VlessFlow
	case "trojan":
		proxy.Type = "trojan"
		proxy.Uuid = n.UUID
		proxy.Sni = n.Sni
	case "shadowsocks":
		proxy.Type = "ss"
		proxy.Cipher = n.Scy
		proxy.Password = n.UUID

	}
	proxy.Name = n.Remarks
	proxy.Server = n.Address
	proxy.Port = int(n.Port)
	proxy.Udp = true
	proxy.Network = n.Network
	proxy.SkipCertVerify = n.AllowInsecure
	switch proxy.Network {
	case "ws":
		proxy.WsOpts.Path = n.Path
		proxy.WsOpts.Headers = make(map[string]string, 1)
		proxy.WsOpts.Headers["Host"] = n.Host
	case "grpc":
		proxy.GrpcOpts.GrpcServiceName = "grpc"
	case "tcp":
	case "h2":
		proxy.H2Opts.Path = n.Path
		proxy.H2Opts.Host = append(proxy.H2Opts.Host, n.Host)
	}
	switch n.Security {
	case "tls":
		proxy.Tls = true
		proxy.Servername = n.Sni
		proxy.ClientFingerprint = n.Fingerprint
		proxy.Alpn = append(proxy.Alpn, n.Alpn)
		proxy.ClientFingerprint = n.Fingerprint
	case "reality":
		proxy.Tls = true
		proxy.Servername = n.Sni
		proxy.RealityOpts.PublicKey = n.PublicKey
		proxy.RealityOpts.ShortID = n.ShortId
		proxy.ClientFingerprint = n.Fingerprint
		proxy.Alpn = append(proxy.Alpn, n.Alpn)
	}
	return proxy
}
