package sflogger

import (
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
	"zgin/global"
	"zgin/pkg/env"
)

type ClearLogs struct {
}

func NewClearLogs() *ClearLogs {
	return &ClearLogs{}
}

func (this *ClearLogs) TimerClearLogs() {
	go func() {
		t := time.NewTicker(time.Second)
		for {
			select {
			case <-t.C:
				tt := time.Now()
				if tt.Hour() == 0 && tt.Minute() == 0 && tt.Second() == 0 {
					this.ClearHistoryLog()
				}
			}
		}
	}()
}

func (this *ClearLogs) ClearHistoryLog() {
	var logs []string
	logPath := fmt.Sprintf("%s/%s", env.PROJECT_PATH, global.Config.App.LogSavePath)
	filepath.Walk(logPath, func(path string, info fs.FileInfo, err error) error {
		if info.Name() == global.Config.App.LogSavePath {
			return nil
		}
		if info.IsDir() {
			logs = append(logs, info.Name())
		}
		return nil
	})
	log.Println(fmt.Sprintf("当前日志文件数:%d,最大支持文件数:%d", len(logs), global.Config.App.LogMaxNum))
	if len(logs) > global.Config.App.LogMaxNum {
		var n = len(logs) - global.Config.App.LogMaxNum
		filepath.Walk(logPath, func(path string, info fs.FileInfo, err error) error {
			if n <= 0 {
				return nil
			}
			if strings.Contains(global.Config.App.LogSavePath, info.Name()) {
				return nil
			}
			if info.IsDir() {
				p := fmt.Sprintf("%s/%s", logPath, info.Name())
				e := os.RemoveAll(p)
				log.Println(fmt.Sprintf("删除文件：%s，结果:%v", p, e))
				n--
			}
			return nil
		})
	}

}
