package env

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"log"
	"os"
	"zgin/global"
	"zgin/pkg/util"
)

var PROJECT_PATH string

func InitEnv(configName string) {
	PROJECT_PATH = util.GetProjectPath()
	vpInit := viper.New()
	vpInit.SetConfigName(configName)
	vpInit.SetConfigType("yaml") // 高版本才支持ini后缀，早起版本不支持,(经测试：1.4不支持)
	vpInit.AddConfigPath(fmt.Sprintf("%s/%s", PROJECT_PATH, "conf"))

	err := vpInit.ReadInConfig()
	if err != nil {
		log.Println("环境配置读取失败:", err.Error())
		os.Exit(1)
	}

	switch configName {
	case "config":
		err = vpInit.Unmarshal(&global.Config)
	default:
		log.Println("环境配置文件名有误:", err.Error())
		os.Exit(1)
	}

	vpInit.WatchConfig()
	vpInit.OnConfigChange(func(in fsnotify.Event) {
		err = vpInit.Unmarshal(&global.Config)
		if err != nil {
			log.Println(in.Name, " change fail")
			return
		}
		log.Println(in.Name, " change success")
	})

	if err != nil {
		log.Println("配置读取失败2:", err.Error())
		os.Exit(1)
	}
	log.Println("环境变量初始化成功")

}
