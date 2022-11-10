package global

import (
	ut "github.com/go-playground/universal-translator"
	redis2 "github.com/go-redis/redis/v8"
	"github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
	"go.mongodb.org/mongo-driver/mongo"
	"gorm.io/gorm"
	"sync"
)

var (
	Config              allConfig
	Mysql               *gorm.DB
	MysqlSub            *gorm.DB // 分表
	RabbitmqClient      *amqp.Connection
	MongoClient         *mongo.Client
	Logger              *logrus.Logger
	Trans               ut.Translator
	GoroutinePool       chan int
	RabbitmqChannalPool chan int
)

var (
	RedisClients     map[int]*redis2.Client
	RedisClient      *redis2.Client
	RedisLocalClient *redis2.Client
)

var (
	LoggerLock sync.Mutex
)

type allConfig struct {
	App         app
	Database    DatabaseConfig
	DatabaseSub DatabaseConfig
	Redis       RedisConfig
	RedisLocal  RedisConfig
	Sms         smsConfig
	Rabbitmq    rabbitmqConfig
	Goroutine   goroutineConfig
}

type app struct {
	Env         string
	Host        string
	PhpHost     string
	HttpPort    string
	PprofPort   string
	LogSavePath string
	LogFileExt  string
}

type DatabaseConfig struct {
	Type         string
	User         string
	Password     string
	Host         string
	Port         string
	Name         string
	TablePrefix  string
	Charset      string
	ParseTime    bool
	MaxIdleConns int
	MaxOpenConns int
	LogSavePath  string
	LogFileName  string
	LogFileExt   string
	SlowMinTime  int // 慢查询时间，单位毫秒，大于该时间则算是慢查询
}

type RedisConfig struct {
	Host        string
	Password    string
	Port        int
	DbIndexs    string // 允许初始化多个库,用逗号隔开
	DbIndex     int
	MaxIdle     int
	MaxActive   int
	IdleTimeout int
	MinIdleConn int
	MaxConn     int
}

// 配置后期可考虑加入到数据库里查询获取
type smsConfig struct {
	IsSend   bool
	Domain   string
	Account  string
	Password string
}

type rabbitmqConfig struct {
	RabbitmqHost       string
	RabbitmqPwd        string
	RabbitmqName       string
	RabbitmqMaxChannel int
}

type goroutineConfig struct {
	MaxPullGoruntine int
}
