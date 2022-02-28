package producters

import (
	"zgin/global"
	"zgin/pkg/sflogger"
	"github.com/streadway/amqp"
)

const (
	DLXS_EXCHANGE_NAME      = "go.sms.exchange.dlxs"
	DLXS_ROUTE_NAME         = "go.sms.route.dlxs"
	UN_EXISTS_EXCHANGE_NAME = "go.un_exists_exchange.dlxs"
	UN_EXISTS_QUEUE_NAME    = "go.un_exists_queue.dlxs"
	UN_EXISTS_ROUTE_NAME    = "go.un_exists_route.dlxs"
)

/**
  QueueCreateDlxs
  @Description: 创建死信队列--队列过期前未消费，则会进入死信消费者
  @param msg 消息
  @param expire 死信有效期，单位秒
  @return error
**/
func QueueCreateDlxs(dlxsExchangeName, dlxsQueueName, dlxsRouteName, msg string, expire int) error {
	ch, err := global.RabbitmqClient.Channel()
	if err != nil {
		sflogger.QueueAddErrorLog("连接管道", msg, err.Error())
		return err
	}

	err = ch.ExchangeDeclare(dlxsExchangeName, amqp.ExchangeFanout, true, false, false, false, nil)
	if err != nil {
		sflogger.QueueAddErrorLog("创建交换器", msg, err.Error())
		return err
	}
	var (
		args = make(map[string]interface{})
	)
	//设置队列的过期时间
	args["x-message-ttl"] = expire * 1000 //单位为ms
	//设置死信交换器
	args["x-dead-letter-exchange"] = dlxsExchangeName
	//设置死信交换器Key
	args["x-dead-letter-routing-key"] = dlxsRouteName
	normalQueue, err := ch.QueueDeclare(dlxsQueueName, true, false, false, true, args)
	if err != nil {
		sflogger.QueueAddErrorLog("创建队列", msg, err.Error())
		return err
	}

	err = ch.QueueBind(normalQueue.Name, dlxsRouteName, dlxsExchangeName, false, nil)
	if err != nil {
		sflogger.QueueAddErrorLog("绑定队列", msg, err.Error())
		return err
	}

	err = ch.Publish(dlxsExchangeName, dlxsRouteName, false, false, amqp.Publishing{
		ContentType: "text/plain",
		Body:        []byte(msg),
	})
	if err != nil {
		sflogger.QueueAddErrorLog("推入队列", msg, err.Error())
		return nil
	}
	return nil
}
