package dbdriver

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
	"time"
	"zgin/global"
	"zgin/pkg/env"
)

func InitMysqlDriver(maxConn, minConn int) {
	dbConfig := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=%s&parseTime=%t&loc=Local", global.Config.Database.User,
		global.Config.Database.Password,
		global.Config.Database.Host,
		global.Config.Database.Name,
		global.Config.Database.Charset,
		global.Config.Database.ParseTime,
	)

	filepath := fmt.Sprintf("%s/%s/%s", os.Getenv("GOPATH"), env.PROJECT_NAME, global.Config.Database.LogSavePath)
	filename := global.Config.Database.LogFileName
	fileext := global.Config.Database.LogFileExt
	logFile := filepath + "/" + filename + fileext
	file, _ := os.OpenFile(logFile, os.O_CREATE|os.O_APPEND|os.O_RDWR, os.ModeAppend|os.ModePerm)

	newLogger := logger.New(
		log.New(file, "", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Millisecond * time.Duration(global.Config.Database.SlowMinTime), // 慢 SQL 阈值
			LogLevel:                  logger.Info,                                                          // Log level
			IgnoreRecordNotFoundError: true,                                                                 // Ignore ErrRecordNotFound error for logger
			Colorful:                  false,                                                                // Disable color
		},
	)

	db, err := gorm.Open(mysql.New(mysql.Config{
		DSN:        dbConfig,
		DriverName: global.Config.Database.Type,
	}), &gorm.Config{
		Logger: newLogger,
	})
	if err != nil {
		log.Println("数据库初始化失败", err.Error())
		os.Exit(1)
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Println("数据库初始化失败2")
		os.Exit(1)
	}

	sqlDB.SetMaxIdleConns(minConn)
	sqlDB.SetMaxOpenConns(maxConn)
	sqlDB.SetConnMaxLifetime(time.Minute * 10)
	log.Println("数据库初始化成功")

	global.Mysql = db
}
