package initialize

import (
	"AirGo/global"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

// 读取配置文件，并转换成 struct结构
func InitViper() *viper.Viper {
	v := viper.New()
	v.SetConfigFile("config.yaml") //config路径
	v.SetConfigType("yaml")        //设置文件的类型
	err := v.ReadInConfig()
	if err != nil {
		global.Logrus.Panic("Fatal error config file:", err)
		//panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}
	v.WatchConfig()

	v.OnConfigChange(func(e fsnotify.Event) {
		global.Logrus.Info("config file changed:", e.Name)
		if err = v.Unmarshal(&global.Config); err != nil {
			global.Logrus.Error(err)
			//fmt.Println(err)
		}
	})
	if err = v.Unmarshal(&global.Config); err != nil { //解析到全局配置
		global.Logrus.Error(err)
		//fmt.Println(err)
	}
	return v
}
