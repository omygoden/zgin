package redisclient

import (
	"context"
	"fmt"
	"github.com/omygoden/gotools/sfconst"
	"github.com/omygoden/gotools/sfrand"
	"log"
	"time"
	"zgin/global"
	"zgin/pkg/sflogger"
	"zgin/pkg/util"
)

/**
默认第0个库
*/

func Set(key string, value interface{}, time time.Duration) string {
	res, err := global.RedisClient.Set(context.Background(), key, value, time).Result()
	if err != nil && err.Error() != "redis: nil" {
		sflogger.RedisErrorLog(util.GetMyFuncName(), err.Error())
		return ""
	}
	return res
}

func Get(key string) string {
	res, err := global.RedisClient.Get(context.Background(), key).Result()
	if err != nil && err.Error() != "redis: nil" {
		sflogger.RedisErrorLog(util.GetMyFuncName(), err.Error())
		return ""
	}
	return res
}

func SetNX(key string, value interface{}, time time.Duration) bool {
	res, err := global.RedisClient.SetNX(context.Background(), key, value, time).Result()
	if err != nil && err.Error() != "redis: nil" {
		sflogger.RedisErrorLog(util.GetMyFuncName(), err.Error())
		return false
	}
	return res
}

func SetNxWait(key string, value interface{}, t time.Duration) bool {
	for {
		if SetNX(key, value, t) {
			return true
		} else {
			time.Sleep(time.Millisecond * 100)
		}
	}
}

func SetEX(key string, value interface{}, time time.Duration) string {
	res, err := global.RedisClient.SetEX(context.Background(), key, value, time).Result()
	if err != nil && err.Error() != "redis: nil" {
		sflogger.RedisErrorLog(util.GetMyFuncName(), err.Error())
		return ""
	}
	return res
}

func SAdd(key string, members ...interface{}) int64 {
	res, err := global.RedisClient.SAdd(context.Background(), key, members...).Result()
	if err != nil && err.Error() != "redis: nil" {
		sflogger.RedisErrorLog(util.GetMyFuncName(), err.Error())
		return 0
	}
	return res
}

func SMembers(key string) []string {
	res, err := global.RedisClient.SMembers(context.Background(), key).Result()
	if err != nil && err.Error() != "redis: nil" {
		sflogger.RedisErrorLog(util.GetMyFuncName(), err.Error())
		return nil
	}
	return res
}

func SIsMember(key string, value interface{}) bool {
	res, err := global.RedisClient.SIsMember(context.Background(), key, value).Result()
	if err != nil && err.Error() != "redis: nil" {
		sflogger.RedisErrorLog(util.GetMyFuncName(), err.Error())
		return false
	}
	return res
}

func SRem(key string, value ...interface{}) int64 {
	res, err := global.RedisClient.SRem(context.Background(), key, value...).Result()
	if err != nil && err.Error() != "redis: nil" {
		sflogger.RedisErrorLog(util.GetMyFuncName(), err.Error())
		return 0
	}
	return res
}

func Llen(key string) int64 {
	res, err := global.RedisClient.LLen(context.Background(), key).Result()
	if err != nil && err.Error() != "redis: nil" {
		sflogger.RedisErrorLog(util.GetMyFuncName(), err.Error())
		return 0
	}
	return res
}

func Lpush(key string, value interface{}) int64 {
	res, err := global.RedisClient.LPush(context.Background(), key, value).Result()
	if err != nil && err.Error() != "redis: nil" {
		sflogger.RedisErrorLog(util.GetMyFuncName(), err.Error())
		return 0
	}
	return res
}

func Rpush(key string, value interface{}) int64 {
	res, err := global.RedisClient.RPush(context.Background(), key, value).Result()
	if err != nil && err.Error() != "redis: nil" {
		sflogger.RedisErrorLog(util.GetMyFuncName(), err.Error())
		return 0
	}
	return res
}

func RPop(key string) string {
	res, err := global.RedisClient.RPop(context.Background(), key).Result()
	if err != nil && err.Error() != "redis: nil" {
		sflogger.RedisErrorLog(util.GetMyFuncName(), err.Error())
		return ""
	}
	return res
}

func LPop(key string) string {
	res, err := global.RedisClient.LPop(context.Background(), key).Result()
	if err != nil && err.Error() != "redis: nil" {
		sflogger.RedisErrorLog(util.GetMyFuncName(), err.Error())
		return ""
	}
	return res
}

func BRPop(key string, waitTime int) string {
	res, err := global.RedisClient.BRPop(context.Background(), time.Duration(waitTime), key).Result()
	if err != nil && err.Error() != "redis: nil" {
		sflogger.RedisErrorLog(util.GetMyFuncName(), err.Error())
		return ""
	}
	if len(res) == 2 {
		return res[1]
	}
	return ""
}

func BLPop(key string, waitTime int) string {
	res, err := global.RedisClient.BLPop(context.Background(), time.Duration(waitTime), key).Result()
	if err != nil && err.Error() != "redis: nil" {
		sflogger.RedisErrorLog(util.GetMyFuncName(), err.Error())
		return ""
	}
	if len(res) == 2 {
		return res[1]
	}
	return ""
}

func HIncrBy(key string, field string, incr int64) int64 {
	res, err := global.RedisClient.HIncrBy(context.Background(), key, field, incr).Result()
	if err != nil && err.Error() != "redis: nil" {
		sflogger.RedisErrorLog(util.GetMyFuncName(), err.Error())
		return 0
	}
	return res
}

func HSet(key string, field string, value interface{}) int64 {
	res, err := global.RedisClient.HSet(context.Background(), key, field, value).Result()
	if err != nil && err.Error() != "redis: nil" {
		sflogger.RedisErrorLog(util.GetMyFuncName(), err.Error())
		return 0
	}
	return res
}

func Incr(key string) int64 {
	res, err := global.RedisClient.Incr(context.Background(), key).Result()
	if err != nil && err.Error() != "redis: nil" {
		sflogger.RedisErrorLog(util.GetMyFuncName(), err.Error())
		return 0
	}
	return res
}

func Ttl(key string) float64 {
	res, err := global.RedisClient.TTL(context.Background(), key).Result()
	if err != nil && err.Error() != "redis: nil" {
		sflogger.RedisErrorLog(util.GetMyFuncName(), err.Error())
		return 0
	}
	return res.Seconds()
}

func Expire(key string, time time.Duration) bool {
	res, err := global.RedisClient.Expire(context.Background(), key, time).Result()
	if err != nil && err.Error() != "redis: nil" {
		sflogger.RedisErrorLog(util.GetMyFuncName(), err.Error())
		return false
	}
	return res
}

func Del(key ...string) int64 {
	res, err := global.RedisClient.Del(context.Background(), key...).Result()
	if err != nil && err.Error() != "redis: nil" {
		sflogger.RedisErrorLog(util.GetMyFuncName(), err.Error())
		return 0
	}
	return res
}

func Select(index int) error {
	_, err := global.RedisClient.Do(context.Background(), "select", index).Result()
	if err != nil && err.Error() != "redis: nil" {
		sflogger.RedisErrorLog(util.GetMyFuncName(), err.Error())
		return err
	}
	return nil
}

//接口请求qps限制
func QpsLimit(fKey string, limitNum int) {
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
