package dbdriver

import (
	"context"
	"errors"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/utils"
	"log"
	"os"
	"time"
	"zgin/global"
	"zgin/pkg/sflogger"
)

type mysqlLogger struct {
	logger.Writer
	logger.Config
	infoStr, warnStr, errStr            string
	traceStr, traceErrStr, traceWarnStr string
}

type myWriter struct {
}

func InitMysqlDriver(maxConn, minConn int) {
	if global.Config.Database.Host != "" {
		global.Mysql = InitMysql(global.Config.Database, maxConn, minConn)
	}
}

func InitMysql(DatabaseConfig global.DatabaseConfig, maxConn, minConn int) *gorm.DB {
	dbConfig := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=%t&loc=Local&multiStatements=true",
		DatabaseConfig.User,
		DatabaseConfig.Password,
		DatabaseConfig.Host,
		DatabaseConfig.Port,
		DatabaseConfig.Name,
		DatabaseConfig.Charset,
		DatabaseConfig.ParseTime,
	)

	if maxConn == 0 || minConn == 0 {
		maxConn = DatabaseConfig.MaxOpenConns
		minConn = DatabaseConfig.MaxIdleConns
	}
	/*
	* 1默认，2错误，3警告，4所有
	* 2：只打印错误日志
	* 3：只打印错误和慢日志
	* 4：每次执行都打印
	* 生产环境设置为3即可
	 */
	var logLevel = logger.Info
	if global.Config.App.Env == "production" {
		logLevel = logger.Warn
	}

	newLogger := &mysqlLogger{
		Writer: &myWriter{},
		Config: logger.Config{
			SlowThreshold:             time.Millisecond * time.Duration(DatabaseConfig.SlowMinTime), // 慢 SQL 阈值
			LogLevel:                  logLevel,                                                     // Log level
			IgnoreRecordNotFoundError: true,                                                         // Ignore ErrRecordNotFound error for logger
			Colorful:                  false,                                                        // Disable color
		},
	}

	db, err := gorm.Open(mysql.New(mysql.Config{
		DSN:        dbConfig,
		DriverName: DatabaseConfig.Type,
	}), &gorm.Config{
		Logger: newLogger,
	})
	if err != nil {
		log.Println(fmt.Sprintf("[%s]mysql初始化失败:%s", DatabaseConfig.Host, err.Error()))
		os.Exit(1)
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Println(fmt.Sprintf("[%s]mysql链接失败:%s", DatabaseConfig.Host, err.Error()))
		os.Exit(1)
	}

	sqlDB.SetMaxIdleConns(minConn)
	sqlDB.SetMaxOpenConns(maxConn)
	sqlDB.SetConnMaxLifetime(time.Minute * 10)
	log.Println(fmt.Sprintf("[%s]mysql初始化成功", DatabaseConfig.Host))

	return db
}

func (w *myWriter) Printf(format string, args ...interface{}) {
	log.Println(fmt.Sprintf(format, args...))
}

// LogMode log mode
func (l *mysqlLogger) LogMode(level logger.LogLevel) logger.Interface {
	newlogger := *l
	newlogger.LogLevel = level
	return &newlogger
}

// Info print info
func (l mysqlLogger) Info(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= logger.Info {
		l.Printf(l.infoStr+msg, append([]interface{}{utils.FileWithLineNum()}, data...)...)
	}
}

// Warn print warn messages
func (l mysqlLogger) Warn(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= logger.Warn {
		l.Printf(l.warnStr+msg, append([]interface{}{utils.FileWithLineNum()}, data...)...)
	}
}

// Error print error messages
func (l mysqlLogger) Error(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= logger.Error {
		l.Printf(l.errStr+msg, append([]interface{}{utils.FileWithLineNum()}, data...)...)
	}
}

// Trace print sql message
func (l mysqlLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	if l.LogLevel <= logger.Silent {
		return
	}

	elapsed := time.Since(begin)
	switch {
	case err != nil && l.LogLevel >= logger.Error && (!errors.Is(err, logger.ErrRecordNotFound) || !l.IgnoreRecordNotFoundError):
		sql, rows := fc()
		go sflogger.MysqlErrorLog(utils.FileWithLineNum(), err.Error(), sql, float64(elapsed.Nanoseconds())/1e6, rows)
	case elapsed > l.SlowThreshold && l.SlowThreshold != 0 && l.LogLevel >= logger.Warn:
		sql, rows := fc()
		go sflogger.MysqlSlowLog(utils.FileWithLineNum(), sql, float64(elapsed.Nanoseconds())/1e6, rows)
	case l.LogLevel == logger.Info:
		if global.Config.App.Debug {
			sql, rows := fc()
			go sflogger.MysqlLog(utils.FileWithLineNum(), sql, float64(elapsed.Nanoseconds())/1e6, rows)
		}
	}
}
