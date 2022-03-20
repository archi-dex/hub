package db

import (
	"context"
	"time"

	"github.com/archi-dex/ingester/pkg/util"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func getCol() *mongo.Collection {
	opts := util.GetOptions()
	return client.Database(opts.DbName).Collection(opts.DbCollection)
}

func CreateEntity(ctx context.Context, entity Entity) error {
	col := getCol()

	_, err := col.InsertOne(ctx, entity)
	if err != nil {
		return err
	}

	return nil
}

func ReadEntity(ctx context.Context, id string) (*Entity, error) {
	col := getCol()

	var entity Entity
	if err := col.FindOne(ctx, bson.D{{Key: "_id", Value: id}}).Decode(&entity); err != nil {
		return nil, err
	}

	return &entity, nil
}

func UpdateEntity(ctx context.Context, id string, attributes map[string]string) (*Entity, error) {
	col := getCol()

	filter := bson.D{{Key: "_id", Value: id}}
	update := bson.D{{Key: "$set", Value: bson.D{
		{Key: "attributes", Value: attributes},
		{Key: "updated_at", Value: time.Now().UTC()},
	}}}

	var result Entity
	returnType := options.After
	options := &options.FindOneAndUpdateOptions{ReturnDocument: &returnType}
	if err := col.FindOneAndUpdate(ctx, filter, update, options).Decode(&result); err != nil {
		return nil, err
	}

	return &result, nil
}

func DeleteEntity(ctx context.Context, id string) error {
	col := getCol()

	filter := bson.D{{Key: "_id", Value: id}}
	return col.FindOneAndDelete(ctx, filter).Err()
}

func ListEntities(ctx context.Context, filter util.Map) ([]*Entity, error) {
	return nil, nil
}

func ListEntityRoot(ctx context.Context, filter util.Map) ([]*Entity, error) {
	return nil, nil
}

func ListEntityChildren(ctx context.Context, id string, filter util.Map) ([]*Entity, error) {
	return nil, nil
}
