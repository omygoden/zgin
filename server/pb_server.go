package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/balancer/roundrobin"
	"net"
)
import pb "zgin/pb"

type Greeter struct {
}

func (this *Greeter) SayHello(ctx context.Context, r *pb.HelloRequest) (*pb.HelloReply, error) {
	return &pb.HelloReply{Message: r.Name}, nil
}

type Login struct {

}

func (this *Login)UserLogin(ctx context.Context,r *pb.LoginRequest)  (*pb.LoginResponse,error){
	return &pb.LoginResponse{Code: "200",Message: "",UserId: "1"},nil
}

func main()  {
	server := grpc.NewServer()

	pb.RegisterGreeterServer(server,&Greeter{})
	pb.RegisterLoginServer(server,&Login{})
	grpc.WithDefaultServiceConfig(roundrobin.Name)
	lis,_ := net.Listen("tcp",":11000")
	fmt.Println("listen start ...")
	server.Serve(lis)
}
