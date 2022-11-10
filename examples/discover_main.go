package main

import (
	"log"
	"zgin/microservice"
)

func main()  {
	dis,err := microservice.NewDiscoverService([]string{"127.0.0.1:2379"})
	if err != nil {
		log.Println(err)
		return
	}
	addrs,err := dis.GetService("test")
	if err != nil {
		log.Println(err)
		return
	}
	log.Println(addrs)
	select {

	}
}
