package customers

import (
	"github.com/streadway/amqp"
	"log"
	"sync"
	"time"
	"zgin/pkg/rabbitmq/rbmqserver"
)

var cFanoutChan *amqp.Channel
var cFanoutChanOnce sync.Once

/**
  QueueCustomerFanout
  @Description: 消费队列--广播模式
  @param exchangeName
  @param qps
  @return error
**/
func QueueCustomerFanout(exchangeName, queueName string, qps int, function func(msg *amqp.Delivery)) error {
	var err error
	cFanoutChanOnce.Do(func() {
		cFanoutChan = rbmqserver.InitFanoutChan(exchangeName, "")
	})

retry:
	//当前消费者一次能接受的最大消息数量
	//服务器传递的最大容量
	//如果为true 对channel可用 false则只对当前队列可用
	_ = cFanoutChan.Qos(qps, 0, false)

	q, err := cFanoutChan.QueueDeclare(
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

	err = cFanoutChan.QueueBind(q.Name, "", exchangeName, false, nil)
	if err != nil {
		log.Println("绑定队列失败：", err.Error())
		return err
	}

	msgs, err := cFanoutChan.Consume(
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

	//如果一直连接不上rabbitmq管道，表示rabbitmq服务已经不可用
	log.Println("customer is close,begin to reconnection rabbitmq channel")
	cFanoutChan = rbmqserver.InitFanoutChan(exchangeName, "")
	time.Sleep(time.Second)
	goto retry
}
