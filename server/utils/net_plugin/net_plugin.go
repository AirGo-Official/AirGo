package net_plugin

import (
	"compress/gzip"
	"crypto/tls"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"net"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/AirGo-Official/AirGo/utils/response"
	"github.com/dustin/go-humanize"
	"github.com/gin-gonic/gin"
	"github.com/go-resty/resty/v2"
	"golang.org/x/net/context"
)

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
		Timeout: timeOut, //超时时间
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

// DownloadFile 下载文件
func DownloadFile(url, fileName string, mode fs.FileMode) error {
	client := resty.New()
	resp, err := client.R().SetOutput(fileName).Get(url)
	//resp.RawResponse.ContentLength
	//fmt.Println("resp.RawResponse.ContentLength:", resp.RawResponse.ContentLength)
	if err != nil {
		return err
	}
	if resp.StatusCode() != http.StatusOK {
		return errors.New(fmt.Sprintf("%s%d", "Download failed: code = ", resp.StatusCode()))
	}
	err = os.Chmod(fileName, mode)
	if err != nil {
		return err
	}
	return nil
}

type WriteCounter struct {
	Total uint64
	Ctx   *gin.Context
	num   uint
}

func (wc *WriteCounter) Write(p []byte) (int, error) {
	n := len(p)
	wc.Total += uint64(n)
	if wc.num == 80 { //间隔发送
		wc.PrintProgress()
	}
	if wc.num--; wc.num == 0 {
		wc.num = 80
	}
	return n, nil
}
func (wc WriteCounter) PrintProgress() {
	response.ResponseSSE("message", fmt.Sprintf("Downloading........................ %s completed", humanize.Bytes(wc.Total)), wc.Ctx)
}

func DownloadFileWithProgress(url, fileName string, mode fs.FileMode, ctx *gin.Context) error {
	out, err := os.Create(fileName + ".tmp")
	if err != nil {
		return err
	}
	resp, err := http.Get(url)
	if err != nil {
		out.Close()
		return err
	}
	defer resp.Body.Close()
	counter := &WriteCounter{Ctx: ctx, num: 80}
	if _, err = io.Copy(out, io.TeeReader(resp.Body, counter)); err != nil {
		out.Close()
		return err
	}
	out.Close()
	if err = os.Rename(fileName+".tmp", fileName); err != nil {
		return err
	}
	_ = os.Chmod(fileName, mode)
	return nil
}
