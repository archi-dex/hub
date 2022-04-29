package api

import (
	"context"
	"net/http"

	"github.com/archi-dex/ingester/pkg/util"
	"github.com/gorilla/mux"
)

func NewRouter(ctx context.Context, logger util.Logger) *mux.Router {
	r := mux.NewRouter()
	r.StrictSlash(true)

	// middleware

	r.Use(loggingMiddlware(logger))
	r.Use(defaultHeadersMiddleware)

	// endpoints

	r.HandleFunc("/api/ingest", ingest(ctx, logger)).Schemes("http", "ws")

	r.HandleFunc("/api/entities", createEntity(ctx, logger)).Methods(http.MethodPost).Schemes("http")
	r.HandleFunc("/api/entities/{id}", readEntity(ctx, logger)).Methods(http.MethodGet).Schemes("http")
	r.HandleFunc("/api/entities/{id}", updateEntity(ctx, logger)).Methods(http.MethodPost).Schemes("http")
	r.HandleFunc("/api/entities/{id}", deleteEntity(ctx, logger)).Methods(http.MethodDelete).Schemes("http")
	r.HandleFunc("/api/entities", listEntities(ctx, logger)).Methods(http.MethodGet).Schemes("http")

	r.HandleFunc("/api/facets", createFacet(ctx, logger)).Methods(http.MethodPost).Schemes("http")
	r.HandleFunc("/api/facets/{id}", deleteFacet(ctx, logger)).Methods(http.MethodDelete).Schemes("http")
	r.HandleFunc("/api/facets", listFacets(ctx, logger)).Methods(http.MethodGet).Schemes("http")

	return r
}
