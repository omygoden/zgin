package producters

import (
	"github.com/streadway/amqp"
	"zgin/global"
	"zgin/pkg/sflogger"
)

/**
  QueueCreateDirect
  @Description: 创建队列 -- 直连模式
  @param exchangeName
  @param queueName
  @param routeName
  @param msg
  @return error
**/
func QueueCreateDirect(exchangeName, queueName, routeName, msg string) error {
	ch, err := global.RabbitmqClient.Channel()
	if err != nil {
		sflogger.QueueAddErrorLog("连接管道", msg, err.Error())
		return err
	}

	err = ch.ExchangeDeclare(exchangeName, amqp.ExchangeDirect, true, false, false, false, nil)
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

	//_ = ch.Tx()
	err = ch.Publish(exchangeName, routeName, false, false, amqp.Publishing{
		ContentType: "text/plain",
		Body:        []byte(msg),
	})
	if err != nil {
		//_ = ch.TxRollback()
		sflogger.QueueAddErrorLog("推入队列", msg, err.Error())
		return err
	}
	//_ = ch.TxCommit()
	return nil
}
