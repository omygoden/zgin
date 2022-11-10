package mongo

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"zgin/global"
	"zgin/internal/model/mongo/db"
)

type User struct {
	db.DemoDB  `bson:"-"`
	ID         primitive.ObjectID `bson:"_id" json:"_id,omitempty"`
	Name       string             `bson:"name"`
	Mobile     string             `bson:"mobile"`
	Age        int                `bson:"age"`
	From       int                `bson:"from"`
	CreateTime int64              `bson:"create_time"`
}

func (this *User) Collection() *mongo.Collection {
	return global.MongoClient.Database(this.DB()).Collection("user")
}
