package file_plugin

import (
	"archive/tar"
	"compress/gzip"
	"io"
	"log"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"time"
)

// BinaryPath 二进制文件的绝对路径
func BinaryPath() (string, error) {
	exePath, err := os.Executable()
	if err != nil {
		return "", err
	}
	return filepath.EvalSymlinks(filepath.Dir(exePath))
}

// DoBinaryPath 执行二进制文件的绝对路径
func DoBinaryPath() (string, error) {
	return os.Getwd()
}

// DoCodePath 当前执行代码的绝对路径
func DoCodePath() (string, bool) {
	_, filename, _, ok := runtime.Caller(0)
	return filename, ok
}

// TarGzUnzip 解压tar.gz
func TarGzUnzip(zipFile, dest string) error {
	fr, err := os.Open(zipFile)
	if err != nil {
		return err
	}
	defer fr.Close()
	gr, err := gzip.NewReader(fr)
	if err != nil {
		return err
	}
	defer gr.Close()
	tr := tar.NewReader(gr)
	// 读取文件
	for {
		h, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		fw, err := os.OpenFile(path.Join(dest, h.Name), os.O_CREATE|os.O_WRONLY, 0666 /*os.FileMode(h.Mode)*/)
		if err != nil {
			return err
		}
		defer fw.Close()
		_, err = io.Copy(fw, tr)
		if err != nil {
			return err
		}
	}
	return nil
}
func SetOutputFile() (*os.File, error) {
	now := time.Now()
	logFileName := now.Format("2006-01-02") + ".log" //日志文件名
	logFilePath := ""                                //路径
	if dir, err := os.Getwd(); err == nil {          //当前工作目录
		logFilePath = dir + "/logs/"
	}
	_, err := os.Stat(logFilePath)
	if os.IsNotExist(err) {
		if err := os.MkdirAll(logFilePath, 0777); err != nil {
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
