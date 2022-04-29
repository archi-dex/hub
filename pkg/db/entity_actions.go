package db

import (
	"context"
	"regexp"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func GetEntity(ctx context.Context, id string) (*Entity, error) {
	if id == "" {
		return nil, ErrorInvalidEntity.Tracef("must specify id")
	}

	var entity Entity
	if err := getCol().FindOne(ctx, bson.M{"_id": id}).Decode(&entity); err != nil {
		return nil, err
	}

	return &entity, nil
}

func CreateEntity(ctx context.Context, path string, attributes map[string]interface{}) (*Entity, error) {
	if path == "" || attributes == nil {
		return nil, ErrorInvalidEntity.Tracef("must specify attributes")
	}

	r := regexp.MustCompile(`\/`)
	r.Split(path[1:len(path)-1], -1)

	entity := NewEntity(path, attributes)
	var result *mongo.InsertOneResult
	var err error
	if result, err = getCol().InsertOne(ctx, entity); err != nil {
		return nil, err
	}

	return GetEntity(ctx, (result.InsertedID).(string))
}

func UpdateEntity(ctx context.Context, id string, attributes map[string]interface{}) (*Entity, error) {
	if id == "" || attributes == nil {
		return nil, ErrorInvalidEntity.Tracef("must specify both id and attributes")
	}

	filter := bson.M{"_id": id}
	update := bson.M{"$set": bson.M{
		"attributes": attributes,
		"updated_at": time.Now().UTC(),
	}}

	returnType := options.After
	opts := &options.FindOneAndUpdateOptions{ReturnDocument: &returnType}
	if err := getCol().FindOneAndUpdate(ctx, filter, update, opts).Err(); err != nil {
		return nil, err
	}

	return GetEntity(ctx, id)
}

func DeleteEntity(ctx context.Context, id string) error {
	if id == "" {
		return ErrorInvalidEntity.Tracef("must specify id")
	}

	filter := bson.M{"_id": id}
	return getCol().FindOneAndDelete(ctx, filter).Err()
}

func ListEntities(ctx context.Context, filter interface{}) ([]*Entity, error) {
	cursor, err := getCol().Find(ctx, filter)
	if err != nil {
		return nil, err
	}

	var results []*Entity
	if err := cursor.All(ctx, &results); err != nil {
		return nil, err
	}

	return results, nil
}
