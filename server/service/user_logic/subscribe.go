package user_logic

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/ppoonk/AirGo/constant"
	"github.com/ppoonk/AirGo/global"
	"github.com/ppoonk/AirGo/model"
	"github.com/ppoonk/AirGo/utils/encrypt_plugin"
	uuid "github.com/satori/go.uuid"
	"gopkg.in/ini.v1"
	"gopkg.in/yaml.v3"
	"net/url"
	"strconv"
	"strings"
)

func (c *CustomerService) GetSubscribe(uuidStr string, clientType string) (string, string) {
	var nodeArr []model.Node
	//查找用户
	//var u model.User
	//err := global.DB.Where("subscribe_url = ? and sub_status = 1 and d + u < t", url).First(&u).Error
	//if err != nil {
	//	return "", ""
	//}

	// 查找用户服务
	subUUID, err := uuid.FromString(uuidStr)
	if err != nil {
		return "", ""
	}
	cs, err := customerService.FirstCustomerService(&model.CustomerService{SubUUID: subUUID})
	if err != nil {
		return "", ""
	}
	//根据goodsID 查找具体的节点
	var goods model.Goods
	err = global.DB.
		Where(&model.Goods{ID: cs.GoodsID}).
		Preload("Nodes", "enabled = 1 AND ORDER BY node_order").
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
			Address:  global.Server.Website.SubName,
			Port:     6666,
			Aid:      0,
			Network:  "ws",
			Enabled:  true,
			Protocol: "vmess",
		}
		secondNode = model.Node{
			Remarks:  "剩余流量:" + expiredBd2 + "GB",
			Address:  global.Server.Website.SubName,
			Port:     6666,
			Aid:      0,
			Network:  "ws",
			Enabled:  true,
			Protocol: "vmess",
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
		case constant.NODE_TYPE_TRANSFER: //中转节点修改ip和端口
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
	case "v2rayNG", "V2rayU":
		return v2rayNG(&nodeArr), headerInfo

	case "NekoBox", "v2rayN":
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

	}
	return "", headerInfo
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

		case constant.NODE_PROTOCOL_VLESS, constant.NODE_PROTOCOL_TROJAN, constant.NODE_PROTOCOL_HYSTERIA:
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
		Name:    global.Server.Website.SubName,
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
	clashYaml.RuleProviders.CN = model.RuleProvidersItem{
		Behavior: "domain",
		Interval: 86400,
		Path:     "./rule-set/cn_domain.yaml",
		Type:     "http",
		Url:      "https://testingcf.jsdelivr.net/gh/MetaCubeX/meta-rules-dat@release/cn_domain.yaml",
	}

	clashYaml.RuleProviders.Proxy = model.RuleProvidersItem{
		Behavior: "domain",
		Interval: 86400,
		Path:     "./rule-set/proxy.yaml",
		Type:     "http",
		Url:      "https://testingcf.jsdelivr.net/gh/MetaCubeX/meta-rules-dat@release/proxy.yaml",
	}
	clashYaml.Rules = []string{
		"RULE-SET,cn,DIRECT",
		"RULE-SET,proxy," + global.Server.Website.SubName,
		"GEOSITE,category-ads-all,REJECT",

		"GEOSITE,private,DIRECT",
		"GEOSITE,onedrive,DIRECT",
		"GEOSITE,microsoft@cn,DIRECT",
		"GEOSITE,apple-cn,DIRECT",
		"GEOSITE,steam@cn,DIRECT",
		"GEOSITE,category-games@cn,DIRECT",
		"GEOSITE,cn,DIRECT",
		"GEOIP,CN,DIRECT",
		"GEOIP,private,DIRECT,no-resolve",

		"GEOSITE,youtube," + global.Server.Website.SubName,
		"GEOSITE,google," + global.Server.Website.SubName,
		"GEOSITE,twitter," + global.Server.Website.SubName,
		"GEOSITE,pixiv," + global.Server.Website.SubName,
		"GEOSITE,category-scholar-!cn," + global.Server.Website.SubName,
		"GEOSITE,biliintl," + global.Server.Website.SubName,
		"GEOSITE,geolocation-!cn," + global.Server.Website.SubName,
		"GEOIP,telegram," + global.Server.Website.SubName,
		"MATCH," + global.Server.Website.SubName,
	}
	res, err := yaml.Marshal(clashYaml)
	if err != nil {
		global.Logrus.Error("yaml.Marshal error:", err)
		return ""
	}
	return string(res)
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
		case constant.NODE_PROTOCOL_HYSTERIA:
			if res := Hy2UrlForShadowrocket(v); res != "" {
				nodeArr = append(nodeArr, res)
			}
		default:
			continue
		}
	}
	return base64.StdEncoding.EncodeToString([]byte(strings.Join(nodeArr, "\r\n")))
}

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
		case constant.NODE_PROTOCOL_VMESS:
			var nodeItem []string
			//
			nodeItem = append(nodeItem, v.Remarks+"="+"vmess")
			nodeItem = append(nodeItem, v.Address)
			nodeItem = append(nodeItem, fmt.Sprintf("%d", v.Port))
			nodeItem = append(nodeItem, fmt.Sprintf("username=%s", v.UUID))
			nodeItem = append(nodeItem, "tfo=true")
			//nodeItem = append(nodeItem, "udp-relay=true")
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

			}
			nodeArr = append(nodeArr, strings.Join(nodeItem, ", "))
			proxyGroupProxy = append(proxyGroupProxy, v.Remarks)
			proxyGroupAuto = append(proxyGroupAuto, v.Remarks)
			proxyGroupFallback = append(proxyGroupFallback, v.Remarks)
		case constant.NODE_PROTOCOL_TROJAN:
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
		case constant.NODE_PROTOCOL_HYSTERIA:
			var nodeItem []string
			nodeItem = append(nodeItem, v.Remarks+" = "+"hysteria2")
			nodeItem = append(nodeItem, v.Address)
			nodeItem = append(nodeItem, fmt.Sprintf("%d", v.Port))
			nodeItem = append(nodeItem, fmt.Sprintf("password=%s", v.UUID))
			//nodeItem = append(nodeItem, "tfo=true")
			//nodeItem = append(nodeItem, "udp-relay=true")
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
			nodeItem = append(nodeItem, fmt.Sprintf("password=%s", SSPasswordHandler(v)))
			nodeItem = append(nodeItem, "tfo=true")
			//nodeItem = append(nodeItem, "udp-relay=true")
			nodeArr = append(nodeArr, strings.Join(nodeItem, ", "))
			proxyGroupProxy = append(proxyGroupProxy, v.Remarks)
			proxyGroupAuto = append(proxyGroupAuto, v.Remarks)
			proxyGroupFallback = append(proxyGroupFallback, v.Remarks)

		}
	}
	//
	subscribeInfo = fmt.Sprintf("title=%s, content=%s\\n%s, style=info", global.Server.Website.SubName, (*nodes)[0].Remarks, (*nodes)[1].Remarks)
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
			Test_timeout:               4,
			Proxy_test_url:             "http://www.gstatic.com/generate_204",
			Geoip_maxmind_url:          "https://unpkg.zhimg.com/rulestatic@1.0.1/Country.mmdb",
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
			RuleText: DefaultSurgeRules,
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
	//fmt.Println("text:", text)
	return text

}

func Quantumult(nodes *[]model.Node) string {
	var nodeArr []string
	for _, v := range *nodes {
		switch v.Protocol {
		case constant.NODE_PROTOCOL_VMESS:
			var nodeItem []string
			nodeItem = append(nodeItem, fmt.Sprintf("vmess=%s:%d", v.Address, v.Port))
			nodeItem = append(nodeItem, fmt.Sprintf("method=%s", "chacha20-ietf-poly1305")) //surge不能为auto
			nodeItem = append(nodeItem, fmt.Sprintf("password=%s", v.UUID))

			switch v.Security {
			case "tls":
				if v.Network == constant.NETWORK_WS {
					nodeItem = append(nodeItem, fmt.Sprintf("obfs=over-tls"))
					//nodeItem = append(nodeItem, fmt.Sprintf("obfs-uri=%s", v.Path))
					nodeItem = append(nodeItem, fmt.Sprintf("obfs-host=%s", v.Host))
					if v.AllowInsecure {
						nodeItem = append(nodeItem, fmt.Sprintf("tls-verification=false"))
					} else {
						nodeItem = append(nodeItem, fmt.Sprintf("tls-verification=true"))
					}
				}

			default:
				if v.Network == constant.NETWORK_WS {
					nodeItem = append(nodeItem, fmt.Sprintf("obfs=ws"))
					nodeItem = append(nodeItem, fmt.Sprintf("obfs-uri=%s", v.Path))
					nodeItem = append(nodeItem, fmt.Sprintf("obfs-host=%s", v.Host))
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
			nodeItem = append(nodeItem, fmt.Sprintf("password=%s", SSPasswordHandler(v)))
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
	case constant.NODE_PROTOCOL_HYSTERIA:
		nodeUrl.Scheme = "hy2"
	}

	nodeUrl.User = url.UserPassword(node.UUID, "")
	nodeUrl.Host = node.Address + ":" + strconv.FormatInt(node.Port, 10)

	values := url.Values{}
	switch nodeUrl.Scheme {
	case "vless":
		values.Add("encryption", "none")
		values.Add("type", node.Network)
		values.Add("flow", node.VlessFlow)
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
			values.Add("security", "none")
		}
	case "hy2":
		values.Add("sni", node.Sni)

	case "trojan":
	}
	if node.AllowInsecure {
		values.Add("allowInsecure", "1")
	}

	nodeUrl.RawQuery = values.Encode()
	nodeUrl.Fragment = node.Remarks
	return strings.ReplaceAll(nodeUrl.String(), ":@", "@")

}

func ShadowsocksUrl(node model.Node) string {
	if node.NodeType != constant.NODE_TYPE_SHARED {
		var ss url.URL
		ss.Scheme = "ss"
		ss.User = url.UserPassword(base64.StdEncoding.EncodeToString([]byte(node.Scy+":"+node.UUID)), "")
		ss.Host = node.Address + ":" + fmt.Sprintf("%d", node.Port)
		ss.Fragment = node.Remarks
		return strings.ReplaceAll(ss.String(), ":@", "@")
	}

	var ss url.URL
	ss.Scheme = "ss"
	ss.User = url.UserPassword(SSPasswordHandler(node), "")
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
	case constant.NODE_PROTOCOL_HYSTERIA:
		proxy.Type = "hysteria2"
		proxy.Uuid = node.UUID
		proxy.Password = node.UUID
		proxy.Sni = node.Sni
	case constant.NODE_PROTOCOL_SHADOWSOCKS:
		proxy.Type = "ss"
		proxy.Cipher = node.Scy
		proxy.Password = node.UUID
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
	}
	return proxy
}

func VmessUrlForShadowrocket(node model.Node) string {
	var nodeUrl url.URL
	var user string

	switch node.Protocol {
	case constant.NODE_PROTOCOL_VMESS:
		nodeUrl.Scheme = "vmess"
		user = fmt.Sprintf("%s:%s@%s:%d", node.Scy, node.UUID, node.Address, node.Port)
	case constant.NODE_PROTOCOL_VLESS:
		nodeUrl.Scheme = "vless"
		node.Scy = "auto"
		user = fmt.Sprintf("%s:%s@%s:%d", node.Scy, node.UUID, node.Address, node.Port)
	}

	user = base64.StdEncoding.EncodeToString([]byte(user))
	user = strings.ReplaceAll(user, "=", "")
	nodeUrl.User = url.UserPassword(user, "")

	values := url.Values{}
	//基础参数
	//values.Add("tfo", "1") //tcp快速打开
	//values.Add("mux", "1")             //多路复用
	values.Add("remark", node.Remarks) //节点名

	switch node.Protocol {
	case constant.NODE_PROTOCOL_VMESS:
		values.Add("alterId", fmt.Sprintf("%d", node.Aid)) //vmess alterId
	}

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
	case "none", "":
	case "xtls-rprx-direct":
		values.Add("xtls", "1")
	case "xtls-rprx-vision", "xtls-rprx-vision-udp443":
		values.Add("xtls", "2")

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

	sni := node.Address
	if node.Sni != "" {
		sni = node.Sni
	}
	values.Add("peer", sni)

	switch node.AllowInsecure {
	case true:
		values.Add("insecure", "1")
	}

	values.Add("obfs", "none")
	nodeUrl.RawQuery = values.Encode()
	return strings.ReplaceAll(nodeUrl.String(), ":@", "@")
}

func SSPasswordHandler(node model.Node) string {
	switch strings.HasPrefix(node.Scy, "2022") {
	case true:
		p1 := node.Scy
		p2 := node.ServerKey
		if p2 == "" {
			p2 = encrypt_plugin.RandomString(32)
		}
		p3 := node.UUID
		if p1 == "2022-blake3-aes-128-gcm" {
			p2 = p2[:16]
			p3 = p3[0:16]
		}
		p2 = base64.StdEncoding.EncodeToString([]byte(p2))
		p3 = base64.StdEncoding.EncodeToString([]byte(p3))
		p := base64.StdEncoding.EncodeToString([]byte(p1 + ":" + p2 + ":" + p3))
		return p
	default:
		p1 := node.Scy
		p2 := node.UUID
		p := base64.StdEncoding.EncodeToString([]byte(p1 + ":" + p2))
		return p
	}

}

const DefaultSurgeRules = `
## 您可以在此处插入自定义规则
# 强制订阅域名直连
DOMAIN,192.168.0.92:16667,DIRECT
# Google 中国服务
DOMAIN-SUFFIX,services.googleapis.cn,Proxy
DOMAIN-SUFFIX,xn--ngstr-lra8j.com,Proxy
# Apple
DOMAIN,developer.apple.com,Proxy
DOMAIN-SUFFIX,digicert.com,Proxy
USER-AGENT,com.apple.trustd*,Proxy
DOMAIN-SUFFIX,apple-dns.net,Proxy
DOMAIN,testflight.apple.com,Proxy
DOMAIN,sandbox.itunes.apple.com,Proxy
DOMAIN,itunes.apple.com,Proxy
DOMAIN-SUFFIX,apps.apple.com,Proxy
DOMAIN-SUFFIX,blobstore.apple.com,Proxy
DOMAIN,cvws.icloud-content.com,Proxy
DOMAIN,safebrowsing.urlsec.qq.com,DIRECT
DOMAIN,safebrowsing.googleapis.com,DIRECT
USER-AGENT,com.apple.appstored*,DIRECT
USER-AGENT,AppStore*,DIRECT
DOMAIN-SUFFIX,mzstatic.com,DIRECT
DOMAIN-SUFFIX,itunes.apple.com,DIRECT
DOMAIN-SUFFIX,icloud.com,DIRECT
DOMAIN-SUFFIX,icloud-content.com,DIRECT
USER-AGENT,cloudd*,DIRECT
USER-AGENT,*com.apple.WebKit*,DIRECT
USER-AGENT,*com.apple.*,DIRECT
DOMAIN-SUFFIX,me.com,DIRECT
DOMAIN-SUFFIX,aaplimg.com,DIRECT
DOMAIN-SUFFIX,cdn20.com,DIRECT
DOMAIN-SUFFIX,cdn-apple.com,DIRECT
DOMAIN-SUFFIX,akadns.net,DIRECT
DOMAIN-SUFFIX,akamaiedge.net,DIRECT
DOMAIN-SUFFIX,edgekey.net,DIRECT
DOMAIN-SUFFIX,mwcloudcdn.com,DIRECT
DOMAIN-SUFFIX,mwcname.com,DIRECT
DOMAIN-SUFFIX,apple.com,DIRECT
DOMAIN-SUFFIX,apple-cloudkit.com,DIRECT
DOMAIN-SUFFIX,apple-mapkit.com,DIRECT
# 国内网站
USER-AGENT,MicroMessenger Client*,DIRECT
USER-AGENT,WeChat*,DIRECT
DOMAIN-SUFFIX,126.com,DIRECT
DOMAIN-SUFFIX,126.net,DIRECT
DOMAIN-SUFFIX,127.net,DIRECT
DOMAIN-SUFFIX,163.com,DIRECT
DOMAIN-SUFFIX,360buyimg.com,DIRECT
DOMAIN-SUFFIX,36kr.com,DIRECT
DOMAIN-SUFFIX,acfun.tv,DIRECT
DOMAIN-SUFFIX,air-matters.com,DIRECT
DOMAIN-SUFFIX,aixifan.com,DIRECT
DOMAIN-KEYWORD,alicdn,DIRECT
DOMAIN-KEYWORD,alipay,DIRECT
DOMAIN-KEYWORD,aliyun,DIRECT
DOMAIN-KEYWORD,taobao,DIRECT
DOMAIN-SUFFIX,amap.com,DIRECT
DOMAIN-SUFFIX,autonavi.com,DIRECT
DOMAIN-KEYWORD,baidu,DIRECT
DOMAIN-SUFFIX,bdimg.com,DIRECT
DOMAIN-SUFFIX,bdstatic.com,DIRECT
DOMAIN-SUFFIX,bilibili.com,DIRECT
DOMAIN-SUFFIX,bilivideo.com,DIRECT
DOMAIN-SUFFIX,caiyunapp.com,DIRECT
DOMAIN-SUFFIX,clouddn.com,DIRECT
DOMAIN-SUFFIX,cnbeta.com,DIRECT
DOMAIN-SUFFIX,cnbetacdn.com,DIRECT
DOMAIN-SUFFIX,cootekservice.com,DIRECT
DOMAIN-SUFFIX,csdn.net,DIRECT
DOMAIN-SUFFIX,ctrip.com,DIRECT
DOMAIN-SUFFIX,dgtle.com,DIRECT
DOMAIN-SUFFIX,dianping.com,DIRECT
DOMAIN-SUFFIX,douban.com,DIRECT
DOMAIN-SUFFIX,doubanio.com,DIRECT
DOMAIN-SUFFIX,duokan.com,DIRECT
DOMAIN-SUFFIX,easou.com,DIRECT
DOMAIN-SUFFIX,ele.me,DIRECT
DOMAIN-SUFFIX,feng.com,DIRECT
DOMAIN-SUFFIX,fir.im,DIRECT
DOMAIN-SUFFIX,frdic.com,DIRECT
DOMAIN-SUFFIX,g-cores.com,DIRECT
DOMAIN-SUFFIX,godic.net,DIRECT
DOMAIN-SUFFIX,gtimg.com,DIRECT
DOMAIN-SUFFIX,hongxiu.com,DIRECT
DOMAIN-SUFFIX,hxcdn.net,DIRECT
DOMAIN-SUFFIX,iciba.com,DIRECT
DOMAIN-SUFFIX,ifeng.com,DIRECT
DOMAIN-SUFFIX,ifengimg.com,DIRECT
DOMAIN-SUFFIX,ipip.net,DIRECT
DOMAIN-SUFFIX,iqiyi.com,DIRECT
DOMAIN-SUFFIX,jd.com,DIRECT
DOMAIN-SUFFIX,jianshu.com,DIRECT
DOMAIN-SUFFIX,knewone.com,DIRECT
DOMAIN-SUFFIX,le.com,DIRECT
DOMAIN-SUFFIX,lecloud.com,DIRECT
DOMAIN-SUFFIX,lemicp.com,DIRECT
DOMAIN-SUFFIX,licdn.com,DIRECT
DOMAIN-SUFFIX,luoo.net,DIRECT
DOMAIN-SUFFIX,meituan.com,DIRECT
DOMAIN-SUFFIX,meituan.net,DIRECT
DOMAIN-SUFFIX,mi.com,DIRECT
DOMAIN-SUFFIX,miaopai.com,DIRECT
DOMAIN-SUFFIX,microsoft.com,DIRECT
DOMAIN-SUFFIX,microsoftonline.com,DIRECT
DOMAIN-SUFFIX,miui.com,DIRECT
DOMAIN-SUFFIX,miwifi.com,DIRECT
DOMAIN-SUFFIX,mob.com,DIRECT
DOMAIN-SUFFIX,netease.com,DIRECT
DOMAIN-SUFFIX,office.com,DIRECT
DOMAIN-KEYWORD,officecdn,DIRECT
DOMAIN-SUFFIX,office365.com,DIRECT
DOMAIN-SUFFIX,oschina.net,DIRECT
DOMAIN-SUFFIX,ppsimg.com,DIRECT
DOMAIN-SUFFIX,pstatp.com,DIRECT
DOMAIN-SUFFIX,qcloud.com,DIRECT
DOMAIN-SUFFIX,qdaily.com,DIRECT
DOMAIN-SUFFIX,qdmm.com,DIRECT
DOMAIN-SUFFIX,qhimg.com,DIRECT
DOMAIN-SUFFIX,qhres.com,DIRECT
DOMAIN-SUFFIX,qidian.com,DIRECT
DOMAIN-SUFFIX,qihucdn.com,DIRECT
DOMAIN-SUFFIX,qiniu.com,DIRECT
DOMAIN-SUFFIX,qiniucdn.com,DIRECT
DOMAIN-SUFFIX,qiyipic.com,DIRECT
DOMAIN-SUFFIX,qq.com,DIRECT
DOMAIN-SUFFIX,qqurl.com,DIRECT
DOMAIN-SUFFIX,rarbg.to,DIRECT
DOMAIN-SUFFIX,ruguoapp.com,DIRECT
DOMAIN-SUFFIX,segmentfault.com,DIRECT
DOMAIN-SUFFIX,sinaapp.com,DIRECT
DOMAIN-SUFFIX,smzdm.com,DIRECT
DOMAIN-SUFFIX,snapdrop.net,DIRECT
DOMAIN-SUFFIX,sogou.com,DIRECT
DOMAIN-SUFFIX,sogoucdn.com,DIRECT
DOMAIN-SUFFIX,sohu.com,DIRECT
DOMAIN-SUFFIX,soku.com,DIRECT
DOMAIN-SUFFIX,speedtest.net,DIRECT
DOMAIN-SUFFIX,sspai.com,DIRECT
DOMAIN-SUFFIX,suning.com,DIRECT
DOMAIN-SUFFIX,taobao.com,DIRECT
DOMAIN-SUFFIX,tencent.com,DIRECT
DOMAIN-SUFFIX,tenpay.com,DIRECT
DOMAIN-SUFFIX,tianyancha.com,DIRECT
DOMAIN-KEYWORD,.tmall.com,DIRECT
DOMAIN-SUFFIX,tudou.com,DIRECT
DOMAIN-SUFFIX,umetrip.com,DIRECT
DOMAIN-SUFFIX,upaiyun.com,DIRECT
DOMAIN-SUFFIX,upyun.com,DIRECT
DOMAIN-SUFFIX,veryzhun.com,DIRECT
DOMAIN-SUFFIX,weather.com,DIRECT
DOMAIN-SUFFIX,weibo.com,DIRECT
DOMAIN-SUFFIX,xiami.com,DIRECT
DOMAIN-SUFFIX,xiami.net,DIRECT
DOMAIN-SUFFIX,xiaomicp.com,DIRECT
DOMAIN-SUFFIX,ximalaya.com,DIRECT
DOMAIN-SUFFIX,xmcdn.com,DIRECT
DOMAIN-SUFFIX,xunlei.com,DIRECT
DOMAIN-SUFFIX,yhd.com,DIRECT
DOMAIN-SUFFIX,yihaodianimg.com,DIRECT
DOMAIN-SUFFIX,yinxiang.com,DIRECT
DOMAIN-SUFFIX,ykimg.com,DIRECT
DOMAIN-SUFFIX,youdao.com,DIRECT
DOMAIN-SUFFIX,youku.com,DIRECT
DOMAIN-SUFFIX,zealer.com,DIRECT
DOMAIN-SUFFIX,zhihu.com,DIRECT
DOMAIN-SUFFIX,zhimg.com,DIRECT
DOMAIN-SUFFIX,zimuzu.tv,DIRECT
DOMAIN-SUFFIX,zoho.com,DIRECT

# 常见广告
DOMAIN-KEYWORD,admarvel,REJECT-TINYGIF
DOMAIN-KEYWORD,admaster,REJECT-TINYGIF
DOMAIN-KEYWORD,adsage,REJECT-TINYGIF
DOMAIN-KEYWORD,adsmogo,REJECT-TINYGIF
DOMAIN-KEYWORD,adsrvmedia,REJECT-TINYGIF
DOMAIN-KEYWORD,adwords,REJECT-TINYGIF
DOMAIN-KEYWORD,adservice,REJECT-TINYGIF
DOMAIN-SUFFIX,appsflyer.com,REJECT-TINYGIF
DOMAIN-KEYWORD,domob,REJECT-TINYGIF
DOMAIN-SUFFIX,doubleclick.net,REJECT-TINYGIF
DOMAIN-KEYWORD,duomeng,REJECT-TINYGIF
DOMAIN-KEYWORD,dwtrack,REJECT-TINYGIF
DOMAIN-KEYWORD,guanggao,REJECT-TINYGIF
DOMAIN-KEYWORD,lianmeng,REJECT-TINYGIF
DOMAIN-SUFFIX,mmstat.com,REJECT-TINYGIF
DOMAIN-KEYWORD,mopub,REJECT-TINYGIF
DOMAIN-KEYWORD,omgmta,REJECT-TINYGIF
DOMAIN-KEYWORD,openx,REJECT-TINYGIF
DOMAIN-KEYWORD,partnerad,REJECT-TINYGIF
DOMAIN-KEYWORD,pingfore,REJECT-TINYGIF
DOMAIN-KEYWORD,supersonicads,REJECT-TINYGIF
DOMAIN-KEYWORD,uedas,REJECT-TINYGIF
DOMAIN-KEYWORD,umeng,REJECT-TINYGIF
DOMAIN-KEYWORD,usage,REJECT-TINYGIF
DOMAIN-SUFFIX,vungle.com,REJECT-TINYGIF
DOMAIN-KEYWORD,wlmonitor,REJECT-TINYGIF
DOMAIN-KEYWORD,zjtoolbar,REJECT-TINYGIF

## 抗 DNS 污染
DOMAIN-KEYWORD,amazon,Proxy
DOMAIN-KEYWORD,google,Proxy
DOMAIN-KEYWORD,gmail,Proxy
DOMAIN-KEYWORD,youtube,Proxy
DOMAIN-KEYWORD,facebook,Proxy
DOMAIN-SUFFIX,fb.me,Proxy
DOMAIN-SUFFIX,fbcdn.net,Proxy
DOMAIN-KEYWORD,twitter,Proxy
DOMAIN-KEYWORD,instagram,Proxy
DOMAIN-KEYWORD,dropbox,Proxy
DOMAIN-SUFFIX,twimg.com,Proxy
DOMAIN-KEYWORD,blogspot,Proxy
DOMAIN-SUFFIX,youtu.be,Proxy

## 常见国外域名列表
DOMAIN-SUFFIX,9to5mac.com,Proxy
DOMAIN-SUFFIX,abpchina.org,Proxy
DOMAIN-SUFFIX,adblockplus.org,Proxy
DOMAIN-SUFFIX,adobe.com,Proxy
DOMAIN-SUFFIX,akamaized.net,Proxy
DOMAIN-SUFFIX,alfredapp.com,Proxy
DOMAIN-SUFFIX,amplitude.com,Proxy
DOMAIN-SUFFIX,ampproject.org,Proxy
DOMAIN-SUFFIX,android.com,Proxy
DOMAIN-SUFFIX,angularjs.org,Proxy
DOMAIN-SUFFIX,aolcdn.com,Proxy
DOMAIN-SUFFIX,apkpure.com,Proxy
DOMAIN-SUFFIX,appledaily.com,Proxy
DOMAIN-SUFFIX,appshopper.com,Proxy
DOMAIN-SUFFIX,appspot.com,Proxy
DOMAIN-SUFFIX,arcgis.com,Proxy
DOMAIN-SUFFIX,archive.org,Proxy
DOMAIN-SUFFIX,armorgames.com,Proxy
DOMAIN-SUFFIX,aspnetcdn.com,Proxy
DOMAIN-SUFFIX,att.com,Proxy
DOMAIN-SUFFIX,awsstatic.com,Proxy
DOMAIN-SUFFIX,azureedge.net,Proxy
DOMAIN-SUFFIX,azurewebsites.net,Proxy
DOMAIN-SUFFIX,bing.com,Proxy
DOMAIN-SUFFIX,bintray.com,Proxy
DOMAIN-SUFFIX,bit.com,Proxy
DOMAIN-SUFFIX,bit.ly,Proxy
DOMAIN-SUFFIX,bitbucket.org,Proxy
DOMAIN-SUFFIX,bjango.com,Proxy
DOMAIN-SUFFIX,bkrtx.com,Proxy
DOMAIN-SUFFIX,blog.com,Proxy
DOMAIN-SUFFIX,blogcdn.com,Proxy
DOMAIN-SUFFIX,blogger.com,Proxy
DOMAIN-SUFFIX,blogsmithmedia.com,Proxy
DOMAIN-SUFFIX,blogspot.com,Proxy
DOMAIN-SUFFIX,blogspot.hk,Proxy
DOMAIN-SUFFIX,bloomberg.com,Proxy
DOMAIN-SUFFIX,box.com,Proxy
DOMAIN-SUFFIX,box.net,Proxy
DOMAIN-SUFFIX,cachefly.net,Proxy
DOMAIN-SUFFIX,chromium.org,Proxy
DOMAIN-SUFFIX,cl.ly,Proxy
DOMAIN-SUFFIX,cloudflare.com,Proxy
DOMAIN-SUFFIX,cloudfront.net,Proxy
DOMAIN-SUFFIX,cloudmagic.com,Proxy
DOMAIN-SUFFIX,cmail19.com,Proxy
DOMAIN-SUFFIX,cnet.com,Proxy
DOMAIN-SUFFIX,cocoapods.org,Proxy
DOMAIN-SUFFIX,comodoca.com,Proxy
DOMAIN-SUFFIX,crashlytics.com,Proxy
DOMAIN-SUFFIX,culturedcode.com,Proxy
DOMAIN-SUFFIX,d.pr,Proxy
DOMAIN-SUFFIX,danilo.to,Proxy
DOMAIN-SUFFIX,dayone.me,Proxy
DOMAIN-SUFFIX,db.tt,Proxy
DOMAIN-SUFFIX,deskconnect.com,Proxy
DOMAIN-SUFFIX,disq.us,Proxy
DOMAIN-SUFFIX,disqus.com,Proxy
DOMAIN-SUFFIX,disquscdn.com,Proxy
DOMAIN-SUFFIX,dnsimple.com,Proxy
DOMAIN-SUFFIX,docker.com,Proxy
DOMAIN-SUFFIX,dribbble.com,Proxy
DOMAIN-SUFFIX,droplr.com,Proxy
DOMAIN-SUFFIX,duckduckgo.com,Proxy
DOMAIN-SUFFIX,dueapp.com,Proxy
DOMAIN-SUFFIX,dytt8.net,Proxy
DOMAIN-SUFFIX,edgecastcdn.net,Proxy
DOMAIN-SUFFIX,edgekey.net,Proxy
DOMAIN-SUFFIX,edgesuite.net,Proxy
DOMAIN-SUFFIX,engadget.com,Proxy
DOMAIN-SUFFIX,entrust.net,Proxy
DOMAIN-SUFFIX,eurekavpt.com,Proxy
DOMAIN-SUFFIX,evernote.com,Proxy
DOMAIN-SUFFIX,fabric.io,Proxy
DOMAIN-SUFFIX,fast.com,Proxy
DOMAIN-SUFFIX,fastly.net,Proxy
DOMAIN-SUFFIX,fc2.com,Proxy
DOMAIN-SUFFIX,feedburner.com,Proxy
DOMAIN-SUFFIX,feedly.com,Proxy
DOMAIN-SUFFIX,feedsportal.com,Proxy
DOMAIN-SUFFIX,fiftythree.com,Proxy
DOMAIN-SUFFIX,firebaseio.com,Proxy
DOMAIN-SUFFIX,flexibits.com,Proxy
DOMAIN-SUFFIX,flickr.com,Proxy
DOMAIN-SUFFIX,flipboard.com,Proxy
DOMAIN-SUFFIX,g.co,Proxy
DOMAIN-SUFFIX,gabia.net,Proxy
DOMAIN-SUFFIX,geni.us,Proxy
DOMAIN-SUFFIX,gfx.ms,Proxy
DOMAIN-SUFFIX,ggpht.com,Proxy
DOMAIN-SUFFIX,ghostnoteapp.com,Proxy
DOMAIN-SUFFIX,git.io,Proxy
DOMAIN-KEYWORD,github,Proxy
DOMAIN-SUFFIX,globalsign.com,Proxy
DOMAIN-SUFFIX,gmodules.com,Proxy
DOMAIN-SUFFIX,godaddy.com,Proxy
DOMAIN-SUFFIX,golang.org,Proxy
DOMAIN-SUFFIX,gongm.in,Proxy
DOMAIN-SUFFIX,goo.gl,Proxy
DOMAIN-SUFFIX,goodreaders.com,Proxy
DOMAIN-SUFFIX,goodreads.com,Proxy
DOMAIN-SUFFIX,gravatar.com,Proxy
DOMAIN-SUFFIX,gstatic.com,Proxy
DOMAIN-SUFFIX,gvt0.com,Proxy
DOMAIN-SUFFIX,hockeyapp.net,Proxy
DOMAIN-SUFFIX,hotmail.com,Proxy
DOMAIN-SUFFIX,icons8.com,Proxy
DOMAIN-SUFFIX,ifixit.com,Proxy
DOMAIN-SUFFIX,ift.tt,Proxy
DOMAIN-SUFFIX,ifttt.com,Proxy
DOMAIN-SUFFIX,iherb.com,Proxy
DOMAIN-SUFFIX,imageshack.us,Proxy
DOMAIN-SUFFIX,img.ly,Proxy
DOMAIN-SUFFIX,imgur.com,Proxy
DOMAIN-SUFFIX,imore.com,Proxy
DOMAIN-SUFFIX,instapaper.com,Proxy
DOMAIN-SUFFIX,ipn.li,Proxy
DOMAIN-SUFFIX,is.gd,Proxy
DOMAIN-SUFFIX,issuu.com,Proxy
DOMAIN-SUFFIX,itgonglun.com,Proxy
DOMAIN-SUFFIX,itun.es,Proxy
DOMAIN-SUFFIX,ixquick.com,Proxy
DOMAIN-SUFFIX,j.mp,Proxy
DOMAIN-SUFFIX,js.revsci.net,Proxy
DOMAIN-SUFFIX,jshint.com,Proxy
DOMAIN-SUFFIX,jtvnw.net,Proxy
DOMAIN-SUFFIX,justgetflux.com,Proxy
DOMAIN-SUFFIX,kat.cr,Proxy
DOMAIN-SUFFIX,klip.me,Proxy
DOMAIN-SUFFIX,libsyn.com,Proxy
DOMAIN-SUFFIX,linkedin.com,Proxy
DOMAIN-SUFFIX,line-apps.com,Proxy
DOMAIN-SUFFIX,linode.com,Proxy
DOMAIN-SUFFIX,lithium.com,Proxy
DOMAIN-SUFFIX,littlehj.com,Proxy
DOMAIN-SUFFIX,live.com,Proxy
DOMAIN-SUFFIX,live.net,Proxy
DOMAIN-SUFFIX,livefilestore.com,Proxy
DOMAIN-SUFFIX,llnwd.net,Proxy
DOMAIN-SUFFIX,macid.co,Proxy
DOMAIN-SUFFIX,macromedia.com,Proxy
DOMAIN-SUFFIX,macrumors.com,Proxy
DOMAIN-SUFFIX,mashable.com,Proxy
DOMAIN-SUFFIX,mathjax.org,Proxy
DOMAIN-SUFFIX,medium.com,Proxy
DOMAIN-SUFFIX,mega.co.nz,Proxy
DOMAIN-SUFFIX,mega.nz,Proxy
DOMAIN-SUFFIX,megaupload.com,Proxy
DOMAIN-SUFFIX,microsofttranslator.com,Proxy
DOMAIN-SUFFIX,mindnode.com,Proxy
DOMAIN-SUFFIX,mobile01.com,Proxy
DOMAIN-SUFFIX,modmyi.com,Proxy
DOMAIN-SUFFIX,msedge.net,Proxy
DOMAIN-SUFFIX,myfontastic.com,Proxy
DOMAIN-SUFFIX,name.com,Proxy
DOMAIN-SUFFIX,nextmedia.com,Proxy
DOMAIN-SUFFIX,nsstatic.net,Proxy
DOMAIN-SUFFIX,nssurge.com,Proxy
DOMAIN-SUFFIX,nyt.com,Proxy
DOMAIN-SUFFIX,nytimes.com,Proxy
DOMAIN-SUFFIX,omnigroup.com,Proxy
DOMAIN-SUFFIX,onedrive.com,Proxy
DOMAIN-SUFFIX,onenote.com,Proxy
DOMAIN-SUFFIX,ooyala.com,Proxy
DOMAIN-SUFFIX,openvpn.net,Proxy
DOMAIN-SUFFIX,openwrt.org,Proxy
DOMAIN-SUFFIX,orkut.com,Proxy
DOMAIN-SUFFIX,osxdaily.com,Proxy
DOMAIN-SUFFIX,outlook.com,Proxy
DOMAIN-SUFFIX,ow.ly,Proxy
DOMAIN-SUFFIX,paddleapi.com,Proxy
DOMAIN-SUFFIX,parallels.com,Proxy
DOMAIN-SUFFIX,parse.com,Proxy
DOMAIN-SUFFIX,pdfexpert.com,Proxy
DOMAIN-SUFFIX,periscope.tv,Proxy
DOMAIN-SUFFIX,pinboard.in,Proxy
DOMAIN-SUFFIX,pinterest.com,Proxy
DOMAIN-SUFFIX,pixelmator.com,Proxy
DOMAIN-SUFFIX,pixiv.net,Proxy
DOMAIN-SUFFIX,playpcesor.com,Proxy
DOMAIN-SUFFIX,playstation.com,Proxy
DOMAIN-SUFFIX,playstation.com.hk,Proxy
DOMAIN-SUFFIX,playstation.net,Proxy
DOMAIN-SUFFIX,playstationnetwork.com,Proxy
DOMAIN-SUFFIX,pushwoosh.com,Proxy
DOMAIN-SUFFIX,rime.im,Proxy
DOMAIN-SUFFIX,servebom.com,Proxy
DOMAIN-SUFFIX,sfx.ms,Proxy
DOMAIN-SUFFIX,shadowsocks.org,Proxy
DOMAIN-SUFFIX,sharethis.com,Proxy
DOMAIN-SUFFIX,shazam.com,Proxy
DOMAIN-SUFFIX,skype.com,Proxy
DOMAIN-SUFFIX,smartdnsProxy.com,Proxy
DOMAIN-SUFFIX,smartmailcloud.com,Proxy
DOMAIN-SUFFIX,sndcdn.com,Proxy
DOMAIN-SUFFIX,sony.com,Proxy
DOMAIN-SUFFIX,soundcloud.com,Proxy
DOMAIN-SUFFIX,sourceforge.net,Proxy
DOMAIN-SUFFIX,spotify.com,Proxy
DOMAIN-SUFFIX,squarespace.com,Proxy
DOMAIN-SUFFIX,sstatic.net,Proxy
DOMAIN-SUFFIX,st.luluku.pw,Proxy
DOMAIN-SUFFIX,stackoverflow.com,Proxy
DOMAIN-SUFFIX,startpage.com,Proxy
DOMAIN-SUFFIX,staticflickr.com,Proxy
DOMAIN-SUFFIX,steamcommunity.com,Proxy
DOMAIN-SUFFIX,symauth.com,Proxy
DOMAIN-SUFFIX,symcb.com,Proxy
DOMAIN-SUFFIX,symcd.com,Proxy
DOMAIN-SUFFIX,tapbots.com,Proxy
DOMAIN-SUFFIX,tapbots.net,Proxy
DOMAIN-SUFFIX,tdesktop.com,Proxy
DOMAIN-SUFFIX,techcrunch.com,Proxy
DOMAIN-SUFFIX,techsmith.com,Proxy
DOMAIN-SUFFIX,thepiratebay.org,Proxy
DOMAIN-SUFFIX,theverge.com,Proxy
DOMAIN-SUFFIX,time.com,Proxy
DOMAIN-SUFFIX,timeinc.net,Proxy
DOMAIN-SUFFIX,tiny.cc,Proxy
DOMAIN-SUFFIX,tinypic.com,Proxy
DOMAIN-SUFFIX,tmblr.co,Proxy
DOMAIN-SUFFIX,todoist.com,Proxy
DOMAIN-SUFFIX,trello.com,Proxy
DOMAIN-SUFFIX,trustasiassl.com,Proxy
DOMAIN-SUFFIX,tumblr.co,Proxy
DOMAIN-SUFFIX,tumblr.com,Proxy
DOMAIN-SUFFIX,tweetdeck.com,Proxy
DOMAIN-SUFFIX,tweetmarker.net,Proxy
DOMAIN-SUFFIX,twitch.tv,Proxy
DOMAIN-SUFFIX,txmblr.com,Proxy
DOMAIN-SUFFIX,typekit.net,Proxy
DOMAIN-SUFFIX,ubertags.com,Proxy
DOMAIN-SUFFIX,ublock.org,Proxy
DOMAIN-SUFFIX,ubnt.com,Proxy
DOMAIN-SUFFIX,ulyssesapp.com,Proxy
DOMAIN-SUFFIX,urchin.com,Proxy
DOMAIN-SUFFIX,usertrust.com,Proxy
DOMAIN-SUFFIX,v.gd,Proxy
DOMAIN-SUFFIX,v2ex.com,Proxy
DOMAIN-SUFFIX,vimeo.com,Proxy
DOMAIN-SUFFIX,vimeocdn.com,Proxy
DOMAIN-SUFFIX,vine.co,Proxy
DOMAIN-SUFFIX,vivaldi.com,Proxy
DOMAIN-SUFFIX,vox-cdn.com,Proxy
DOMAIN-SUFFIX,vsco.co,Proxy
DOMAIN-SUFFIX,vultr.com,Proxy
DOMAIN-SUFFIX,w.org,Proxy
DOMAIN-SUFFIX,w3schools.com,Proxy
DOMAIN-SUFFIX,webtype.com,Proxy
DOMAIN-SUFFIX,wikiwand.com,Proxy
DOMAIN-SUFFIX,wikileaks.org,Proxy
DOMAIN-SUFFIX,wikimedia.org,Proxy
DOMAIN-SUFFIX,wikipedia.com,Proxy
DOMAIN-SUFFIX,wikipedia.org,Proxy
DOMAIN-SUFFIX,windows.com,Proxy
DOMAIN-SUFFIX,windows.net,Proxy
DOMAIN-SUFFIX,wire.com,Proxy
DOMAIN-SUFFIX,wordpress.com,Proxy
DOMAIN-SUFFIX,workflowy.com,Proxy
DOMAIN-SUFFIX,wp.com,Proxy
DOMAIN-SUFFIX,wsj.com,Proxy
DOMAIN-SUFFIX,wsj.net,Proxy
DOMAIN-SUFFIX,xda-developers.com,Proxy
DOMAIN-SUFFIX,xeeno.com,Proxy
DOMAIN-SUFFIX,xiti.com,Proxy
DOMAIN-SUFFIX,yahoo.com,Proxy
DOMAIN-SUFFIX,yimg.com,Proxy
DOMAIN-SUFFIX,ying.com,Proxy
DOMAIN-SUFFIX,yoyo.org,Proxy
DOMAIN-SUFFIX,ytimg.com,Proxy

# Telegram
DOMAIN-SUFFIX,telegra.ph,Proxy
DOMAIN-SUFFIX,telegram.org,Proxy

IP-CIDR,91.108.4.0/22,Proxy,no-resolve
IP-CIDR,91.108.8.0/21,Proxy,no-resolve
IP-CIDR,91.108.16.0/22,Proxy,no-resolve
IP-CIDR,91.108.56.0/22,Proxy,no-resolve
IP-CIDR,149.154.160.0/20,Proxy,no-resolve
IP-CIDR6,2001:67c:4e8::/48,Proxy,no-resolve
IP-CIDR6,2001:b28:f23d::/48,Proxy,no-resolve
IP-CIDR6,2001:b28:f23f::/48,Proxy,no-resolve

# Google 中国服务 services.googleapis.cn
IP-CIDR,120.232.181.162/32,Proxy,no-resolve
IP-CIDR,120.241.147.226/32,Proxy,no-resolve
IP-CIDR,120.253.253.226/32,Proxy,no-resolve
IP-CIDR,120.253.255.162/32,Proxy,no-resolve
IP-CIDR,120.253.255.34/32,Proxy,no-resolve
IP-CIDR,120.253.255.98/32,Proxy,no-resolve
IP-CIDR,180.163.150.162/32,Proxy,no-resolve
IP-CIDR,180.163.150.34/32,Proxy,no-resolve
IP-CIDR,180.163.151.162/32,Proxy,no-resolve
IP-CIDR,180.163.151.34/32,Proxy,no-resolve
IP-CIDR,203.208.39.0/24,Proxy,no-resolve
IP-CIDR,203.208.40.0/24,Proxy,no-resolve
IP-CIDR,203.208.41.0/24,Proxy,no-resolve
IP-CIDR,203.208.43.0/24,Proxy,no-resolve
IP-CIDR,203.208.50.0/24,Proxy,no-resolve
IP-CIDR,220.181.174.162/32,Proxy,no-resolve
IP-CIDR,220.181.174.226/32,Proxy,no-resolve
IP-CIDR,220.181.174.34/32,Proxy,no-resolve

RULE-SET,LAN,DIRECT

# 剩余未匹配的国内网站
DOMAIN-SUFFIX,cn,DIRECT
DOMAIN-KEYWORD,-cn,DIRECT

# 最终规则
GEOIP,CN,DIRECT
FINAL,Proxy,dns-failed
`
