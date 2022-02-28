package customers

import (
	"zgin/global"
	"github.com/streadway/amqp"
	"log"
)

/**
  QueueCustomerFanout
  @Description: 消费队列--广播模式
  @param exchangeName
  @param qps
  @return error
**/
func QueueCustomerFanout(exchangeName, queueName string, qps int, function func(msg *amqp.Delivery)) error {
	var forever = make(chan int)
	ch, err := global.RabbitmqClient.Channel()
	if err != nil {
		log.Println("连接管道失败：", err.Error())
		return err
	}

	//当前消费者一次能接受的最大消息数量
	//服务器传递的最大容量
	//如果为true 对channel可用 false则只对当前队列可用
	_ = ch.Qos(qps, 0, false)

	err = ch.ExchangeDeclare(exchangeName, amqp.ExchangeFanout, true, false, false, false, nil)
	if err != nil {
		log.Println("创建交换器失败：", err.Error())
		return err
	}

	q, err := ch.QueueDeclare(
		queueName, // name
		true,      // durable
		false,     // delete when unused
		false,     // exclusive
		true,      // no-wait
		nil,
	)
	if err != nil {
		log.Println("创建队列失败：", err.Error())
		return err
	}

	err = ch.QueueBind(q.Name, "", exchangeName, false, nil)
	if err != nil {
		log.Println("绑定队列失败：", err.Error())
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
		log.Println("消费失败：", err.Error())
		return err
	}

	for msg := range msgs {
		go func(message amqp.Delivery) {
			function(&message)
		}(msg)
		//_ = msg.Ack(false)
	}
	log.Println("customer is close")
	<-forever
	return nil
}
