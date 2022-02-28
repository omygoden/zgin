package goroutinepool

import (
	"fmt"
	"runtime"
	"time"
	"zgin/global"
	"zgin/pkg/sflogger"
)

func GoroutineListen() {
	go func() {
		for {
			sflogger.GoroutineLog(fmt.Sprintf("【当前总协程数】:%d，【当前协程池】:%d", runtime.NumGoroutine(), len(global.GoroutinePool)))

			time.Sleep(time.Second * 3)
		}
	}()
}

func InitGoroutinePool() {
	global.GoroutinePool = make(chan int, global.Config.Goroutine.MaxPullGoruntine)
	for i := 0; i < global.Config.Goroutine.MaxPullGoruntine; i++ {
		global.GoroutinePool <- 1
	}
}
