package api

import (
	"context"
	"net/http"

	"github.com/archi-dex/ingester/pkg/util"
	"github.com/gorilla/mux"
)

func NewRouter(ctx context.Context, logger util.Logger) *mux.Router {
	r := mux.NewRouter()

	// middleware

	r.Use(loggingMiddlware(logger))
	r.Use(defaultHeadersMiddleware)

	// endpoints

	r.HandleFunc("/api/ingest", ingest(ctx, logger)).Schemes("ws")

	r.HandleFunc("/api/entity", createEntity(ctx, logger)).Methods(http.MethodPost).Schemes("http")
	r.HandleFunc("/api/entity/{id}", readEntity(ctx, logger)).Methods(http.MethodGet).Schemes("http")
	r.HandleFunc("/api/entity/{id}", updateEntity(ctx, logger)).Methods(http.MethodPost).Schemes("http")
	r.HandleFunc("/api/entity/{id}", deleteEntity(ctx, logger)).Methods(http.MethodDelete).Schemes("http")
	r.HandleFunc("/api/entities", listEntities(ctx, logger)).Methods(http.MethodGet).Schemes("http")
	r.HandleFunc("/api/entities/{id}", listEntities(ctx, logger)).Methods(http.MethodGet).Schemes("http")

	r.HandleFunc("/api/index", createIndex(ctx, logger)).Methods(http.MethodPost).Schemes("http")
	r.HandleFunc("/api/index/{id}", deleteIndex(ctx, logger)).Methods(http.MethodDelete).Schemes("http")
	r.HandleFunc("/api/index", listIndicies(ctx, logger)).Methods(http.MethodGet).Schemes("http")

	return r
}
