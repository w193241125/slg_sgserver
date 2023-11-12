package config

import (
	"errors"
	"github.com/Unknwon/goconfig"
	"log"
	"os"
)

const configFile = "/conf/conf.ini"

var File *goconfig.ConfigFile

// 加载这个文件的时候,先init初始化...
func init() {
	//拿到当前文件所在目录
	currentDir, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	configPath := currentDir + configFile

	if !fileExist(configPath) {
		panic(errors.New("配置文件不存在")) // 无配置, 应用终止
	}

	//构建后启动时添加配置目录 sgserver.exe  M:/sgserver/conf/conf.ini
	len := len(os.Args)
	if len > 1 {
		dir := os.Args[1]
		if dir != "" {
			configPath = dir + configFile
		}
	}

	//文件系统的读取
	File, err = goconfig.LoadConfigFile(configPath)
	if err != nil {
		log.Fatal("读取配置文件出错:", err)
	}
}

func fileExist(fileName string) bool {
	_, err := os.Stat(fileName)

	return err == nil || os.IsExist(err)
}
