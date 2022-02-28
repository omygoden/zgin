package main

import (
	"github.com/omygoden/gotools/sfconst"
	"github.com/robfig/cron"
	"log"
	"os"
	"time"
	"zgin/console/cmd"
	"zgin/pkg/setting"
)

func init() {
	setting.CronInitSetting(200, 100)
}

/**
*从左到右含义
Seconds      | Yes        | 0-59            | * / , -
Minutes      | Yes        | 0-59            | * / , -
Hours        | Yes        | 0-23            | * / , -
Day of month | Yes        | 1-31            | * / , -
Month        | Yes        | 1-12 or JAN-DEC | * / , -
Day of week  | Yes        | 0-6 or SUN-SAT  | * / , -
*/
func main() {
	var forever = make(chan int)
	log.Println("cron ready ...")

	crons := cron.New()
	_ = crons.AddFunc("* * * * * *", func() {
		file, _ := os.OpenFile("cron.log", os.O_APPEND|os.O_CREATE|os.O_RDWR, os.ModeAppend|os.ModePerm)
		log.Println("exec ...")
		_, _ = file.WriteString(time.Now().Format(sfconst.GO_TIME_FULL) + "\n")
	})
	_ = crons.AddFunc("0 * * * * *", cmd.Test)

	crons.Start()
	log.Println("cron start all ...")

	<-forever
}
