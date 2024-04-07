package model

// config配置
type Config struct {
	SystemParams SystemParams `mapstructure:"system" json:"system" yaml:"system"`
	Mysql        Mysql        `mapstructure:"mysql" json:"mysql" yaml:"mysql"` // gorm
	Sqlite       Sqlite       `mapstructure:"sqlite" json:"sqlite" yaml:"sqlite"`
}

// mysql配置
type Mysql struct {
	Address      string `mapstructure:"address" json:"address" yaml:"address"`                      // 服务器地址:端口
	Port         string `mapstructure:"port" json:"port" yaml:"port"`                               //:端口
	Config       string `mapstructure:"config" json:"config" yaml:"config"`                         // 高级配置
	Dbname       string `mapstructure:"db-name" json:"db-name" yaml:"db-name"`                      // 数据库名
	Username     string `mapstructure:"username" json:"username" yaml:"username"`                   // 数据库用户名
	Password     string `mapstructure:"password" json:"password" yaml:"password"`                   // 数据库密码
	Engine       string `mapstructure:"engine" json:"engine" yaml:"engine" default:"InnoDB"`        //数据库引擎，默认InnoDB
	MaxIdleConns int64  `mapstructure:"max-idle-conns" json:"max-idle-conns" yaml:"max-idle-conns"` // 空闲中的最大连接数
	MaxOpenConns int64  `mapstructure:"max-open-conns" json:"max-open-conns" yaml:"max-open-conns"` // 打开到数据库的最大连接数

}

type Sqlite struct {
	Path string `mapstructure:"path" yaml:"path"`
}
type SystemParams struct {
	Mode          string `mapstructure:"mode"    yaml:"mode"` //release=正常 dev=开发
	AdminEmail    string `mapstructure:"admin-email"    yaml:"admin-email"`
	AdminPassword string `mapstructure:"admin-password" yaml:"admin-password"`
	HTTPPort      int    `mapstructure:"http-port"      yaml:"http-port"`
	HTTPSPort     int    `mapstructure:"https-port"     yaml:"https-port"`
	GRPCPort      int    `mapstructure:"grpc-port"      yaml:"grpc-port"`
	DbType        string `mapstructure:"db-type"        yaml:"db-type"`
}
