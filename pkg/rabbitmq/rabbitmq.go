package rabbitmq

import (
	"fmt"
	"github.com/streadway/amqp"
	"log"
	"os"
	"time"
	"zgin/global"
)

func InitRabbitmq() {
	rabbitmqClient, err := amqp.DialConfig(fmt.Sprintf("amqp://%s:%s@%s", global.Config.Rabbitmq.RabbitmqName, global.Config.Rabbitmq.RabbitmqPwd, global.Config.Rabbitmq.RabbitmqHost), amqp.Config{
		ChannelMax: global.Config.Rabbitmq.RabbitmqMaxChannel,
		Heartbeat:  time.Second,
	})
	if err != nil {
		log.Println("rabbitmq初始化失败：", err.Error())
		os.Exit(1)
	}

	global.RabbitmqClient = rabbitmqClient
	//initChannalPool()
	go connCloseListen()

	//初始化信道
	//if needInitChan {
	//	minChannalLen := global.Config.Rabbitmq.RabbitmqMaxChannel/2
	//	global.RabbitmqChannal = make(chan *amqp.Channel,global.Config.Rabbitmq.RabbitmqMaxChannel)
	//	for i:=0;i<minChannalLen;i++ {
	//		ch,err := global.RabbitmqClient.Channel()
	//		if err != nil {
	//			log.Println("rabbitmq初始化信道失败：",err.Error())
	//			os.Exit(1)
	//		}
	//		go chanCloseListen(ch)
	//		global.RabbitmqChannal <- ch
	//	}
	//}

	//global.RabbitmqClient.Config.ChannelMax = 10 //动态设置管道数，管道数对应线程数
	log.Println("rabbitmq初始化成功")
}

//监听rabbitmq的连接关闭情况
func connCloseListen() {
	closeNotify := global.RabbitmqClient.NotifyClose(make(chan *amqp.Error))
	for {
		select {
		//收到关闭通知，重新建立连接,并且return当前老的连接监听协程
		case res := <-closeNotify:
			log.Println("RabbitqMQ连接关闭通知，触发重连接机制", res.Error(), res.Reason)
			InitRabbitmq()
			return
		}
	}
}

//信道池
func initChannalPool() {
	global.RabbitmqChannalPool = make(chan int, global.Config.Rabbitmq.RabbitmqMaxChannel)
	for i := 0; i < global.Config.Rabbitmq.RabbitmqMaxChannel; i++ {
		global.RabbitmqChannalPool <- i
	}
}

//监听rabbitmq的信道关闭情况
//func chanCloseListen(ch *amqp.Channel)  {
//	closeNotify := ch.NotifyClose(make(chan *amqp.Error))
//	for {
//		select {
//		//收到关闭通知，新建信道并推入信道池,并且return当前信道的监控协程
//		case res := <- closeNotify:
//			log.Println("RabbitqMQ信道关闭通知，触发重连接机制",res.Error(),res.Reason)
//			newCh,err := global.RabbitmqClient.Channel()
//			if err != nil {
//				log.Println("rabbitmq初始化信道失败：",err.Error())
//				return
//			}
//			global.RabbitmqChannal <- newCh
//			go chanCloseListen(newCh)
//			log.Println("RabbitqMQ信道新建成功")
//			return
//		}
//	}
//}
