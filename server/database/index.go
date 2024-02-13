package database

import (
	"fmt"
	"github.com/glebarez/sqlite"
	"github.com/ppoonk/AirGo/global"
	"github.com/ppoonk/AirGo/model"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

// StartGorm 初始化数据库并产生数据库全局变量
func StartGorm() {
	var err error
	switch global.Config.SystemParams.DbType {
	case "mysql":
		global.DB, err = GormMysql()
	case "sqlite":
		global.DB, err = GormSqlite()
	default:
		global.DB, err = GormMysql()
	}
	if err != nil {
		panic("Database connection failed:" + err.Error())
	}
	if global.DB != nil {
		if !global.DB.Migrator().HasTable(&model.User{}) {
			global.Logrus.Info("Start creating database and initializing data...")
			RegisterTables() //创建table
			InsertInto()     //导入数据
		} else {
			RegisterTables() //AutoMigrate 自动迁移 schema
		}
	} else {
		panic("Database connection failed")
	}
}

// GormSqlite 初始化sqlite数据库
func GormSqlite() (*gorm.DB, error) {
	//db, err := sql.Open("sqlite", ":memory:")

	if db, err := gorm.Open(sqlite.Open(global.Config.Sqlite.Path), &gorm.Config{
		//SkipDefaultTransaction: true, //关闭事务
		PrepareStmt: true,
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true, //单数表名
		},
	}); err != nil {
		global.Logrus.Error("gorm open error:", err)
		return nil, err
	} else {
		sqlDB, _ := db.DB()
		sqlDB.SetMaxIdleConns(int(global.Config.Mysql.MaxIdleConns))
		sqlDB.SetMaxOpenConns(int(global.Config.Mysql.MaxOpenConns))
		return db, nil
	}
}

// GormMysql 初始化Mysql数据库
func GormMysql() (*gorm.DB, error) {
	mysqlConfig := mysql.Config{
		DSN:                       fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?%s", global.Config.Mysql.Username, global.Config.Mysql.Password, global.Config.Mysql.Address, global.Config.Mysql.Port, global.Config.Mysql.Dbname, global.Config.Mysql.Config),
		DefaultStringSize:         191,
		SkipInitializeWithVersion: false,
	}
	if db, err := gorm.Open(mysql.New(mysqlConfig), &gorm.Config{
		//SkipDefaultTransaction: true, //关闭事务
		PrepareStmt: true,
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	}); err != nil {
		global.Logrus.Error("gorm open error:", err)
		return nil, err
	} else {
		db.InstanceSet("gorm:table_options", "ENGINE="+global.Config.Mysql.Engine)
		sqlDB, _ := db.DB()
		sqlDB.SetMaxIdleConns(int(global.Config.Mysql.MaxIdleConns))
		sqlDB.SetMaxOpenConns(int(global.Config.Mysql.MaxOpenConns))
		return db, nil
	}
}

// GormPgSql 初始化 Postgresql 数据库
func GormPgSql() *gorm.DB {
	pgsqlConfig := postgres.Config{
		DSN:                  "host=127.0.0.1 user=yourusername password=yourpassword dbname=userDB port=5432 sslmode=disable",
		PreferSimpleProtocol: false,
	}
	if db, err := gorm.Open(postgres.New(pgsqlConfig), &gorm.Config{
		//SkipDefaultTransaction: true, //关闭事务
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	}); err != nil {
		return nil
	} else {
		sqlDB, _ := db.DB()
		sqlDB.SetMaxIdleConns(int(global.Config.Mysql.MaxIdleConns))
		sqlDB.SetMaxOpenConns(int(global.Config.Mysql.MaxOpenConns))
		return db
	}
}
