package mongodb

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"log"
	"time"
	"zgin/global"
)

func InitMongodb()  {
	ctx,cancel := context.WithTimeout(context.TODO(),time.Second * 3)
	defer cancel()
	optionClient := options.Client()
	optionClient.ApplyURI("mongodb://121.37.163.104:2717")
	optionClient.SetAuth(options.Credential{Username: "rootAdmin",Password: "Czy//299"})


	client,err := mongo.Connect(ctx,optionClient)
	if err != nil {
		panic("mongo连接失败："+err.Error())
	}
	ctx, cancel = context.WithTimeout(context.Background(), time.Second * 2)
	defer cancel()
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		panic("mongo连接成功，ping失败："+err.Error())
	}


	global.MongoClient = client

	log.Println("mongdo初始化成功")
}
