package sflogger

import (
	"encoding/json"
	"fmt"
	"github.com/omygoden/gotools/sfconst"
	"os"
	"time"
	"zgin/global"
)

type CustomLogger struct {
	Prefix   string
	FileName string
}

func NewCustomLogger(fileName string) *CustomLogger {
	return &CustomLogger{
		FileName: fileName,
	}
}

func (this *CustomLogger) SetPrefix(prefix string) {
	this.Prefix = prefix
}

func (this *CustomLogger) Start() {
	this.write(fmt.Sprintf("---------- [任务开始] ----------\n"))
}

func (this *CustomLogger) Info(msg interface{}) {
	this.write(fmt.Sprintf("[%s] %s INFO:%v\n", time.Now().Format(sfconst.GO_TIME_FULL), this.Prefix, msg))
}

func (this *CustomLogger) InfoMap(m map[string]interface{}) {
	b, _ := json.Marshal(m)
	this.write(fmt.Sprintf("[%s] %s INFO:%v\n", time.Now().Format(sfconst.GO_TIME_FULL), this.Prefix, string(b)))
}

func (this *CustomLogger) Error(msg interface{}) {
	this.write(fmt.Sprintf("[%s] %s ERROR:%v\n", time.Now().Format(sfconst.GO_TIME_FULL), this.Prefix, msg))
}

func (this *CustomLogger) ErrorMap(m map[string]interface{}) {
	b, _ := json.Marshal(m)
	this.write(fmt.Sprintf("[%s] %s ERROR:%v\n", time.Now().Format(sfconst.GO_TIME_FULL), this.Prefix, string(b)))
}

func (this *CustomLogger) Finish() {
	this.write(fmt.Sprintf("---------- [任务结束] ----------\n"))
}

func (this *CustomLogger) write(msg string) {
	global.LoggerLock.Lock()
	defer global.LoggerLock.Unlock()

	fd, _ := os.OpenFile(getFilePath(this.FileName), os.O_CREATE|os.O_APPEND|os.O_WRONLY, os.ModePerm)
	defer fd.Close()
	fd.WriteString(msg)
}
