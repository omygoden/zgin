package setting

import (
	"zgin/global"
	"zgin/pkg/dbdriver"
	"zgin/pkg/env"
	"zgin/pkg/formtranslator"
	"zgin/pkg/goroutinepool"
	settings "zgin/pkg/rabbitmq"
	"zgin/pkg/redisclient"
	"zgin/pkg/sflogger"
)

func InitSetting() {
	setupEnv()
	//setupDb(global.Config.Database.MaxOpenConns, global.Config.Database.MaxIdleConns)
	//setupRedis(global.Config.Redis.MaxConn, global.Config.Redis.MinIdleConn)
	setupLogger()
	setupTrans()
	//setupGoruntimePool()
	//setupRabbitmq()
	goroutinepool.GoroutineListen()
}

func QueueInitSettting(dbMaxConn, dbMinConn int) {
	setupEnv()
	setupDb(dbMaxConn, dbMinConn)
	setupRedis(global.Config.Redis.MaxConn, global.Config.Redis.MinIdleConn)
	setupLogger()
	setupGoruntimePool()
	setupRabbitmq()
}

func CronInitSetting(dbMaxConn, dbMinConn int) {
	setupEnv()
	setupDb(dbMaxConn, dbMinConn)
	setupRedis(global.Config.Redis.MaxConn, global.Config.Redis.MinIdleConn)
	setupLogger()
	setupGoruntimePool()
}

//初始化环境变量
func setupEnv() {
	env.InitEnv("config")
}

//初始化数据库
func setupDb(maxConn, minConn int) {
	dbdriver.InitMysqlDriver(maxConn, minConn)
}

//初始化redis连接池
func setupRedis(maxConn, minConn int) {
	redisclient.InitRedisPool(maxConn, minConn)
}

//初始化日志配置
func setupLogger() {
	sflogger.InitLogger()
}

//初始化翻译器
func setupTrans() {
	formtranslator.InitTrans()
}

//初始化rabbitmq
func setupRabbitmq() {
	settings.InitRabbitmq()
}

//初始化协程池
func setupGoruntimePool() {
	goroutinepool.InitGoroutinePool()
}
