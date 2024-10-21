package sflogger

import (
	"fmt"
	"github.com/omygoden/gotools/sfconst"
	"github.com/sirupsen/logrus"
	"log"
	"os"
	"time"
	"zgin/global"
	"zgin/pkg/env"
)

func InitLogger() {
	logs := logrus.New()

	logs.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:   true,
		DisableColors:   true,
		TimestampFormat: "2006-01-02 15:04:05",
	})
	//logs.SetReportCaller(true)
	if _, err := os.Stat(fmt.Sprintf("%s/%s", env.PROJECT_PATH, global.Config.App.LogSavePath)); err != nil {
		err = os.Mkdir(fmt.Sprintf("%s/%s", env.PROJECT_PATH, global.Config.App.LogSavePath), os.ModePerm)
		if err != nil {
			log.Println("日志目录创建失败：", err.Error())
			os.Exit(1)
		}
	}

	global.Logger = logs
	log.Println("日志初始化成功")
}

func Info(filename, content string, extra map[string]interface{}) {
	global.LoggerLock.Lock()
	defer global.LoggerLock.Unlock()

	fd, _ := os.OpenFile(getFilePath(filename), os.O_CREATE|os.O_APPEND|os.O_WRONLY, os.ModePerm)
	defer fd.Close()

	global.Logger.SetOutput(fd)
	global.Logger.SetFormatter(&logrus.JSONFormatter{})
	global.Logger.WithFields(extra).Info(content)
}

func Error(filename, content string, extra map[string]interface{}) {
	global.LoggerLock.Lock()
	defer global.LoggerLock.Unlock()

	fd, _ := os.OpenFile(getFilePath(filename), os.O_CREATE|os.O_APPEND|os.O_WRONLY, os.ModePerm)
	defer fd.Close()

	global.Logger.SetOutput(fd)
	global.Logger.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: sfconst.GO_TIME_FULL,
	})

	global.Logger.WithFields(extra).Error(content, content, content)
}

var fileMap = make(map[string]string)

func getFilePath(fileName string) string {
	filePath := fmt.Sprintf("%s/%s/%s", env.PROJECT_PATH, global.Config.App.LogSavePath, time.Now().Format(sfconst.GO_TIME_YMD))
	fullPath := filePath + "/" + fileName

	if v, ok := fileMap[fullPath]; ok {
		return v
	}

	if _, err := os.Stat(filePath); err != nil {
		os.Mkdir(filePath, os.ModePerm)
	}

	fileMap[fullPath] = fullPath

	return fullPath
}
