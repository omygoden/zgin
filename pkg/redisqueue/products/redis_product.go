package products

import (
	"context"
	"encoding/json"
	"github.com/go-redis/redis/v8"
	"sph-banjia-go/global"
	"sph-banjia-go/pkg/redisqueue"
	"sph-banjia-go/pkg/sflogger"
)

func RedisStreamProduct(streamName string, msg map[string]interface{}) error {
	res := global.RedisLocalClient.XAdd(context.TODO(), &redis.XAddArgs{
		Stream: streamName,
		Values: msg,
		MaxLen: redisqueue.QUEUE_MAX_LEN,
	})
	if res.Err() != nil {
		b, _ := json.Marshal(msg)
		sflogger.QueueAddErrorLog("推入队列失败", string(b), res.Err().Error())
		return res.Err()
	}
	return nil
}