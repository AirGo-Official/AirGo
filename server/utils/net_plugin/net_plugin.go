package net_plugin

import (
	"compress/gzip"
	"crypto/tls"
	"fmt"
	"golang.org/x/net/context"
	"io"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type Ip3030 struct {
	IP       string `json:"ip"`
	Location string `json:"location"`
}

// 通过http代理访问网站
func GetByHTTPProxy(objUrl, proxyAddress string, proxyPort int, timeOut time.Duration) (*http.Response, error) {
	proxy := func(_ *http.Request) (*url.URL, error) {
		return url.Parse(fmt.Sprintf("http://%s:%d", proxyAddress, proxyPort))
	}
	transport := &http.Transport{Proxy: proxy}
	client := &http.Client{
		Transport: transport,
		Timeout:   timeOut,
	}
	return client.Get(objUrl)
}

// 通过Socks5代理访问网站
func GetBySocks5Proxy(objUrl, proxyAddress string, proxyPort int, timeOut time.Duration) (*http.Response, error) {

	proxy := func(_ *http.Request) (*url.URL, error) {
		return url.Parse(fmt.Sprintf("socks5://%s:%d", proxyAddress, proxyPort))
	}
	transport := &http.Transport{Proxy: proxy}
	client := &http.Client{
		Transport: transport,
		Timeout:   timeOut,
	}
	return client.Get(objUrl)
}

// 不通过代理访问网站
func GetNoProxy(url string, timeOut time.Duration) (*http.Response, error) {
	client := &http.Client{
		Timeout: timeOut,
	}
	return client.Get(url)
}

// 发送tcp
func Tcp(address string, port int) {
	dialer := net.Dialer{Timeout: 3 * time.Second}
	conn, err := dialer.Dial("tcp", fmt.Sprintf("%s:%d", address, port))
	if err != nil {
		return
	}
	conn.Close()
}

// socks5 http.Client
func ClientWithSocks5(proxyAddress string, proxyPort int, timeOut time.Duration) *http.Client {
	proxy := func(_ *http.Request) (*url.URL, error) {
		return url.Parse(fmt.Sprintf("socks5://%s:%d", proxyAddress, proxyPort))
	}
	transport := &http.Transport{
		Proxy: proxy,
	}
	return &http.Client{
		Transport: transport,
		Timeout:   timeOut,
	}
}

// 自定义dns http.Client
func ClientWithDNS(dns string, timeOut time.Duration) *http.Client {
	dialer := &net.Dialer{
		Timeout: timeOut,
	}
	dialer.Resolver = &net.Resolver{
		PreferGo: true,
		Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
			return dialer.DialContext(ctx, "udp", fmt.Sprintf("%s:%d", dns, 53)) // 请求nameserver解析域名
		},
	}
	return &http.Client{
		Timeout: time.Duration(5) * time.Second, //超时时间
		Transport: &http.Transport{
			MaxIdleConnsPerHost:   200,   //单个路由最大空闲连接数
			MaxConnsPerHost:       10000, //单个路由最大连接数
			IdleConnTimeout:       90 * time.Second,
			TLSHandshakeTimeout:   10 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
			DialContext:           dialer.DialContext,
			TLSClientConfig:       &tls.Config{InsecureSkipVerify: true}, //设置client信任所有证书
		},
	}
}

// 自定义net.Resolver，用于dns 查询
func Resolver(dns string, timeOut time.Duration) *net.Resolver {
	dialer := &net.Dialer{
		Timeout: timeOut,
	}

	return &net.Resolver{
		PreferGo: true,
		Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
			return dialer.DialContext(ctx, "udp", fmt.Sprintf("%s:%d", dns, 53)) // 请求nameserver解析域名
		},
	}
}

// 自定义net.Dialer，用于dns 查询,节点tcping
func Dialer(dns string, timeOut time.Duration) *net.Dialer {
	return &net.Dialer{
		Timeout:  timeOut,
		Resolver: Resolver(dns, timeOut),
	}
}

// 读取http响应的内容
func ReadDate(resp *http.Response) string {
	// 是否有 gzip
	gzipFlag := false
	for k, v := range resp.Header {
		if strings.ToLower(k) == "content-encoding" && strings.ToLower(v[0]) == "gzip" {
			gzipFlag = true
		}
	}

	var content []byte
	if gzipFlag {
		// 创建 gzip.Reader
		gr, err := gzip.NewReader(resp.Body)
		defer gr.Close()
		if err != nil {
			fmt.Println(err.Error())
		}
		content, _ = io.ReadAll(gr)
	} else {
		content, _ = io.ReadAll(resp.Body)
	}

	return string(content)
}
