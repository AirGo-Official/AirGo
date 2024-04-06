package admin_logic

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-resty/resty/v2"
	"github.com/ppoonk/AirGo/constant"
	"github.com/ppoonk/AirGo/global"
	"github.com/ppoonk/AirGo/model"
	"github.com/ppoonk/AirGo/utils/file_plugin"
	"github.com/ppoonk/AirGo/utils/net_plugin"
	"github.com/ppoonk/AirGo/utils/response"
	"github.com/ppoonk/AirGo/utils/websocket_plugin"
	"github.com/tidwall/gjson"
	"gorm.io/gorm"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"runtime"
	"syscall"
	"time"
)

type System struct{}

// 修改系统配置
func (s *System) UpdateSetting(setting *model.Server) error {
	err := global.DB.Transaction(func(tx *gorm.DB) error {
		return tx.Save(&setting).Error
	})
	if err != nil {
		return err
	}
	global.Server = *setting
	//重新加载email
	global.GoroutinePool.Submit(func() {
		global.Logrus.Info("重新加载 email")
		EmailSvc.Reload()
	})
	//重新加载tg bot
	if global.Server.Notice.EnableTGBot {
		if global.Server.Notice.BotToken != "" {
			global.GoroutinePool.Submit(func() {
				global.Logrus.Info("重新加载tg bot")
				//关闭
				TgBotSvc.TGBotCloseListen()
				//重启
				TgBotSvc.TGBotStart()
			})
		}
	} else {
		global.GoroutinePool.Submit(func() {
			global.Logrus.Info("停止 tg bot")
			//关闭
			TgBotSvc.TGBotCloseListen()
		})
	}
	//重新加载通知消息时的管理员id
	global.GoroutinePool.Submit(func() {
		global.Logrus.Info("重新加载通知消息时的管理员id")
		PushMessageSvc.AdminAccountHandler()
	})
	return nil
}

// 查询最新版本AirGo
func (s *System) GetLatestVersion() (string, error) {
	client := resty.New()
	resp, err := client.R().Get(constant.AIRGO_GITHUB_API)
	if err != nil {
		return "", err
	}
	tag := resp.String()
	version := gjson.Get(tag, "tag_name")
	return version.String(), nil

}
func send(wsMessage *websocket_plugin.WsMessage, msgChannel chan<- *websocket_plugin.WsMessage, code int, msg string) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("panic:", err) //在websocket异常关闭后，向管道中写入会panic
		}
	}()
	wsMessage.Data = model.WebsocketResponse{
		Code: code,
		Msg:  msg,
	}
	msgChannel <- wsMessage
}

// 下载最新版本AirGo
func (s *System) DownloadLatestVersion(ctx *gin.Context) error {
	//防止重复操作
	_, ok := global.LocalCache.Get("DownloadLatestVersion")
	if ok {
		return errors.New("Frequent operations")
	} else {
		global.LocalCache.Set("DownloadLatestVersion", nil, 20*time.Second)
	}

	//1、获取版本
	response.ResponseSSE("message", "获取版本...", ctx)
	version, err := s.GetLatestVersion()
	if err != nil {
		return err
	}
	if version == "" {
		return errors.New("version is empty")
	}
	response.ResponseSSE("message", "版本："+version, ctx)
	// 2、下载
	response.ResponseSSE("message", "开始下载...", ctx)
	currentPath, err := file_plugin.DoBinaryPath()
	if err != nil {
		return err
	}
	filePath := path.Join(currentPath, "AirGo.tar.gz")
	// example:https://github.com/ppoonk/AirGo/releases/download/v0.2.1/AirGo-v0.2.1-darwin-arm64.tar.gz
	downloadFilePath := fmt.Sprintf("%s/%s/%s-%s-%s-%s%s", constant.AIRGO_GITHUB_DOWNLOAD_PRE, version, "AirGo", version, runtime.GOOS, runtime.GOARCH, ".tar.gz")
	//err = net_plugin.DownloadFile(downloadFilePath, filePath, 0666)
	err = net_plugin.DownloadFileWithProgress(downloadFilePath, filePath, 0666, ctx)
	if err != nil {
		return err
	}
	// 3、解压
	response.ResponseSSE("message", "开始解压文件...", ctx)
	tempPath := path.Join(currentPath, "temp")
	if err = os.MkdirAll(tempPath, 0777); err != nil { //创建
		return err
	}
	err = file_plugin.TarGzUnzip(filePath, tempPath)
	if err != nil {
		return err
	}
	// 判断是否文件是否存在、大小是否合法
	tempAirGoPath := filepath.Join(tempPath, "AirGo")
	tempFileInfo, err := os.Stat(tempAirGoPath)
	if err != nil {
		return err
	}
	if tempFileInfo.Size() < 40000000 {
		//v0.2.1版本是49697040，太小说明有问题
		return err
	}
	//权限
	err = os.Chmod(tempAirGoPath, 0777)
	if err != nil {
		return err
	}
	// 替换核心
	currentAirGoPath := filepath.Join(currentPath, "AirGo")
	currentOldAirGoPath := filepath.Join(currentPath, "AirGo-old")
	err = os.Rename(currentAirGoPath, currentOldAirGoPath)
	if err != nil {
		return err
	}
	err = os.Rename(tempAirGoPath, currentAirGoPath)
	if err != nil {
		return err
	}
	// 4、删除旧的二进制，删除临时目录
	response.ResponseSSE("message", "清理临时文件...", ctx)
	_ = os.RemoveAll(tempPath)
	_ = os.RemoveAll(currentOldAirGoPath)

	//5、重新初始化数据库 role_and_menu 、 menu 以及 casbin_rule
	response.ResponseSSE("message", "初始化数据库关键数据...", ctx)
	//使用新下载的二进制核心，通过命令行执行
	cmd := exec.Command(currentAirGoPath, "update")
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout // 标准输出
	cmd.Stderr = &stderr // 标准错误
	err = cmd.Run()
	if err != nil {
		return err
	}
	outStr, errStr := string(stdout.Bytes()), string(stderr.Bytes())
	if outStr != "" {
		response.ResponseSSE("message", outStr, ctx)
	}
	if errStr != "" {
		response.ResponseSSE("message", errStr, ctx)
	}
	if err != nil {
		return err
	}
	// 5、额外需要处理的数据
	s.ChangeDataForUpdate()
	// 6、重启（新进程）
	response.ResponseSSE("message", "3秒后自动重启...", ctx)
	err = s.Reload()
	if err != nil {
		return err
	}
	return nil
}

// 重启
func (s *System) Reload() error {
	p, err := os.FindProcess(syscall.Getpid())
	if err != nil {
		return err
	}
	go func() {
		ticker := time.NewTicker(3 * time.Second)
		defer ticker.Stop()
		<-ticker.C
		err = p.Signal(syscall.SIGHUP)
		if err != nil {
			global.Logrus.Error("Send signal error:", err)
		}
	}()
	return nil
}

// 版本升级时，额外需要处理的数据，例如数据库字段值批量修改
func (s *System) ChangeDataForUpdate() error {
	return global.DB.Transaction(func(tx *gorm.DB) error {
		err := tx.Exec("UPDATE node SET protocol = 'hysteria2' WHERE protocol = 'hysteria' ").Error
		if err != nil {
			return err
		}
		err = tx.Exec("UPDATE node SET vless_flow = '' WHERE vless_flow = 'none' ").Error
		if err != nil {
			return err
		}
		err = tx.Exec("UPDATE node SET security = '' WHERE security = 'none' ").Error
		if err != nil {
			return err
		}
		return nil
	})

}
