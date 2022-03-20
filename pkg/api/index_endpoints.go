package api

import (
	"context"
	"net/http"

	"github.com/archi-dex/ingester/pkg/db"
	"github.com/archi-dex/ingester/pkg/util"
)

func createIndex(ctx context.Context, logger util.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) { db.CreateIndex(ctx, logger) }
}

func deleteIndex(ctx context.Context, logger util.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) { db.DeleteIndex(ctx, logger) }
}

func listIndicies(ctx context.Context, logger util.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) { db.ListIndicies(ctx, logger) }
}
