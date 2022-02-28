package redisclient

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
	"zgin/global"
)

var ctx = context.Background()

//根据自定义的库，连接指定的redis库
func InitRedisPool(maxConn, minConn int) {
	var redisIndexs = strings.Split(global.Config.Redis.DbIndexs, ",")
	var indexs string
	var redisDbs = make([]int, len(redisIndexs))

	for k, v := range redisIndexs {
		i, _ := strconv.Atoi(v)
		redisDbs[k] = i
	}

	global.RedisClients = make(map[int]*redis.Client)
	for _, k := range redisDbs {
		global.RedisClients[k] = initRedis(maxConn, minConn, k)
		indexs += fmt.Sprintf("%d,", k)
	}

	if _, ok := global.RedisClients[0]; !ok {
		global.RedisClient = initRedis(maxConn, minConn, 0)
		indexs += fmt.Sprintf("%d,", 0)
	}

	log.Println(fmt.Sprintf("初始化redis库成功--共%d库，分别是:%s", len(redisIndexs), strings.TrimRight(indexs, ",")))
}

func initRedis(maxConn, minConn, index int) *redis.Client {
	client := redis.NewClient(&redis.Options{
		//连接信息
		Network:  "tcp",                                                                    //网络类型，tcp or unix，默认tcp
		Addr:     fmt.Sprintf("%s:%d", global.Config.Redis.Host, global.Config.Redis.Port), //主机名+冒号+端口，默认localhost:6379，正式：10.23.68.87
		Password: global.Config.Redis.Password,                                             //密码
		DB:       index,                                                                    // redis数据库index

		//连接池容量及闲置连接数量
		PoolSize:     maxConn, // 连接池最大socket连接数，默认为4倍CPU数， 4 * runtime.NumCPU
		MinIdleConns: minConn, //在启动阶段创建指定数量的Idle连接，并长期维持idle状态的连接数不少于指定数量；。

		//超时
		DialTimeout:  5 * time.Second, //连接建立超时时间，默认5秒。
		ReadTimeout:  3 * time.Second, //读超时，默认3秒， -1表示取消读超时
		WriteTimeout: 3 * time.Second, //写超时，默认等于读超时
		PoolTimeout:  4 * time.Second, //当所有连接都处在繁忙状态时，客户端等待可用连接的最大等待时长，默认为读超时+1秒。

		//闲置连接检查包括IdleTimeout，MaxConnAge
		IdleCheckFrequency: 60 * time.Second, //闲置连接检查的周期，默认为1分钟，-1表示不做周期性检查，只在客户端获取连接时对闲置连接进行处理。
		IdleTimeout:        1 * time.Minute,  //闲置超时，默认5分钟，-1表示取消闲置超时检查
		MaxConnAge:         0 * time.Second,  //连接存活时长，从创建开始计时，超过指定时长则关闭连接，默认为0，即不关闭存活时长较长的连接

		//命令执行失败时的重试策略
		MaxRetries:      1,                      // 命令执行失败时，最多重试多少次，默认为0即不重试
		MinRetryBackoff: 8 * time.Millisecond,   //每次计算重试间隔时间的下限，默认8毫秒，-1表示取消间隔
		MaxRetryBackoff: 512 * time.Millisecond, //每次计算重试间隔时间的上限，默认512毫秒，-1表示取消间隔

		//钩子函数，如果连接池需要新建连接时则会调用此钩子函数
		//OnConnect: func(ctx context.Context, conn *redis.Conn) error {
		//	fmt.Printf("conn=%v\n", conn)
		//	return nil
		//},
	})
	if err := client.Ping(ctx).Err(); err != nil {
		log.Println(fmt.Sprintf("初始化redis失败: %s", err))
		os.Exit(1)
	}
	return client
}
