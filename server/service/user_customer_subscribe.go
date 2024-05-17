package service

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"github.com/ppoonk/AirGo/constant"
	"github.com/ppoonk/AirGo/global"
	"github.com/ppoonk/AirGo/model"
	uuid "github.com/satori/go.uuid"
	"gopkg.in/ini.v1"
	"gopkg.in/yaml.v3"
)

func (c *CustomerService) GetSubscribe(uuidStr string, clientType string) (string, string) {
	var nodeArr []model.Node
	// 查找用户服务
	subUUID, err := uuid.FromString(uuidStr)
	if err != nil {
		return "", ""
	}
	cs, err := CustomerServiceSvc.FirstCustomerService(&model.CustomerService{SubUUID: subUUID})
	if err != nil {
		return "", ""
	}
	//根据goodsID 查找具体的节点
	var goods model.Goods
	err = global.DB.
		Where(&model.Goods{ID: cs.GoodsID}).
		Preload("Nodes", "enabled = 1 ORDER BY node_order ASC").
		Find(&goods).
		Error
	// 计算剩余天数，流量
	expiredTime := cs.ServiceEndAt.Format("2006-01-02 15:04:05")
	expiredBd1 := (float64(cs.TotalBandwidth - cs.UsedUp - cs.UsedDown)) / 1024 / 1024 / 1024
	expiredBd2 := strconv.FormatFloat(expiredBd1, 'f', 2, 64)
	// 处理header参数
	headerInfo := fmt.Sprintf("upload=%d;download=%d;total=%d;expire=%d",
		cs.UsedUp, cs.UsedDown, cs.TotalBandwidth, cs.ServiceEndAt.Unix())
	var firstNode, secondNode model.Node
	if len(goods.Nodes) > 0 {
		firstNode = goods.Nodes[0]
		firstNode.Remarks = "到期时间:" + expiredTime

		secondNode = goods.Nodes[0]
		secondNode.Remarks = "剩余流量:" + expiredBd2 + "GB"

	} else {
		firstNode = model.Node{
			Remarks:  "到期时间:" + expiredTime,
			Address:  global.Server.Subscribe.SubName,
			Port:     6666,
			Aid:      0,
			Network:  "ws",
			Enabled:  true,
			Protocol: "vmess",
			NodeType: constant.NODE_TYPE_NORMAL,
		}
		secondNode = model.Node{
			Remarks:  "剩余流量:" + expiredBd2 + "GB",
			Address:  global.Server.Subscribe.SubName,
			Port:     6666,
			Aid:      0,
			Network:  "ws",
			Enabled:  true,
			Protocol: "vmess",
			NodeType: constant.NODE_TYPE_NORMAL,
		}

	}

	// 判断订阅是否有效，服务是否有效
	if !cs.ServiceStatus {
		nodeArr = append(goods.Nodes, firstNode, secondNode)
		//fmt.Println("goto subHandler")
		goto subHandler
	}

	//插入计算剩余天数，流量
	goods.Nodes = append(goods.Nodes, model.Node{}, model.Node{})
	copy(goods.Nodes[2:], goods.Nodes[0:])
	goods.Nodes[0] = firstNode
	goods.Nodes[1] = secondNode
	//最后处理一些参数
	for k, _ := range goods.Nodes {
		switch goods.Nodes[k].NodeType {
		case constant.NODE_TYPE_NORMAL: // 替换uuid
			goods.Nodes[k].UUID = cs.SubUUID.String()
		case constant.NODE_TYPE_TRANSFER: //中转节点修改ip和端口,以及uuid
			goods.Nodes[k].UUID = cs.SubUUID.String()
			goods.Nodes[k].Address = goods.Nodes[k].TransferAddress
			goods.Nodes[k].Port = goods.Nodes[k].TransferPort
		case constant.NODE_TYPE_SHARED:

		}
		//如果 NodeType 不是 transfer，但是transfer_address transfer_port 均不为空，也修改地址和端口
		if goods.Nodes[k].TransferAddress != "" && goods.Nodes[k].TransferPort != 0 {
			goods.Nodes[k].Address = goods.Nodes[k].TransferAddress
			goods.Nodes[k].Port = goods.Nodes[k].TransferPort
		}
		nodeArr = append(nodeArr, goods.Nodes[k])
	}
subHandler:
	//根据clientType生成不同客户端订阅
	switch clientType {
	case "v2rayNG", "V2rayU", "V2Box":
		return v2rayNG(&nodeArr), headerInfo

	case "v2rayN", "NekoBox":
		return NekoBox(&nodeArr), headerInfo

	case "Clash":
		return ClashMeta(&nodeArr), headerInfo

	case "Shadowrocket":
		return Shadowrocket(&nodeArr), headerInfo

	case "Surge":
		return Surge(&nodeArr), headerInfo

	case "Quantumult":
		return Quantumult(&nodeArr), headerInfo

	default:
		return "", headerInfo
	}
}

func v2rayNG(nodes *[]model.Node) string {
	var nodeArr []string
	for _, v := range *nodes {
		switch v.Protocol {
		case constant.NODE_PROTOCOL_VMESS:
			if res := VmessUrl(v); res != "" {
				nodeArr = append(nodeArr, res)
			}
		case constant.NODE_PROTOCOL_VLESS, constant.NODE_PROTOCOL_TROJAN:
			if res := VlessTrojanHysteriaUrl(v); res != "" {
				nodeArr = append(nodeArr, res)
			}
		case constant.NODE_PROTOCOL_SHADOWSOCKS:
			if res := ShadowsocksUrl(v); res != "" {
				nodeArr = append(nodeArr, res)
			}
		default:
			continue
		}
	}
	return base64.StdEncoding.EncodeToString([]byte(strings.Join(nodeArr, "\r\n")))
}

func NekoBox(nodes *[]model.Node) string {
	var nodeArr []string
	for _, v := range *nodes {
		switch v.Protocol {
		case constant.NODE_PROTOCOL_VMESS:
			if res := VmessUrl(v); res != "" {
				nodeArr = append(nodeArr, res)
			}

		case constant.NODE_PROTOCOL_VLESS, constant.NODE_PROTOCOL_TROJAN, constant.NODE_PROTOCOL_HYSTERIA2:
			if res := VlessTrojanHysteriaUrl(v); res != "" {
				nodeArr = append(nodeArr, res)
			}

		case constant.NODE_PROTOCOL_SHADOWSOCKS:
			if res := ShadowsocksUrl(v); res != "" {
				nodeArr = append(nodeArr, res)
			}
		default:
			continue
		}

	}
	return base64.StdEncoding.EncodeToString([]byte(strings.Join(nodeArr, "\r\n")))

}

func ClashMeta(nodes *[]model.Node) string {
	var proxiesArr []model.ClashProxy
	//所有节点名称数组
	var nameArr []string
	for _, v := range *nodes {
		nameArr = append(nameArr, v.Remarks)
		var proxy model.ClashProxy
		proxy = ClashGenerate(v)
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
	res, err := yaml.Marshal(clashYaml)
	if err != nil {
		global.Logrus.Error("yaml.Marshal error:", err)
		return ""
	}
	return fmt.Sprintf("%s\n%s",
		string(res),
		strings.ReplaceAll(global.Server.Subscribe.ClashRule, "AirGo-PROXY", global.Server.Subscribe.SubName)) //分流组名字使用自定义的订阅名
}

func Shadowrocket(nodes *[]model.Node) string {
	var nodeArr []string
	nodeArr = append(nodeArr, "STATUS="+(*nodes)[0].Remarks+"|"+(*nodes)[1].Remarks)
	for _, v := range *nodes {

		switch v.Protocol {
		case constant.NODE_PROTOCOL_VMESS, constant.NODE_PROTOCOL_VLESS:
			if res := VmessUrlForShadowrocket(v); res != "" {
				nodeArr = append(nodeArr, res)
			}

		case constant.NODE_PROTOCOL_SHADOWSOCKS:

			if res := ShadowsocksUrl(v); res != "" {
				nodeArr = append(nodeArr, res)
			}
		case constant.NODE_PROTOCOL_TROJAN:
			if res := TrojanUrlForShadowrocket(v); res != "" {
				nodeArr = append(nodeArr, res)
			}
		case constant.NODE_PROTOCOL_HYSTERIA2:
			if res := Hy2UrlForShadowrocket(v); res != "" {
				nodeArr = append(nodeArr, res)
			}
		default:
			continue
		}
	}
	return base64.StdEncoding.EncodeToString([]byte(strings.Join(nodeArr, "\r\n")))
}

// Surge 客户端订阅
func Surge(nodes *[]model.Node) string {
	var nodeArr, proxyGroupProxy, proxyGroupAuto, proxyGroupFallback []string
	var subscribeInfo, proxyText string
	//
	proxyGroupProxy = append(proxyGroupProxy, "select")
	proxyGroupProxy = append(proxyGroupProxy, "auto")
	proxyGroupProxy = append(proxyGroupProxy, "fallback")
	//
	proxyGroupAuto = append(proxyGroupAuto, "url-test")
	proxyGroupAuto = append(proxyGroupAuto, "url=http://www.gstatic.com/generate_204")
	proxyGroupAuto = append(proxyGroupAuto, "interval=43200")
	//
	proxyGroupFallback = append(proxyGroupFallback, "fallback")
	proxyGroupFallback = append(proxyGroupFallback, "url=http://www.gstatic.com/generate_204")
	proxyGroupFallback = append(proxyGroupFallback, "interval=43200")

	for _, v := range *nodes {

		switch v.Protocol {
		case constant.NODE_PROTOCOL_VMESS: //VMESS协议
			var nodeItem []string
			//
			nodeItem = append(nodeItem, v.Remarks+"="+"vmess")
			nodeItem = append(nodeItem, v.Address)
			nodeItem = append(nodeItem, fmt.Sprintf("%d", v.Port))
			nodeItem = append(nodeItem, fmt.Sprintf("username=%s", v.UUID))
			nodeItem = append(nodeItem, "tfo=true")
			nodeItem = append(nodeItem, "udp-relay=true")
			//
			nodeItem = append(nodeItem, "vmess-aead=true")
			//tls
			if v.Security != "" && v.Security != "none" {
				nodeItem = append(nodeItem, "tls=true")
				sni := v.Address
				if v.Sni != "" {
					sni = v.Sni
				}
				nodeItem = append(nodeItem, fmt.Sprintf("sni=%s", sni))
				if v.AllowInsecure {
					nodeItem = append(nodeItem, "skip-cert-verify=true")
				}
			}
			//ws
			if v.Network == constant.NETWORK_WS {
				nodeItem = append(nodeItem, "ws=true")
				if v.Path != "" {
					nodeItem = append(nodeItem, fmt.Sprintf("ws-path=%s", v.Path))
				}
				if v.Host != "" {
					nodeItem = append(nodeItem, fmt.Sprintf("ws-headers=Host:%s", v.Host))
				}
			}
			nodeArr = append(nodeArr, strings.Join(nodeItem, ", "))
			proxyGroupProxy = append(proxyGroupProxy, v.Remarks)
			proxyGroupAuto = append(proxyGroupAuto, v.Remarks)
			proxyGroupFallback = append(proxyGroupFallback, v.Remarks)
		case constant.NODE_PROTOCOL_TROJAN: //Trojan协议
			var nodeItem []string
			nodeItem = append(nodeItem, v.Remarks+"="+"trojan")
			nodeItem = append(nodeItem, v.Address)
			nodeItem = append(nodeItem, fmt.Sprintf("%d", v.Port))
			nodeItem = append(nodeItem, fmt.Sprintf("password=%s", v.UUID))
			nodeItem = append(nodeItem, "tfo=true")
			//nodeItem = append(nodeItem, "udp-relay=true")
			//
			sni := v.Address
			if v.Sni != "" {
				sni = v.Sni
			}
			nodeItem = append(nodeItem, fmt.Sprintf("sni=%s", sni))
			if v.AllowInsecure {
				nodeItem = append(nodeItem, "skip-cert-verify=true")
			}
			nodeArr = append(nodeArr, strings.Join(nodeItem, ", "))
			proxyGroupProxy = append(proxyGroupProxy, v.Remarks)
			proxyGroupAuto = append(proxyGroupAuto, v.Remarks)
			proxyGroupFallback = append(proxyGroupFallback, v.Remarks)
		case constant.NODE_PROTOCOL_HYSTERIA2: //hy2协议
			var nodeItem []string
			nodeItem = append(nodeItem, v.Remarks+" = "+"hysteria2")
			nodeItem = append(nodeItem, v.Address)
			nodeItem = append(nodeItem, fmt.Sprintf("%d", v.Port))
			nodeItem = append(nodeItem, fmt.Sprintf("password=%s", v.UUID))
			nodeItem = append(nodeItem, fmt.Sprintf("download-bandwidth=%d", v.HyDownMbps))
			sni := v.Address
			if v.Sni != "" {
				sni = v.Sni
			}
			nodeItem = append(nodeItem, fmt.Sprintf("sni=%s", sni))
			//nodeItem = append(nodeItem, "tfo=true")
			nodeItem = append(nodeItem, "udp-relay=true")
			if v.AllowInsecure {
				nodeItem = append(nodeItem, "skip-cert-verify=true")
			}
			nodeArr = append(nodeArr, strings.Join(nodeItem, ", "))
			proxyGroupProxy = append(proxyGroupProxy, v.Remarks)
			proxyGroupAuto = append(proxyGroupAuto, v.Remarks)
			proxyGroupFallback = append(proxyGroupFallback, v.Remarks)

		case constant.NODE_PROTOCOL_SHADOWSOCKS:
			if strings.HasPrefix(v.Scy, "2022") {
				continue
			}
			var nodeItem []string
			nodeItem = append(nodeItem, v.Remarks+"="+"ss")
			nodeItem = append(nodeItem, v.Address)
			nodeItem = append(nodeItem, fmt.Sprintf("%d", v.Port))
			nodeItem = append(nodeItem, fmt.Sprintf("encrypt-method=%s", v.Scy))
			nodeItem = append(nodeItem, fmt.Sprintf("password=%s", SSPasswordEncodeToString(v)))
			nodeItem = append(nodeItem, "tfo=true")
			nodeItem = append(nodeItem, "udp-relay=true")
			//if v.AllowInsecure {
			//	nodeItem = append(nodeItem, "skip-cert-verify=true")
			//}

			nodeArr = append(nodeArr, strings.Join(nodeItem, ", "))
			proxyGroupProxy = append(proxyGroupProxy, v.Remarks)
			proxyGroupAuto = append(proxyGroupAuto, v.Remarks)
			proxyGroupFallback = append(proxyGroupFallback, v.Remarks)

		}
	}
	//
	subscribeInfo = fmt.Sprintf("title=%s, content=%s\\n%s, style=info", global.Server.Subscribe.SubName, (*nodes)[0].Remarks, (*nodes)[1].Remarks)
	//
	proxyText = strings.Join(nodeArr, "\n")
	//
	surgeConf := model.SurgeConf{
		General: model.General{
			Loglevel:                   "notify",
			Doh_server:                 "https://doh.pub/dns-query",
			Dns_server:                 []string{"223.5.5.5", "114.114.114.114"},
			Tun_excluded_routes:        []string{"0.0.0.0/8", "10.0.0.0/8", "100.64.0.0/10", "127.0.0.0/8", "169.254.0.0/16", "172.16.0.0/12", "192.0.0.0/24", "192.0.2.0/24", "192.168.0.0/16", "192.88.99.0/24", "198.51.100.0/24", "203.0.113.0/24", "224.0.0.0/4", "255.255.255.255/32"},
			Skip_proxy:                 []string{"localhost", "*.local", "injections.adguard.org", "local.adguard.org", "captive.apple.com", "guzzoni.apple.com", "0.0.0.0/8", "10.0.0.0/8", "17.0.0.0/8", "100.64.0.0/10", "127.0.0.0/8", "169.254.0.0/16", "172.16.0.0/12", "192.0.0.0/24", "192.0.2.0/24", "192.168.0.0/16", "192.88.99.0/24", "198.18.0.0/15", "198.51.100.0/24", "203.0.113.0/24", "224.0.0.0/4", "240.0.0.0/4", "255.255.255.255/32"},
			Wifi_assist:                true,
			Allow_wifi_access:          true,
			Wifi_access_http_port:      6152,
			Wifi_access_socks5_port:    6153,
			Http_listen:                "0.0.0.0:6152",
			Socks5_listen:              "0.0.0.0:6153",
			External_controller_access: "surgepasswd@0.0.0.0:6170",
			Replica:                    false,
			Tls_provider:               "openssl",
			Network_framework:          false,
			Exclude_simple_hostnames:   true,
			Ipv6:                       true,
			//Skip_server_cert_verify:    true,
			Test_timeout:      4,
			Proxy_test_url:    "http://www.gstatic.com/generate_204",
			Geoip_maxmind_url: "https://cdn.jsdelivr.net/gh/Hackl0us/GeoIP2-CN@release/Country.mmdb",
		},
		Replica: model.Replica{
			Hide_apple_request:       true,
			Hide_crashlytics_request: true,
			Use_keyword_filter:       false,
			Hide_udp:                 false,
		},
		Panel: model.Panel{
			SubscribeInfo: subscribeInfo,
		},
		Proxy: model.Proxy{
			ProxyText: proxyText,
		},
		ProxyGroup: model.ProxyGroup{
			Proxy:    proxyGroupProxy,
			Auto:     proxyGroupAuto,
			Fallback: proxyGroupFallback,
		},
		Rule: model.Rule{
			RuleText: strings.ReplaceAll(global.Server.Subscribe.SurgeRule, "AirGo-PROXY", global.Server.Subscribe.SubName), //分流组名字使用自定义的订阅名
		},
	}
	cfg := ini.Empty()
	err := cfg.ReflectFrom(&surgeConf)
	if err != nil {
		return ""
	}
	bf := bytes.NewBuffer([]byte{})
	_, err = cfg.WriteTo(bf)
	if err != nil {
		return ""
	}
	text := bf.String()
	text = strings.ReplaceAll(text, "\"\"\"", "")
	text = strings.ReplaceAll(text, "ProxyText = ", "")
	text = strings.ReplaceAll(text, "RuleText = ", "")
	text = strings.ReplaceAll(text, "AirGo-PROXY", fmt.Sprintf("%s", global.Server.Subscribe.SubName)) //分流组名字使用自定义的订阅名
	//fmt.Println("text:", text)
	return text

}

func Quantumult(nodes *[]model.Node) string {
	var nodeArr []string
	for _, v := range *nodes {
		switch v.Protocol {
		case constant.NODE_PROTOCOL_VMESS, constant.NODE_PROTOCOL_VLESS:
			var nodeItem []string
			nodeItem = append(nodeItem, fmt.Sprintf("vmess=%s:%d", v.Address, v.Port))
			nodeItem = append(nodeItem, fmt.Sprintf("password=%s", v.UUID))
			//qx method: none aes-128-gcm chacha20-ietf-poly1305
			//nodeItem = append(nodeItem, fmt.Sprintf("method=%s", "chacha20-ietf-poly1305"))
			if v.Protocol == constant.NODE_PROTOCOL_VMESS {
				nodeItem = append(nodeItem, fmt.Sprintf("method=%s", v.Scy))
			} else {
				nodeItem = append(nodeItem, fmt.Sprintf("method=%s", "none"))
			}

			switch v.Network {
			case constant.NETWORK_WS:
				nodeItem = append(nodeItem, fmt.Sprintf("obfs-uri=%s", v.Path))
				nodeItem = append(nodeItem, fmt.Sprintf("obfs-host=%s", v.Host))
				if v.Security == "tls" {
					nodeItem = append(nodeItem, fmt.Sprintf("obfs=wss"))
					if v.AllowInsecure {
						nodeItem = append(nodeItem, fmt.Sprintf("tls-verification=false"))
					} else {
						nodeItem = append(nodeItem, fmt.Sprintf("tls-verification=true"))
					}
				} else {
					nodeItem = append(nodeItem, fmt.Sprintf("obfs=ws"))
				}

			case constant.NETWORK_TCP:
				if v.Security == "tls" {
					nodeItem = append(nodeItem, fmt.Sprintf("obfs=over-tls"))
				}
				if v.AllowInsecure {
					nodeItem = append(nodeItem, fmt.Sprintf("tls-verification=false"))
				} else {
					nodeItem = append(nodeItem, fmt.Sprintf("tls-verification=true"))
				}

			}

			nodeItem = append(nodeItem, fmt.Sprintf("fast-open=false"))
			nodeItem = append(nodeItem, fmt.Sprintf("udp-relay=false"))
			nodeItem = append(nodeItem, fmt.Sprintf("tag=%s", v.Remarks))
			nodeArr = append(nodeArr, strings.Join(nodeItem, ", "))

		case constant.NODE_PROTOCOL_SHADOWSOCKS:
			if strings.HasPrefix(v.Scy, "2022") {
				continue
			}
			var nodeItem []string
			nodeItem = append(nodeItem, fmt.Sprintf("shadowsocks=%s:%d", v.Address, v.Port))
			nodeItem = append(nodeItem, fmt.Sprintf("method=%s", v.Scy))
			nodeItem = append(nodeItem, fmt.Sprintf("password=%s", SSPasswordEncodeToString(v)))
			switch v.Type {
			case "http":
				nodeItem = append(nodeItem, fmt.Sprintf("obfs=http"))
				nodeItem = append(nodeItem, fmt.Sprintf("obfs-host=%s", v.Host))
				nodeItem = append(nodeItem, fmt.Sprintf("obfs-uri==%s", v.Path))

			}
			nodeItem = append(nodeItem, fmt.Sprintf("fast-open=false"))
			nodeItem = append(nodeItem, fmt.Sprintf("udp-relay=false"))
			nodeItem = append(nodeItem, fmt.Sprintf("tag=%s", v.Remarks))
			nodeArr = append(nodeArr, strings.Join(nodeItem, ", "))

		case constant.NODE_PROTOCOL_TROJAN:
			var nodeItem []string
			nodeItem = append(nodeItem, fmt.Sprintf("trojan=%s:%d", v.Address, v.Port))
			nodeItem = append(nodeItem, fmt.Sprintf("password=%s", v.UUID))
			nodeItem = append(nodeItem, fmt.Sprintf("obfs=over-tls"))
			if v.AllowInsecure {
				nodeItem = append(nodeItem, fmt.Sprintf("tls-verification=false"))
			} else {
				nodeItem = append(nodeItem, fmt.Sprintf("tls-verification=true"))
			}
			nodeItem = append(nodeItem, fmt.Sprintf("fast-open=false"))
			nodeItem = append(nodeItem, fmt.Sprintf("udp-relay=false"))
			nodeItem = append(nodeItem, fmt.Sprintf("tag=%s", v.Remarks))
			nodeArr = append(nodeArr, strings.Join(nodeItem, ", "))

		}
	}
	return strings.Join(nodeArr, "\n")
}

func VmessUrl(node model.Node) string {
	var vmess = &model.Vmess{
		V:            "2",
		Name:         node.Remarks,
		Address:      node.Address,
		Port:         fmt.Sprintf("%d", node.Port),
		Uuid:         node.UUID,
		Aid:          fmt.Sprintf("%d", node.Aid),
		Net:          node.Network,
		Disguisetype: "",
		Host:         node.Host,
		Path:         node.Path,
		Tls:          "",
		Alpn:         "",
		Fp:           "",
		Sni:          "",
	}
	switch node.Network {
	case constant.NETWORK_WS:
		vmess.Disguisetype = node.Type
		vmess.Path = node.Path
	case constant.NETWORK_TCP:
		vmess.Disguisetype = node.Type
		vmess.Path = node.Path
	case constant.NETWORK_GRPC:
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

func VlessTrojanHysteriaUrl(node model.Node) string {
	var nodeUrl url.URL
	switch node.Protocol {
	case constant.NODE_PROTOCOL_VLESS:
		nodeUrl.Scheme = "vless"
	case constant.NODE_PROTOCOL_TROJAN:
		nodeUrl.Scheme = "trojan"
	case constant.NODE_PROTOCOL_HYSTERIA2:
		nodeUrl.Scheme = "hy2"
	}

	nodeUrl.User = url.UserPassword(node.UUID, "")
	nodeUrl.Host = node.Address + ":" + strconv.FormatInt(node.Port, 10)

	values := url.Values{}
	switch nodeUrl.Scheme {
	case "vless":
		values.Add("encryption", "none")
		values.Add("type", node.Network)

		switch node.VlessFlow {
		case "xtls-rprx-vision", "xtls-rprx-vision-udp443":
			values.Add("flow", node.VlessFlow)
		default:
		}

		switch node.Network {
		case "ws":
			values.Add("host", node.Host)
			values.Add("path", node.Path)
		case "tcp":
			values.Add("headerType", node.Type)
			values.Add("host", node.Host)
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
		default:

		}
		if node.AllowInsecure {
			values.Add("allowInsecure", "1")
		}
	case "hy2":
		if node.HyPorts != "" {
			values.Add("mport", node.HyPorts)
		}
		if node.HyObfs != "" {
			values.Add("obfs", node.HyObfs)
			values.Add("obfs-password", node.HyObfsPassword)
		}
		if node.Sni != "" {
			values.Add("sni", node.Sni)
		}
		if node.AllowInsecure {
			values.Add("insecure", "1")
		}
	case "trojan":
		if node.AllowInsecure {
			values.Add("allowInsecure", "1")
		}
	}

	nodeUrl.RawQuery = values.Encode()
	nodeUrl.Fragment = node.Remarks
	return strings.ReplaceAll(nodeUrl.String(), ":@", "@")

}

func ShadowsocksUrl(node model.Node) string {
	if node.NodeType == constant.NODE_TYPE_SHARED {
		var ss url.URL
		ss.Scheme = "ss"
		ss.User = url.UserPassword(base64.StdEncoding.EncodeToString([]byte(node.Scy+":"+node.UUID)), "")
		ss.Host = node.Address + ":" + fmt.Sprintf("%d", node.Port)
		ss.Fragment = node.Remarks
		return strings.ReplaceAll(ss.String(), ":@", "@")
	}

	var ss url.URL
	ss.Scheme = "ss"
	ss.User = url.UserPassword(SSPasswordEncodeToString(node), "")
	ss.Host = node.Address + ":" + fmt.Sprintf("%d", node.Port)
	ss.Fragment = node.Remarks
	return strings.ReplaceAll(ss.String(), ":@", "@")

}

func ClashGenerate(node model.Node) model.ClashProxy {
	var proxy model.ClashProxy
	switch node.Protocol {
	case constant.NODE_PROTOCOL_VMESS:
		proxy.Type = "vmess"
		proxy.Uuid = node.UUID
		proxy.Alterid = node.Aid
		proxy.Cipher = "auto"
	case constant.NODE_PROTOCOL_VLESS:
		proxy.Type = "vless"
		proxy.Uuid = node.UUID
		proxy.Flow = node.VlessFlow
	case constant.NODE_PROTOCOL_TROJAN:
		proxy.Type = "trojan"
		proxy.Uuid = node.UUID
		proxy.Sni = node.Sni
	case constant.NODE_PROTOCOL_HYSTERIA2:
		proxy.Type = "hysteria2"
		//proxy.Uuid = node.UUID
		proxy.Password = node.UUID
		proxy.Sni = node.Sni
		proxy.Ports = node.HyPorts
		proxy.HopInterval = 15 //给个默认值吧，用户自行修改
		proxy.Up = fmt.Sprintf("%d Mbps", node.HyUpMbps)
		proxy.Down = fmt.Sprintf("%d Mbps", node.HyDownMbps)
		proxy.Obfs = node.HyObfs
		proxy.ObfsPassword = node.HyObfsPassword

	case constant.NODE_PROTOCOL_SHADOWSOCKS:
		proxy.Type = "ss"
		proxy.Cipher = node.Scy
		proxy.Password = GetSSPassword(node)
	}
	proxy.Name = node.Remarks
	proxy.Server = node.Address
	proxy.Port = int(node.Port)
	proxy.Udp = true
	proxy.Network = node.Network //传输协议
	proxy.SkipCertVerify = node.AllowInsecure

	switch proxy.Network {
	case "ws":
		proxy.WsOpts = model.WsOpts{
			Path: node.Path,
			Headers: map[string]string{
				"Host": node.Host,
			},
		}
	case "grpc":
		proxy.GrpcOpts = model.GrpcOpts{
			GrpcServiceName: node.ServiceName,
		}
	case "tcp":
		if node.Type == "http" {
			proxy.Network = "http"
			proxy.HttpOpts = model.HttpOpts{
				Method: "GET",
				Path:   []string{node.Path},
				Headers: map[string]model.Connection{
					"Host":       []string{node.Host},
					"Connection": []string{"keep-alive"},
				},
			}
		}
	case "h2":
		proxy.H2Opts = model.H2Opts{
			Host: []string{node.Host},
			Path: node.Path,
		}
	}
	switch node.Security {
	case "tls":
		proxy.Tls = true
		proxy.Servername = node.Sni
		proxy.ClientFingerprint = node.Fingerprint
		proxy.Alpn = append(proxy.Alpn, node.Alpn)
		proxy.ClientFingerprint = node.Fingerprint
	case "reality":
		proxy.Tls = true
		proxy.Servername = node.Sni
		proxy.RealityOpts.PublicKey = node.PublicKey
		proxy.RealityOpts.ShortID = node.ShortId
		proxy.ClientFingerprint = node.Fingerprint
		proxy.Alpn = append(proxy.Alpn, node.Alpn)
	default:
		proxy.Tls = false
	}
	return proxy
}

func VmessUrlForShadowrocket(node model.Node) string {
	var nodeUrl url.URL
	var user string
	values := url.Values{}

	switch node.Protocol {
	case constant.NODE_PROTOCOL_VMESS:
		nodeUrl.Scheme = "vmess"
		values.Add("alterId", fmt.Sprintf("%d", node.Aid)) //vmess alterId
		user = fmt.Sprintf("%s:%s@%s:%d", node.Scy, node.UUID, node.Address, node.Port)
	case constant.NODE_PROTOCOL_VLESS:
		nodeUrl.Scheme = "vless"
		node.Scy = "auto"
		user = fmt.Sprintf("%s:%s@%s:%d", node.Scy, node.UUID, node.Address, node.Port)
	}

	user = base64.StdEncoding.EncodeToString([]byte(user))
	user = strings.ReplaceAll(user, "=", "")
	nodeUrl.User = url.UserPassword(user, "")

	//基础参数
	//values.Add("tfo", "1") //tcp快速打开
	//values.Add("mux", "1")             //多路复用
	values.Add("remark", node.Remarks) //节点名

	//network
	switch node.Network {
	case constant.NETWORK_WS:
		values.Add("obfs", "websocket")
		values.Add("path", node.Path)
		values.Add("obfsParam", node.Host)
	case constant.NETWORK_TCP:
		values.Add("obfs", "none")
	case constant.NETWORK_GRPC:
		values.Add("obfs", "grpc")
		serviceName := node.Address
		host := node.Address
		if node.Host != "" {
			host = node.Host
		}
		if node.ServiceName != "" {
			serviceName = node.ServiceName
		}
		values.Add("host", host)
		values.Add("path", serviceName)
	}
	//tls
	switch node.Security {
	case "tls":
		values.Add("tls", "1")
		sni := node.Address
		if node.Sni != "" {
			sni = node.Sni
		}
		values.Add("peer", sni)

	case "reality":
		values.Add("tls", "1")
		sni := node.Address
		if node.Sni != "" {
			sni = node.Sni
		}
		values.Add("peer", sni)
		values.Add("pbk", node.PublicKey)
		values.Add("sid", node.ShortId)

	}
	//
	switch node.VlessFlow {
	case "xtls-rprx-direct":
		values.Add("xtls", "1")
	case "xtls-rprx-vision", "xtls-rprx-vision-udp443":
		values.Add("xtls", "2")
	default:
	}
	//
	switch node.AllowInsecure {
	case true:
		values.Add("allowInsecure", "1")

	}
	nodeUrl.RawQuery = values.Encode()
	return strings.ReplaceAll(nodeUrl.String(), ":@", "")

}
func TrojanUrlForShadowrocket(node model.Node) string {
	var nodeUrl url.URL
	nodeUrl.Scheme = "trojan"
	nodeUrl.Fragment = node.Remarks
	nodeUrl.User = url.UserPassword(node.UUID, "")
	nodeUrl.Host = node.Address + ":" + strconv.FormatInt(node.Port, 10)
	values := url.Values{}
	//values.Add("tfo", "1") //tcp快速打开
	//values.Add("mux", "1") //多路复用

	sni := node.Address
	if node.Sni != "" {
		sni = node.Sni
	}
	values.Add("peer", sni)
	switch node.AllowInsecure {
	case true:
		values.Add("allowInsecure", "1")
	}
	nodeUrl.RawQuery = values.Encode()
	return strings.ReplaceAll(nodeUrl.String(), ":@", "")
}
func Hy2UrlForShadowrocket(node model.Node) string {
	var nodeUrl url.URL
	nodeUrl.Scheme = "hysteria2"
	nodeUrl.Fragment = node.Remarks
	nodeUrl.User = url.UserPassword(node.UUID, "")
	nodeUrl.Host = node.Address + ":" + strconv.FormatInt(node.Port, 10)

	values := url.Values{}
	//values.Add("tfo", "1") //tcp快速打开
	//values.Add("mux", "1") //多路复用

	if node.HyPorts != "" {
		values.Add("mport", node.HyPorts)
	}
	if node.HyObfs != "" {
		values.Add("obfs", node.HyObfs)
		values.Add("obfs-password", node.HyObfsPassword)
	}
	if node.Sni != "" {
		values.Add("sni", node.Sni)
	}
	if node.AllowInsecure {
		values.Add("insecure", "1")
	}

	// values.Add("obfs", "none")
	nodeUrl.RawQuery = values.Encode()
	return strings.ReplaceAll(nodeUrl.String(), ":@", "@")
}

// 格式：serverKey:userKey
// 参考：ZDZiNDE5YjUzNDFiNzhmOWNjMTA0YTU0ZWJmNDYzNTc=:NzQ2NGIxZWQtMDQ1MS00NzZhLWIxODItN2JiZTU2YmU=
func GetSSPassword(node model.Node) string {
	switch node.Scy {
	case "2022-blake3-aes-128-gcm":
		serverKey := base64.StdEncoding.EncodeToString([]byte(node.ServerKey[:16]))
		userKey := base64.StdEncoding.EncodeToString([]byte(node.UUID[:16]))
		return serverKey + ":" + userKey
	case "2022-blake3-aes-256-gcm", "2022-blake3-chacha20-poly1305":
		serverKey := base64.StdEncoding.EncodeToString([]byte(node.ServerKey))
		userKey := base64.StdEncoding.EncodeToString([]byte(node.UUID))
		return serverKey + ":" + userKey
	default:
		return node.UUID
	}
}

// 原始数据：2022-blake3-aes-256-gcm:ZDZiNDE5YjUzNDFiNzhmOWNjMTA0YTU0ZWJmNDYzNTc=:NzQ2NGIxZWQtMDQ1MS00NzZhLWIxODItN2JiZTU2YmU=
// 输出：MjAyMi1ibGFrZTMtYWVzLTI1Ni1nY206WkRaaU5ERTVZalV6TkRGaU56aG1PV05qTVRBMFlUVTBaV0ptTkRZek5UYz06TnpRMk5HSXhaV1F0TURRMU1TMDBOelpoTFdJeE9ESXROMkppWlRVMlltVT0=
func SSPasswordEncodeToString(node model.Node) string {
	p := base64.StdEncoding.EncodeToString([]byte(node.Scy + ":" + GetSSPassword(node)))
	return p
}
