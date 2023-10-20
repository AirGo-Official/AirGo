package logrus_plugin

import (
	"log"
	"os"
	"path"
	"time"

	"github.com/sirupsen/logrus"
)

// Panic：记录日志，然后panic。
// Fatal：致命错误，出现错误时程序无法正常运转。输出日志后，程序退出；
// Error：错误日志，需要查看原因；
// Warn：警告信息，提醒程序员注意；
// Info：关键操作，核心流程的日志；
// Debug：一般程序中输出的调试信息；
// Trace：很细粒度的信息，一般用不到；
func InitLogrus() *logrus.Logger {
	//实例化
	logger := logrus.New()
	logger.SetReportCaller(true) //在输出日志中添加文件名和方法信息
	src, _ := SetOutputFile()
	//设置输出
	logger.Out = src
	//logger.Out = io.MultiWriter(src, os.Stdout) //同时打印到控制台及日志里
	//设置最低日志级别
	//logger.SetLevel(logrus.DebugLevel)
	logger.SetLevel(logrus.InfoLevel)
	//设置日志格式
	//logger.SetFormatter(&logrus.JSONFormatter{})
	logger.SetFormatter(&logrus.TextFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
	})
	return logger
}

func SetOutputFile() (*os.File, error) {
	now := time.Now()
	logFileName := now.Format("2006-01-02") + ".log" //日志文件名
	logFilePath := ""                                //路径
	if dir, err := os.Getwd(); err == nil {          //当前工作目录
		logFilePath = dir + "/logs/"
	}
	_, err := os.Stat(logFilePath)
	if os.IsNotExist(err) { //isNotExist()判断为true，说明文件或者文件夹不存在
		if err := os.MkdirAll(logFilePath, 0777); err != nil { //创建
			log.Println(err.Error())
			return nil, err
		}
	}

	//日志文件
	fileName := path.Join(logFilePath, logFileName)
	if _, err := os.Stat(fileName); err != nil {
		if _, err := os.Create(fileName); err != nil {
			log.Println(err.Error())
			return nil, err
		}
	}
	//写入文件
	src, err := os.OpenFile(fileName, os.O_APPEND|os.O_WRONLY, os.ModeAppend) //如果已经存在，则在尾部添加写
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return src, nil
}
