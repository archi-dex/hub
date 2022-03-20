package db

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Index struct {
	ID   primitive.ObjectID `bson:"_id"  json:"_id,omitempty"`
	Name string             `bson:"name" json:"name"`
}

func IndexFromBson(input bson.D) {}
