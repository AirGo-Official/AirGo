package main

import (
	"AirGo/global"
	"AirGo/initialize"
	"AirGo/service"
	"AirGo/utils/os_plugin"
	"flag"
	"runtime"
)

var start = flag.Bool("start", false, "启动")
var stop = flag.Bool("stop", false, "停止")
var resetAdmin = flag.Bool("resetAdmin", false, "重置管理员账户密码")

func main() {

	switch runtime.GOOS {
	case "darwin":
		initialize.InitializeAll() //初始化系统资源并启动路由

		//global.VP = initialize.InitViper() //初始化Viper
		//global.DB = initialize.Gorm()      //gorm连接数据库
		//initialize.InitServer()            //加载全局系统配置

	default:
		flag.Parse()
		if *start {
			initialize.InitializeAll() //初始化系统资源并启动路由
		} else if *stop {
			os_plugin.StopProcess("AirGo") //停止
		} else if *resetAdmin {
			global.VP = initialize.InitViper() //初始化Viper
			global.DB = initialize.Gorm()      //gorm连接数据库
			service.ResetAdminPassword()       // 重置管理员密码
		}
	}

}
