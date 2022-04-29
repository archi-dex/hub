package db

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/archi-dex/ingester/pkg/util"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type SortType = int32

const (
	SortAsc SortType = 1
	SortDes SortType = -1
)

type Facet struct {
	Id   string   `json:"id"`
	Key  string   `json:"key"`
	Sort SortType `json:"sort"`
}

func FacetFromBytes(raw []byte) (*Facet, error) {
	var facet Facet
	if err := json.Unmarshal(raw, &facet); err != nil {
		return nil, err
	}

	return &facet, nil
}

func FacetFromBson(raw interface{}) (*Facet, error) {
	var list []interface{}
	var err error
	if list, err = util.CoerceToSlice(raw); err != nil || list == nil {
		return nil, fmt.Errorf("could not coerce into slice: %s", err)
	}

	if len(list) != 3 {
		return nil, fmt.Errorf("index bson must have three elements, got %d", len(list))
	}

	var id string
	var key string
	var sort SortType

	if element, ok := (list[2]).(bson.E); ok && element.Key == "name" {
		if value, ok := (element.Value).(string); ok {
			id = value
		}
	}

	if element, ok := (list[1]).(bson.E); ok && element.Key == "key" {
		if element, err := util.CoerceToSlice(element.Value); err == nil && len(element) == 1 {
			if element, ok := (element[0]).(bson.E); ok {
				key = element.Key
				if value, ok := (element.Value).(SortType); ok {
					sort = value
				}
			}
		}
	}

	if id == "" || key == "" || sort == 0 {
		return nil, fmt.Errorf("facet schema missing attributes, (id: '%s', key: '%s', sort: '%d')", id, key, sort)
	}

	return &Facet{id, key, sort}, nil
}

func AssertFacet(ctx context.Context, key string, sort SortType) (*Facet, error) {
	var id string
	var err error
	index := mongo.IndexModel{Keys: bson.M{key: sort}}
	if id, err = getCol().Indexes().CreateOne(ctx, index); err != nil {
		return nil, err
	}

	return &Facet{id, key, sort}, nil
}

func foo(key string, sort SortType) {
	getCol().Indexes().CreateOne(context.TODO(), mongo.IndexModel{Keys: bson.M{key: sort}})
}
