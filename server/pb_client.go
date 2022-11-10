package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"log"
)
import pb "zgin/pb"

func main()  {
	var conn *grpc.ClientConn
	var err error

	client := pb.NewGreeterClient(conn)

	resp,err := client.SayHello(context.Background(),&pb.HelloRequest{Name: "夏利"})
	log.Println("error2:",err)

	log.Println(fmt.Sprintf("%s",resp.GetMessage()))
	log.Println(conn.GetState())

	return

	loginClient := pb.NewLoginClient(conn)
	resps,err := loginClient.UserLogin(context.Background(),&pb.LoginRequest{Name: "小王",Mobile: "12312312312"})

	log.Println(err)
	log.Println(fmt.Sprintf("%s",resp.GetMessage()))
	log.Println(fmt.Sprintf("%s",resp.Message))
	log.Println(fmt.Sprintf("%s",resp.String()))

	log.Println(resps.GetMessage())
	log.Println(resps.UserId)
	log.Println(resps.String())
}