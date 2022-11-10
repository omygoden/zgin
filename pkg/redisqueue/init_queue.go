package redisqueue

import (
	"context"
	"github.com/go-redis/redis/v8"
	"log"
	"os"
	"zgin/global"
)

const exists_err = "BUSYGROUP Consumer Group name already exists"

func InitRedisStreamQueue(streamName, groupName string) {
	/* 创建消费组 */
	res := global.RedisLocalClient.XGroupCreateMkStream(context.TODO(), streamName, groupName, "0")
	if res.Err() != nil && res.Err().Error() != exists_err {
		log.Println("队列消费组创建失败,失败原因:", res.Err().Error())
		os.Exit(1)
	}

	/* 重新唤起xpending里未消费的队列，主要针对单进程，否则会导致消息重复消费 */
	/* 先获取xpending里的数据 */
	xPArgs := redis.XPendingExtArgs{
		Stream: streamName,
		Group:  groupName,
		Idle:   0,
		Start:  "-",
		End:    "+",
		Count:  100,
	}
	var num int
	var startId string
	for {
		pendingLists := global.RedisLocalClient.XPendingExt(context.TODO(), &xPArgs)
		list, _ := pendingLists.Result()
		if len(list) == 0 {
			break
		}
		log.Println("待重写推入队列消息数量:", len(list))
		for _, v := range list {
			/*  先获取消息体 */
			msg := getStreamMsgById(streamName, v.ID)
			if msg != nil {
				/* 重写推入队列 */
				r := global.RedisLocalClient.XAdd(context.TODO(), &redis.XAddArgs{
					Stream: streamName,
					Values: msg,
				})
				if r.Err() != nil {
					log.Println("重写推入队列失败，失败原因:", r.Err().Error())
					os.Exit(1)
				}
				/* 推入队列成功删除原ID对应数据 */
				global.RedisLocalClient.XAck(context.TODO(), streamName, groupName, v.ID)
				global.RedisLocalClient.XDel(context.TODO(), streamName, v.ID)
				num++
			}
			startId = v.ID
		}
		xPArgs.Start = startId
	}
	log.Println("成功推入队列:", num)

	log.Println("redis-stream队列初始化成功")
}

func getStreamMsgById(streamName, id string) map[string]interface{} {
	res := global.RedisLocalClient.XRange(context.TODO(), streamName, id, id)
	r, _ := res.Result()
	for _, v := range r {
		if v.ID == id {
			return v.Values
		}
	}
	return nil
}