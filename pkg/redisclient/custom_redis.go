package redisclient

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/omygoden/gotools/sfconst"
	"github.com/omygoden/gotools/sfrand"
	"log"
	"time"
	"zgin/pkg/sflogger"
	"zgin/pkg/util"
)

type RedisClient struct {
	Client *redis.Client
}

func NewRedisClient(client *redis.Client) *RedisClient {
	return &RedisClient{
		Client: client,
	}
}

func (this *RedisClient) Set(key string, value interface{}, time time.Duration) string {
	res, err := this.Client.Set(context.Background(), key, value, time).Result()
	if err != nil && err.Error() != "redis: nil" {
		sflogger.RedisErrorLog(util.GetMyFuncName(), err.Error())
		return ""
	}
	return res
}

func (this *RedisClient) Get(key string) string {
	res, err := this.Client.Get(context.Background(), key).Result()
	if err != nil && err.Error() != "redis: nil" {
		sflogger.RedisErrorLog(util.GetMyFuncName(), err.Error())
		return ""
	}
	return res
}

func (this *RedisClient) SetNX(key string, value interface{}, time time.Duration) bool {
	res, err := this.Client.SetNX(context.Background(), key, value, time).Result()
	if err != nil && err.Error() != "redis: nil" {
		sflogger.RedisErrorLog(util.GetMyFuncName(), err.Error())
		return false
	}
	return res
}

func (this *RedisClient) SetNxWait(key string, value interface{}, t time.Duration) bool {
	for {
		if SetNX(key, value, t) {
			return true
		} else {
			time.Sleep(time.Millisecond * 100)
		}
	}
}

func (this *RedisClient) SetEX(key string, value interface{}, time time.Duration) string {
	res, err := this.Client.SetEX(context.Background(), key, value, time).Result()
	if err != nil && err.Error() != "redis: nil" {
		sflogger.RedisErrorLog(util.GetMyFuncName(), err.Error())
		return ""
	}
	return res
}

func (this *RedisClient) SAdd(key string, members ...interface{}) int64 {
	res, err := this.Client.SAdd(context.Background(), key, members...).Result()
	if err != nil && err.Error() != "redis: nil" {
		sflogger.RedisErrorLog(util.GetMyFuncName(), err.Error())
		return 0
	}
	return res
}

func (this *RedisClient) SMembers(key string) []string {
	res, err := this.Client.SMembers(context.Background(), key).Result()
	if err != nil && err.Error() != "redis: nil" {
		sflogger.RedisErrorLog(util.GetMyFuncName(), err.Error())
		return nil
	}
	return res
}

func (this *RedisClient) SIsMember(key string, value interface{}) bool {
	res, err := this.Client.SIsMember(context.Background(), key, value).Result()
	if err != nil && err.Error() != "redis: nil" {
		sflogger.RedisErrorLog(util.GetMyFuncName(), err.Error())
		return false
	}
	return res
}

func (this *RedisClient) SRem(key string, value ...interface{}) int64 {
	res, err := this.Client.SRem(context.Background(), key, value...).Result()
	if err != nil && err.Error() != "redis: nil" {
		sflogger.RedisErrorLog(util.GetMyFuncName(), err.Error())
		return 0
	}
	return res
}

func (this *RedisClient) Llen(key string) int64 {
	res, err := this.Client.LLen(context.Background(), key).Result()
	if err != nil && err.Error() != "redis: nil" {
		sflogger.RedisErrorLog(util.GetMyFuncName(), err.Error())
		return 0
	}
	return res
}

func (this *RedisClient) Lpush(key string, value interface{}) int64 {
	res, err := this.Client.LPush(context.Background(), key, value).Result()
	if err != nil && err.Error() != "redis: nil" {
		sflogger.RedisErrorLog(util.GetMyFuncName(), err.Error())
		return 0
	}
	return res
}

func (this *RedisClient) Rpush(key string, value interface{}) int64 {
	res, err := this.Client.RPush(context.Background(), key, value).Result()
	if err != nil && err.Error() != "redis: nil" {
		sflogger.RedisErrorLog(util.GetMyFuncName(), err.Error())
		return 0
	}
	return res
}

func (this *RedisClient) RPop(key string) string {
	res, err := this.Client.RPop(context.Background(), key).Result()
	if err != nil && err.Error() != "redis: nil" {
		sflogger.RedisErrorLog(util.GetMyFuncName(), err.Error())
		return ""
	}
	return res
}

func (this *RedisClient) LPop(key string) string {
	res, err := this.Client.LPop(context.Background(), key).Result()
	if err != nil && err.Error() != "redis: nil" {
		sflogger.RedisErrorLog(util.GetMyFuncName(), err.Error())
		return ""
	}
	return res
}

func (this *RedisClient) BRPop(key string, waitTime int) string {
	res, err := this.Client.BRPop(context.Background(), time.Duration(waitTime), key).Result()
	if err != nil && err.Error() != "redis: nil" {
		sflogger.RedisErrorLog(util.GetMyFuncName(), err.Error())
		return ""
	}
	if len(res) == 2 {
		return res[1]
	}
	return ""
}

func (this *RedisClient) BLPop(key string, waitTime int) string {
	res, err := this.Client.BLPop(context.Background(), time.Duration(waitTime), key).Result()
	if err != nil && err.Error() != "redis: nil" {
		sflogger.RedisErrorLog(util.GetMyFuncName(), err.Error())
		return ""
	}
	if len(res) == 2 {
		return res[1]
	}
	return ""
}

func (this *RedisClient) HIncrBy(key string, field string, incr int64) int64 {
	res, err := this.Client.HIncrBy(context.Background(), key, field, incr).Result()
	if err != nil && err.Error() != "redis: nil" {
		sflogger.RedisErrorLog(util.GetMyFuncName(), err.Error())
		return 0
	}
	return res
}

func (this *RedisClient) HSet(key string, field string, value interface{}) int64 {
	res, err := this.Client.HSet(context.Background(), key, field, value).Result()
	if err != nil && err.Error() != "redis: nil" {
		sflogger.RedisErrorLog(util.GetMyFuncName(), err.Error())
		return 0
	}
	return res
}

func (this *RedisClient) Incr(key string) int64 {
	res, err := this.Client.Incr(context.Background(), key).Result()
	if err != nil && err.Error() != "redis: nil" {
		sflogger.RedisErrorLog(util.GetMyFuncName(), err.Error())
		return 0
	}
	return res
}

func (this *RedisClient) Ttl(key string) float64 {
	res, err := this.Client.TTL(context.Background(), key).Result()
	if err != nil && err.Error() != "redis: nil" {
		sflogger.RedisErrorLog(util.GetMyFuncName(), err.Error())
		return 0
	}
	return res.Seconds()
}

func (this *RedisClient) Expire(key string, time time.Duration) bool {
	res, err := this.Client.Expire(context.Background(), key, time).Result()
	if err != nil && err.Error() != "redis: nil" {
		sflogger.RedisErrorLog(util.GetMyFuncName(), err.Error())
		return false
	}
	return res
}

func (this *RedisClient) Del(key ...string) int64 {
	res, err := this.Client.Del(context.Background(), key...).Result()
	if err != nil && err.Error() != "redis: nil" {
		sflogger.RedisErrorLog(util.GetMyFuncName(), err.Error())
		return 0
	}
	return res
}

func (this *RedisClient) Select(index int) error {
	_, err := this.Client.Do(context.Background(), "select", index).Result()
	if err != nil && err.Error() != "redis: nil" {
		sflogger.RedisErrorLog(util.GetMyFuncName(), err.Error())
		return err
	}
	return nil
}

//接口请求qps限制
func (this *RedisClient) QpsLimit(fKey string, limitNum int) {
	var key = fmt.Sprintf(fKey, time.Now().Format(sfconst.GO_TIME_WITH_FULL))
	var num = Incr(key)
	if num == 1 {
		go Expire(key, util.GetRedisMdExpire())
	}
	if int(num) > limitNum {
		t := sfrand.RandRange(1, 3)
		log.Println(fmt.Sprintf("key:【%s】，QPS超过%d次，sleep：%ds", key, limitNum, t))
		time.Sleep(time.Second * time.Duration(t))
		QpsLimit(fKey, limitNum)
	}
}
