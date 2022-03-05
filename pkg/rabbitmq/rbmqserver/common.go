package rbmqserver

import (
	"github.com/streadway/amqp"
	"zgin/global"
	"zgin/pkg/sflogger"
)

func InitFanoutChan(exchangeName, msg string) *amqp.Channel {
	ch, err := global.RabbitmqClient.Channel()
	if err != nil {
		sflogger.QueueAddErrorLog("连接管道", msg, err.Error())
		return nil
	}
	err = ch.ExchangeDeclare(exchangeName, amqp.ExchangeFanout, true, false, false, false, nil)
	if err != nil {
		sflogger.QueueAddErrorLog("创建交换器", msg, err.Error())
		return nil
	}
	return ch
}
