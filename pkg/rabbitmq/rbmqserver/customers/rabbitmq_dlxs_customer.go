package customers

import (
	"zgin/global"
	"github.com/streadway/amqp"
)

/**
  QueueCustomerDlxs
  @Description: 死信队列消费者 -- 即延迟队列
  @param dlxsExchangeName
  @param dlxsQueueName
  @param dlxsRouteName
  @param qps
  @param function
  @return error
**/
func QueueCustomerDlxs(dlxsExchangeName, dlxsQueueName, dlxsRouteName string, qps int, function func(msg amqp.Delivery)) error {
	ch, err := global.RabbitmqClient.Channel()
	if err != nil {
		return err
	}

	//当前消费者一次能接受的最大消息数量
	//服务器传递的最大容量
	//如果为true 对channel可用 false则只对当前队列可用
	_ = ch.Qos(qps, 0, false)

	err = ch.ExchangeDeclare(dlxsExchangeName, amqp.ExchangeDirect, true, false, false, false, nil)
	if err != nil {
		return err
	}
	q, err := ch.QueueDeclare(
		dlxsQueueName, // name
		false,         // durable
		false,         // delete when unused
		false,         // exclusive
		true,          // no-wait
		nil,
	)
	if err != nil {
		return err
	}

	err = ch.QueueBind(
		q.Name,
		dlxsRouteName,
		dlxsExchangeName,
		false,
		nil,
	)
	if err != nil {
		return err
	}
	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		false,  // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)

	if err != nil {
		return err
	}

	for msg := range msgs {
		function(msg)
	}
	return nil
}
