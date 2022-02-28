package customers

import (
	"github.com/streadway/amqp"
	"zgin/global"
)

/**
  QueueCustomerDirect
  @Description: 消费队列 -- 直连模式
  @param exchangeName
  @param queueName
  @param routeName
  @param qps qps限制
  @return error
**/
func QueueCustomerDirect(exchangeName, queueName, routeName string, qps int, function func(msg amqp.Delivery)) error {
	ch, err := global.RabbitmqClient.Channel()
	if err != nil {
		return err
	}

	//当前消费者一次能接受的最大消息数量
	//服务器传递的最大容量
	//如果为true 对channel可用 false则只对当前队列可用
	_ = ch.Qos(qps, 0, false)

	err = ch.ExchangeDeclare(exchangeName, amqp.ExchangeDirect, true, false, false, false, nil)
	if err != nil {
		return err
	}

	q, err := ch.QueueDeclare(
		queueName, // name
		false,     // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,
	)
	if err != nil {
		return err
	}

	err = ch.QueueBind(q.Name, routeName, exchangeName, false, nil)
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

	go func() {
		for msg := range msgs {
			function(msg)
			//_ = msg.Ack(false) //false表示响应成功

			//	d.Reject(false)//todo 拒绝消息，具体结果待确定
		}
	}()
	return nil
}
