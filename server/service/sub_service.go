package service

import (
	"AirGo/global"
	"AirGo/model"
	"AirGo/utils/encrypt_plugin"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"gorm.io/gorm"
	"net/url"
	"strconv"
	"strings"

	//"gopkg.in/yaml.v2"
	"gopkg.in/yaml.v2"
)

// 获取订阅
func GetUserSub(url string, subType string) string {
	//查找用户
	fmt.Println(url, subType)
	var u model.User
	err := global.DB.Where("subscribe_url = ? and sub_status = 1 and d + u < t", url).First(&u).Error
	if err != nil {
		fmt.Println(err)
		return ""
	}
	//根据goodsID 查找具体的节点
	var goods model.Goods
	err = global.DB.Where("id = ?", u.SubscribeInfo.GoodsID).Preload("Nodes", func(db *gorm.DB) *gorm.DB { return db.Order("node_order") }).Find(&goods).Error
	// 计算剩余天数，流量
	//fmt.Println("根据goodsID 查找具体的节点", goods)
	expiredTime := u.SubscribeInfo.ExpiredAt.Format("2006-01-02")
	expiredBd1 := (float64(u.SubscribeInfo.T - u.SubscribeInfo.U - u.SubscribeInfo.D)) / 1024 / 1024 / 1024
	expiredBd2 := strconv.FormatFloat(expiredBd1, 'f', 2, 64)
	name := "到期时间:" + expiredTime + "  |  剩余流量:" + expiredBd2 + "GB"
	var firstSubNode = model.Node{
		Remarks:  name,
		Address:  global.Server.System.SubName,
		Port:     6666,
		Aid:      0,
		Network:  "ws",
		Enabled:  true,
		NodeType: "vmess",
	}
	//插入计算剩余天数，流量的第一条节点
	goods.Nodes = append(goods.Nodes, model.Node{})
	copy(goods.Nodes[1:], goods.Nodes[0:])
	goods.Nodes[0] = firstSubNode
	//再插入共享的节点
	nodeList, _, err := CommonSqlFind[model.NodeShared, string, []model.Node]("")
	if err == nil {
		for _, v := range nodeList {
			goods.Nodes = append(goods.Nodes, v)
		}
	}
	//fmt.Println("nodes:", goods.Nodes)
	//根据subType生成不同客户端订阅
	switch subType {
	case "v2ray":
		fmt.Println(1)
		return V2rayNGSubscribe(&goods.Nodes, u)
	case "clash":
		return ClashSubscribe(&goods.Nodes, u.UUID.String(), u.SubscribeInfo.Host)
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
		var newHost, newUUID string
		if v.IsSharedNode {
			newHost = v.Host
			newUUID = v.UUID

		} else {
			newUUID = user.UUID.String()
			if user.SubscribeInfo.Host == "" {
				newHost = v.Host
			} else {
				newHost = user.SubscribeInfo.Host
			}
		}
		switch v.NodeType {
		case "vmess":
			if res := V2rayNGVmess(v, newUUID, newHost); res != "" {
				subArr = append(subArr, res)
			}
			continue
		case "vless":
			if res := V2rayNGVlessTrojan(v, "vless", newUUID, newHost); res != "" {
				subArr = append(subArr, res)
			}
			continue
		case "trojan":
			if res := V2rayNGVlessTrojan(v, "trojan", newUUID, newHost); res != "" {
				subArr = append(subArr, res)
			}
			continue
		case "shadowsocks":
			if res := V2rayNGShadowsocks(v, newUUID, user); res != "" {
				subArr = append(subArr, res)
			}
			continue
		}
	}
	return base64.StdEncoding.EncodeToString([]byte(strings.Join(subArr, "\r\n")))
}

// clash 订阅
func ClashSubscribe(nodes *[]model.Node, uuid, host string) string {
	var proxiesArr []model.ClashProxy
	//所有节点名称数组
	var nameArr []string
	for _, v := range *nodes {
		//剔除禁用节点
		if !v.Enabled {
			continue
		}

		var newHost, newUUID string
		if v.IsSharedNode {
			newHost = v.Host
			newUUID = v.UUID

		} else {
			newUUID = uuid
			if host == "" {
				newHost = v.Host
			} else {
				newHost = host
			}
		}

		nameArr = append(nameArr, v.Remarks)

		proxy := ClashVmessVlessTrojan(v, newUUID, newHost)
		proxiesArr = append(proxiesArr, proxy)

	}
	var proxyGroup model.ClashProxyGroup
	proxyGroup.Name = global.Server.System.SubName
	proxyGroup.Type = "select"
	proxyGroup.Proxies = nameArr

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
	clashYaml.ProxyGroups = append(clashYaml.ProxyGroups, proxyGroup)
	clashYaml.Rules = append(clashYaml.Rules, "MATCH,"+global.Server.System.SubName)
	res, err := yaml.Marshal(clashYaml)
	if err != nil {
		global.Logrus.Error("yaml.Marshal error:", err)
		return ""
	}
	return string(res)

}

// generate v2rayNG vmess
func V2rayNGVmess(node model.Node, uuid, host string) string {
	// {"add":"AirGo","aid":"0","alpn":"h2,http/1.1","fp":"qq","host":"www.baidu.com","id":"e0d5fe65-a5d1-4b8a-8d40-ed92a6a35d8b","net":"ws","path":"/path","port":"6666","ps":"到期时间:2024-03-06  |  剩余流量:20.00GB","scy":"auto","sni":"www.baidu.com","tls":"tls","type":"","v":"2"}
	var vmess model.Vmess
	vmess.V = node.V
	vmess.Name = node.Remarks
	if node.EnableTransfer {
		vmess.Address = node.TransferAddress
		vmess.Port = strconv.FormatInt(node.TransferPort, 10)
	} else {
		vmess.Address = node.Address
		vmess.Port = strconv.FormatInt(node.Port, 10)
	}
	vmess.Uuid = uuid
	vmess.Host = host
	vmess.Aid = strconv.FormatInt(node.Aid, 10)
	vmess.Net = node.Network
	vmess.Disguisetype = node.Type //伪装类型
	vmess.Path = node.Path
	//传输层安全
	switch node.Security {
	case "tls":
		// alpn fp sni
		vmess.Tls = node.Security
		vmess.Alpn = node.Alpn
		vmess.Sni = node.Sni
	}

	vmessMarshal, err := json.Marshal(vmess)
	if err != nil {
		return ""
	}
	vmessStr := base64.StdEncoding.EncodeToString([]byte(vmessMarshal))
	return "vmess://" + vmessStr
}

// generate  v2rayng vless trojan
func V2rayNGVlessTrojan(node model.Node, scheme, uuid, host string) string {
	// vless例子 vless://d342d11e-d424-4583-b36e-524ab1f0afa7@1.6.1.4:443?path=%2F%3Fed%3D2048&security=reality&encryption=none&pbk=ppkk&host=v2.airgoo.link&fp=randomized&spx=ssxx&flow=xtls-rprx-vision-udp443&type=ws&sni=v2.airgoo.link&sid=ssdd#v2.airgoo.link
	// [scheme:][//[userinfo@]host][/]path[?query][#fragment]
	var vlessUrl url.URL
	vlessUrl.Scheme = scheme
	vlessUrl.User = url.UserPassword(uuid, "")
	vlessUrl.Host = node.Address + ":" + strconv.FormatInt(node.Port, 10)
	values := url.Values{}
	switch scheme {
	case "vless":
		values.Add("encryption", node.Scy)
		values.Add("type", node.Network)
		values.Add("flow", node.VlessFlow)
		switch node.Network {
		case "ws":
			values.Add("host", host)
			values.Add("path", node.Path)
		case "tcp":
			values.Add("headerType", node.Type)
			values.Add("host", host)
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
	//return vlessUrl.String()
	return strings.ReplaceAll(vlessUrl.String(), ":@", "@")
}

// generate  Clash vmess vless trojan
func ClashVmessVlessTrojan(v model.Node, uuid, host string) model.ClashProxy {
	var proxy model.ClashProxy
	switch v.NodeType {
	case "vmess":
		proxy.Type = "vmess"
		proxy.Uuid = uuid
		proxy.Alterid = strconv.FormatInt(v.Aid, 10)
		proxy.Cipher = "auto"
	case "vless":
		proxy.Type = "vless"
		proxy.Uuid = uuid
		proxy.Flow = v.VlessFlow
	case "trojan":
		proxy.Type = "trojan"
		proxy.Uuid = uuid
		proxy.Sni = v.Sni
	case "shadowsocks":
		proxy.Type = "ss"
		proxy.Cipher = v.Scy
		proxy.Password = v.UUID
	}
	if v.EnableTransfer {
		proxy.Server = v.TransferAddress
		proxy.Port = int(v.TransferPort)
	} else {
		proxy.Server = v.Address
		proxy.Port = int(v.Port)
	}
	proxy.Name = v.Remarks
	proxy.Udp = true
	proxy.Network = v.Network
	proxy.SkipCertVerify = v.AllowInsecure

	switch proxy.Network {
	case "ws":
		proxy.WsOpts.Path = v.Path
		proxy.WsOpts.Headers = make(map[string]string, 1)
		proxy.WsOpts.Headers["Host"] = host

	case "grpc":
		proxy.GrpcOpts.GrpcServiceName = "grpc"
	case "tcp":
	case "h2":
		proxy.H2Opts.Path = v.Path
		proxy.H2Opts.Host = append(proxy.H2Opts.Host, host)

	}
	switch v.Security {
	case "tls":
		proxy.Tls = true
		proxy.Servername = v.Sni
		proxy.ClientFingerprint = v.Fingerprint
		proxy.Alpn = append(proxy.Alpn, v.Alpn)
		proxy.ClientFingerprint = v.Fingerprint
	case "reality":
		proxy.Tls = true
		proxy.Servername = v.Sni
		proxy.RealityOpts.PublicKey = v.PublicKey
		proxy.RealityOpts.ShortID = v.ShortId
		proxy.ClientFingerprint = v.Fingerprint
		proxy.Alpn = append(proxy.Alpn, v.Alpn)
	}
	return proxy
}

func V2rayNGShadowsocks(n model.Node, uuid string, user model.User) string {
	if !n.IsSharedNode {
		if strings.HasPrefix(n.Scy, "2022") {
			uuid = user.Passwd
		}
	}
	var urlV url.URL
	urlV.Scheme = "ss"
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
		urlV.User = url.UserPassword(p, "")
	default:
		p1 := n.Scy
		p2 := user.UUID.String()
		p := base64.StdEncoding.EncodeToString([]byte(p1 + ":" + p2))
		urlV.User = url.UserPassword(p, "")
	}

	urlV.Host = n.Address + ":" + fmt.Sprintf("%d", n.Port)
	urlV.Fragment = n.Remarks
	return strings.ReplaceAll(urlV.String(), ":@", "@")
}
