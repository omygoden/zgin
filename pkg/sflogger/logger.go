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
	if _, err := os.Stat(fmt.Sprintf("%s/%s/%s", os.Getenv("GOPATH"), env.PROJECT_NAME, global.Config.App.LogSavePath)); err != nil {
		err = os.Mkdir(fmt.Sprintf("%s/%s/%s", os.Getenv("GOPATH"), env.PROJECT_NAME, global.Config.App.LogSavePath), os.ModePerm)
		if err != nil {
			log.Println("日志目录创建失败：", err.Error())
			os.Exit(1)
		}
	}

	global.Logger = logs
}

func Info(filename, content string, extra map[string]interface{}) {
	global.LoggerLock.Lock()
	defer global.LoggerLock.Unlock()

	fd, _ := os.OpenFile(getFilePath(filename), os.O_CREATE|os.O_APPEND|os.O_WRONLY, os.ModePerm)
	defer fd.Close()

	global.Logger.SetOutput(fd)
	global.Logger.WithFields(extra).Info(content)
}

func Error(filename, content string, extra map[string]interface{}) {
	global.LoggerLock.Lock()
	defer global.LoggerLock.Unlock()

	fd, _ := os.OpenFile(getFilePath(filename), os.O_CREATE|os.O_APPEND|os.O_WRONLY, os.ModePerm)
	defer fd.Close()

	global.Logger.SetOutput(fd)
	global.Logger.WithFields(extra).Error(content)
}

func getFilePath(filename string) string {
	filepath := fmt.Sprintf("%s/%s/%s/%s", os.Getenv("GOPATH"), env.PROJECT_NAME, global.Config.App.LogSavePath, time.Now().Format(sfconst.GO_TIME_YMD))

	if _, err := os.Stat(filepath); err != nil {
		os.Mkdir(filepath, os.ModePerm)
	}
	filepath = filepath + "/" + filename

	return filepath
}
