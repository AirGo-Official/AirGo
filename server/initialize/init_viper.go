package initialize

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/ppoonk/AirGo/global"
	"github.com/spf13/viper"
	"path"
)

// InitViper 读取配置文件
func InitViper(startConfigPath string) {
	v := viper.New()
	v.SetConfigFile(path.Join(startConfigPath)) //config路径
	v.SetConfigType("yaml")                     //设置文件的类型
	err := v.ReadInConfig()
	if err != nil {
		global.Logrus.Panic("Fatal error config file:", err)
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}
	v.WatchConfig()

	v.OnConfigChange(func(e fsnotify.Event) {
		global.Logrus.Info("config file changed:", e.Name)
		fmt.Println("config file changed:", e.Name)
		if err = v.Unmarshal(&global.Config); err != nil {
			global.Logrus.Error(err)
			fmt.Println(err)
		}
	})
	if err = v.Unmarshal(&global.Config); err != nil { //解析到全局配置
		global.Logrus.Error(err)
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
		//fmt.Println(err)
	}
	//global.Viper = v
	//return v
}
