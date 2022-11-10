package mongo

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"os"
	"zgin/global"
	"zgin/internal/model/mongo/db"
)

type Spider struct {
	db.DemoDB
	ID         primitive.ObjectID `bson:"_id"`
	Sku        string             `json:"sku" bson:"sku"`
	Content    string             `json:"content" bson:"content"`
	From       int8               `json:"from"`
	CreateTime int64              `json:"create_time" bson:"create_time"`
}

func (this *Spider) Collection() *mongo.Collection {
	return global.MongoClient.Database(this.DB()).Collection("spider")
}

func (this *Spider) InsertOne() *mongo.InsertOneResult {
	id, err := this.Collection().InsertOne(context.TODO(), &Spider{})
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	return id
}

func (this *Spider) InsertMany(list []Spider) *mongo.InsertManyResult {
	var data = make([]interface{}, len(list), len(list))
	for k, v := range list {
		data[k] = v
	}

	ids, err := this.Collection().InsertMany(context.TODO(), data)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	return ids
}

func (this *Spider) FindOne(where bson.D, fields interface{}) *mongo.SingleResult {
	return this.Collection().FindOne(context.TODO(), where, &options.FindOneOptions{
		Projection: fields,
	})
}

func (this *Spider) FindMany(where bson.D, fields interface{}, limit int64) []Spider {
	var findOptions options.FindOptions
	var list = make([]Spider, 0, limit)
	findOptions.Limit = &limit
	if fields != nil {
		findOptions.Projection = fields
	}

	result, err := this.Collection().Find(context.TODO(), where, &findOptions)

	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	for result.Next(context.TODO()) {
		var value Spider
		result.Decode(&value)
		list = append(list, value)
	}

	return list
}

func (this *Spider) UpdateOne(where bson.M, data Spider) *mongo.UpdateResult {
	updateData := structToBsonM(data)
	result, err := this.Collection().UpdateOne(context.TODO(), where, updateData)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	return result
}
