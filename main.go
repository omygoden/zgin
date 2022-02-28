package main

import (
	"fmt"
	"log"
	"net/http"
	_ "net/http/pprof"
	"zgin/global"
	"zgin/pkg/setting"
	"zgin/routers"
)

func init() {
	setting.InitSetting()
}

func main() {
	go func() {
		log.Println(http.ListenAndServe(fmt.Sprintf(":%s", global.Config.App.PprofPort), nil))
	}()

	router := routers.InitRouter()
	router.Run(fmt.Sprintf(":%s", global.Config.App.HttpPort))
}
