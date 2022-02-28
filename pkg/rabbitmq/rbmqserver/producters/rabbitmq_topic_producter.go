package producters

import (
	"zgin/global"
	"zgin/pkg/sflogger"
	"github.com/streadway/amqp"
)

/**
  QueueCreateTopic
  @Description: 创建队列 -- topic模式
  @param exchangeName
  @param routeName 路由名称，消费者通过route.*匹配即可消费
  @return error
**/
func QueueCreateTopic(exchangeName, queueName, routeName, msg string) error {
	ch, err := global.RabbitmqClient.Channel()
	if err != nil {
		sflogger.QueueAddErrorLog("连接管道", msg, err.Error())
		return err
	}

	err = ch.ExchangeDeclare(exchangeName, amqp.ExchangeTopic, true, false, false, false, nil)
	if err != nil {
		sflogger.QueueAddErrorLog("创建交换器", msg, err.Error())
		return err
	}

	normalQueue, err := ch.QueueDeclare(queueName, true, false, false, false, nil)
	if err != nil {
		sflogger.QueueAddErrorLog("创建队列", msg, err.Error())
		return err
	}

	err = ch.QueueBind(normalQueue.Name, routeName, exchangeName, false, nil)
	if err != nil {
		sflogger.QueueAddErrorLog("绑定队列", msg, err.Error())
		return err
	}

	err = ch.Publish(exchangeName, routeName, false, false, amqp.Publishing{
		ContentType: "text/plain",
		Body:        []byte(msg),
	})
	if err != nil {
		sflogger.QueueAddErrorLog("推入队列", msg, err.Error())
		return nil
	}

	return nil

}
