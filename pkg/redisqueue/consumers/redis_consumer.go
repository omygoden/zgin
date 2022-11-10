package consumers

import (
	"context"
	"github.com/go-redis/redis/v8"
	"log"
	"sph-banjia-go/global"
	"time"
)

/*
 	todo stream队列缺点
	1、持久化还是和其他类型一样，走的aof和rdb
	2、不支持延迟队列，如果要使用还是需要配合sort set
	3、消息ack后，只是删除了pending里的数据，实际消息还在，需要del之后才能真正删除
	4、队列异常退出，部分没有ack的消息，重新启动后无法自动重新消费
	5、因为ack之后消息并没有立即删除，所以很容易造成堆积
*/
func RedisStreamConsumers(streamName, groupName, consumerName string, function func(data map[string]interface{}) error) {
	var ch = make(chan int, 1000) // 限制消费速度
	go func() {
		/* 监控消息队列长度，避免造成堆积 */
		for {
			l, _ := global.RedisLocalClient.XLen(context.TODO(), streamName).Result()
			log.Println("队列名称:", streamName, ";当前消息数量:", l)
			time.Sleep(time.Second * 5)
		}
	}()
	for {
		res := global.RedisLocalClient.XReadGroup(context.TODO(), &redis.XReadGroupArgs{
			Group:    groupName,
			Consumer: consumerName,
			Streams:  []string{streamName, ">"},
			// Block: -1,
		})
		s, _ := res.Result()
		for _, v := range s {
			for _, vv := range v.Messages {
				ch <- 1
				go func(vv redis.XMessage, ch chan int) {
					defer func() {
						<-ch
					}()
					err := function(vv.Values)
					if err != nil {
						log.Println("消费队列处理失败", ";消息ID:", vv.ID, ";失败原因:", err.Error())
						return
					}
					/* 消费成功删除对应ID */
					global.RedisLocalClient.XAck(context.TODO(), streamName, groupName, vv.ID)
					global.RedisLocalClient.XDel(context.TODO(), streamName, vv.ID)
				}(vv, ch)
			}
		}
	}
}
