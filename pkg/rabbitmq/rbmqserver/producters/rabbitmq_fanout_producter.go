package producters

import (
	"github.com/streadway/amqp"
	"sync"
	"time"
	"zgin/global"
	"zgin/pkg/sflogger"
)

var fanoutChan *amqp.Channel
var fanoutChanOnce sync.Once

/**
  QueueCreateFanout
  @Description: 创建队列--广播模式
  @param exchangeName
  @param queueName
  @param msg
  @return error
**/
func QueueCreateFanout(exchangeName, queueName, msg string) error {
	var err error
	var retryNum = 3
	fanoutChanOnce.Do(func() {
		initFanoutChan(exchangeName, msg)
	})

retry:
	err = fanoutChan.Publish(exchangeName, queueName, false, false, amqp.Publishing{
		ContentType: "text/plain",
		Body:        []byte(msg),
	})

	if err != nil {
		sflogger.QueueAddErrorLog("推入队列,重新连接信道", msg, err.Error())
		if retryNum > 0 {
			retryNum--
			initFanoutChan(exchangeName, msg)
			time.Sleep(time.Second)
			goto retry
		}
		return err
	}

	return nil
}

func initFanoutChan(exchangeName, msg string) {
	ch, err := global.RabbitmqClient.Channel()
	if err != nil {
		sflogger.QueueAddErrorLog("连接管道", msg, err.Error())
		return
	}
	err = ch.ExchangeDeclare(exchangeName, amqp.ExchangeFanout, true, false, false, false, nil)
	if err != nil {
		sflogger.QueueAddErrorLog("创建交换器", msg, err.Error())
		return
	}
	fanoutChan = ch
}
