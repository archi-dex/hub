package db

import (
	"context"
	"fmt"

	"github.com/archi-dex/ingester/pkg/util"
	"go.mongodb.org/mongo-driver/mongo"
)

func CreateFacet(ctx context.Context, logger util.Logger, key string, sort SortType) (*Facet, error) {
	return AssertFacet(ctx, key, sort)
}

func DeleteFacet(ctx context.Context, logger util.Logger, id string) error {
	if id == "_id_" {
		return fmt.Errorf("cannot delete identity index")
	}

	if _, err := getCol().Indexes().DropOne(ctx, id); err != nil {
		return err
	}

	return nil
}

func ListFacets(ctx context.Context, logger util.Logger) ([]*Facet, error) {
	var cursor *mongo.Cursor
	var err error

	if cursor, err = getCol().Indexes().List(ctx); err != nil {
		return nil, err
	}

	var results []interface{}
	if err = cursor.All(ctx, &results); err != nil {
		return nil, err
	}

	parsed := make([]*Facet, len(results))
	for i, result := range results {
		if parsed[i], err = FacetFromBson(result); err != nil {
			return nil, err
		}
	}

	return parsed, nil
}
