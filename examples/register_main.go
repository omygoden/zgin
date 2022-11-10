package main

import (
	"log"
	"time"
	"zgin/microservice"
)

func main()  {
	ser,_ := microservice.NewRegisterService([]string{"127.0.0.1:2379"},5)
	key := "test/test"
	key2 := "test/test2"
	err := ser.PutService(key,"测试服务")
	err = ser.PutService(key2,"测试服务-2")
	log.Println(err)
	log.Println("start success")
	ser.GetService(key)
	time.Sleep(time.Second * 10)
	ser.PutService(key+"3","测试服务-3")

	//err = ser.RevokeLease()
	//log.Println("撤销租赁：",err)
	select {

	}
}
