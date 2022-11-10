package db

import "context"

type DemoDB struct {
	//ID primitive.ObjectID `bson:"_id"`
	Ctx context.Context `bson:"-"`
}

func (this *DemoDB)DB() string {
	return "demo"
}
