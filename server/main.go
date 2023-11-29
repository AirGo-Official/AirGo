package main

import (
	"flag"
	"fmt"
	"github.com/ppoonk/AirGo/initialize"
	"github.com/ppoonk/AirGo/utils/os_plugin"
	"runtime"
)

const v = "old-version"

var start = flag.Bool("start", false, "启动")
var stop = flag.Bool("stop", false, "停止")
var resetAdmin = flag.Bool("resetAdmin", false, "重置管理员账户密码")
var version = flag.Bool("version", false, "版本")
var update = flag.Bool("update", false, "升级核心")

func main() {
	switch runtime.GOOS {
	case "darwin": //开发环境
		initialize.InitializeAll() //初始化系统资源并启动路由
		//global.VP = initialize.InitViper() //初始化Viper
		//global.DB = initialize.Gorm()      //gorm连接数据库
		//initialize.InitServer()            //加载全局系统配置
	default: //生产环境
		flag.Parse()
		if *start {
			initialize.InitializeAll() //初始化系统资源并启动路由
		} else if *stop {
			os_plugin.StopProcess("AirGo")
		} else if *resetAdmin {
			initialize.InitializeResetAdmin()
		} else if *version {
			fmt.Println(v)
		} else if *update {
			initialize.InitializeUpdate()
		}
	}

}
