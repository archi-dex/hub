package api

import (
	"context"
	"io/ioutil"
	"net/http"

	"github.com/archi-dex/ingester/pkg/db"
	"github.com/archi-dex/ingester/pkg/util"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func createEntity(ctx context.Context, logger util.Logger) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		buf, err := ioutil.ReadAll(r.Body)
		if err != nil {
			logger.Errorw(msgErrorReadingBody, "err", err)
			respondBadRequest(rw, msgErrorReadingBody, nil)
			return
		}

		partialEntity, err := db.EntityAttributesFromBytes(buf)
		if err != nil {
			logger.Errorw(msgErrorParsingBody, "err", err)
			respondBadRequest(rw, msgErrorParsingBody, nil)
			return
		}

		entity, err := db.CreateEntity(ctx, partialEntity.Path, partialEntity.Attributes)
		if err != nil {
			logger.Errorw("error creating entity", "err", err)
			respondError(rw, "error creating entity", err.Error())
			return
		}

		logger.Infow("created entity", "entity", entity)
		respond(rw, http.StatusCreated, "", entity)
	}
}

func readEntity(ctx context.Context, logger util.Logger) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		id, ok := mux.Vars(r)["id"]
		if !ok {
			logger.Errorw(msgMustSpecifyId)
			respondBadRequest(rw, msgMustSpecifyId, nil)
			return
		}

		entity, err := db.GetEntity(ctx, id)
		if err == mongo.ErrNoDocuments {
			logger.Errorw("entity not found", "err", err, "id", id)
			respondNotFound(rw)
			return
		}

		if err != nil {
			logger.Errorw("error reading entity", "err", err)
			respondError(rw, "error reading entity", err.Error())
			return
		}

		logger.Infow("read entity", "entity", entity)
		respondOk(rw, entity)
	}
}

func updateEntity(ctx context.Context, logger util.Logger) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		id, ok := mux.Vars(r)["id"]
		if !ok {
			logger.Errorw(msgMustSpecifyId)
			respondBadRequest(rw, msgMustSpecifyId, nil)
			return
		}

		buf, err := ioutil.ReadAll(r.Body)
		if err != nil {
			logger.Errorw(msgErrorReadingBody, "err", err)
			respondBadRequest(rw, msgErrorReadingBody, nil)
			return
		}

		partialEntity, err := db.EntityAttributesFromBytes(buf)
		if err != nil {
			logger.Errorw(msgErrorParsingBody, "err", err, "body", buf)
			respondBadRequest(rw, msgErrorParsingBody, nil)
			return
		}

		entity, err := db.UpdateEntity(ctx, id, partialEntity.Attributes)
		if err == mongo.ErrNoDocuments {
			logger.Errorw("entity not found", "err", err, "id", id)
			respondNotFound(rw)
			return
		}

		if err != nil {
			logger.Errorw("error updating entity", "err", err, "id", id)
			respondError(rw, "error updating entity", err.Error())
		}

		logger.Infow("updated entity", "entity", entity)
		respondOk(rw, entity)
	}
}

func deleteEntity(ctx context.Context, logger util.Logger) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		id, ok := mux.Vars(r)["id"]
		if !ok {
			logger.Errorw(msgMustSpecifyId)
			respondBadRequest(rw, msgMustSpecifyId, nil)
			return
		}

		err := db.DeleteEntity(ctx, id)
		if err == mongo.ErrNoDocuments {
			logger.Errorw("entity not found", "err", err, "id", id)
			respondNotFound(rw)
			return
		}

		if err != nil {
			logger.Errorw("error deleting entity", "err", err)
			respondError(rw, "error deleting entity", err.Error())
			return
		}

		logger.Infow("deleted entity", "id", id)
		respondOk(rw, nil)
	}
}

func listEntities(ctx context.Context, logger util.Logger) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		var entities []*db.Entity
		var err error
		var detail string

		// TODO support filter options
		entities, err = db.ListEntities(ctx, bson.D{})
		detail = "error listing entities"

		if err == mongo.ErrNoDocuments || entities == nil {
			logger.Errorw(detail + " - none found")
			respondNotFound(rw)
			return
		}

		if err != nil {
			logger.Errorw(detail, "err", err)
			respondError(rw, detail, err.Error())
			return
		}

		logger.Infow("listed entities", "count", len(entities))
		respondOk(rw, entities)
	}
}
