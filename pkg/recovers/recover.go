package recovers

import (
	"fmt"
	"log"
	"runtime/debug"
	"zgin/pkg/constant"
	"zgin/pkg/sflogger"
)

//队列/协程等异步异常处理，记录日志
func GoruntineRecoverHandle(title string) {
	if err := recover(); err != nil {
		log.Println(title, ":", err)
		m := debug.Stack()
		sflogger.Error(constant.LOG_PANIC, fmt.Sprintf("【%s】异常信息", title), map[string]interface{}{
			"错误信息": err,
			"异常信息": string(m),
		})
	}
}
