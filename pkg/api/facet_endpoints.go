package api

import (
	"context"
	"io/ioutil"
	"net/http"

	"github.com/archi-dex/ingester/pkg/db"
	"github.com/archi-dex/ingester/pkg/util"
	"github.com/gorilla/mux"
)

func createFacet(ctx context.Context, logger util.Logger) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		var buf []byte
		var err error
		if buf, err = ioutil.ReadAll(r.Body); err != nil {
			logger.Errorw(msgErrorReadingBody, "err", err)
			respondBadRequest(rw, msgErrorReadingBody, nil)
			return
		}

		var facet *db.Facet
		if facet, err = db.FacetFromBytes(buf); err != nil {
			logger.Errorw(msgErrorParsingBody, "err", err)
			respondBadRequest(rw, msgErrorParsingBody, nil)
			return
		}

		if facet, err = db.CreateFacet(ctx, logger, facet.Key, facet.Sort); err != nil {
			logger.Errorw("failed to create facet", "err", err)
			respondError(rw, "failed to create facet", err.Error())
			return
		}

		logger.Infow("create facet", "facet", facet)
		respondOk(rw, facet)
	}
}

func deleteFacet(ctx context.Context, logger util.Logger) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		var id string
		var ok bool
		if id, ok = mux.Vars(r)["id"]; !ok || id == "" {
			logger.Errorw(msgMustSpecifyId)
			respondBadRequest(rw, msgMustSpecifyId, nil)
			return
		}

		if err := db.DeleteFacet(ctx, logger, id); err != nil {
			logger.Errorw("failed to delete facet", "err", err)
			respondError(rw, "failed to delete facet", err.Error())
			return
		}

		logger.Infow("delete facet", "id", id)
		respondOk(rw, nil)
	}
}

func listFacets(ctx context.Context, logger util.Logger) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		var facets []*db.Facet
		var err error
		if facets, err = db.ListFacets(ctx, logger); err != nil {
			logger.Errorw("failed to list facets", "err", err)
			respondError(rw, "failed to list facets", err.Error())
			return
		}

		logger.Infow("list facets", "facets", facets)
		respondOk(rw, facets)
	}
}
