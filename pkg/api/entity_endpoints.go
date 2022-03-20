package api

import (
	"context"
	"io/ioutil"
	"net/http"

	"github.com/archi-dex/ingester/pkg/db"
	"github.com/archi-dex/ingester/pkg/util"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	msgErrorCreatingEntity     responseDetail = "error creating entity"
	msgCreatedEntity           responseDetail = "created entity"
	msgErrorReadingEntity      responseDetail = "error reading entity"
	msgReadEntity              responseDetail = "read entity"
	msgErrorUpdatingEntity     responseDetail = "error updating entity"
	msgUpdatedEntity           responseDetail = "updated entity"
	msgErrorDeletingEntity     responseDetail = "error deleting entity"
	msgDeletedEntity           responseDetail = "deleted entity"
	msgErrorListingEntities    responseDetail = "error listing entities"
	msgListedEntities          responseDetail = "listed entities"
	msgErrorEntityNotFound     responseDetail = "entity not found"
	msgErrorAttributesRequired responseDetail = "must specify attributes"
)

func createEntity(ctx context.Context, logger util.Logger) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		buf, err := ioutil.ReadAll(r.Body)
		if err != nil {
			logger.Errorw(msgErrorReadingBody, "err", err)
			respondBadRequest(rw, msgErrorReadingBody, nil)
			return
		}

		entity, err := db.EntityFromBytes(buf)
		if err != nil {
			logger.Errorw(msgErrorParsingBody, "err", err)
			respondBadRequest(rw, msgErrorParsingBody, nil)
			return
		}

		if err := db.CreateEntity(ctx, *entity); err != nil {
			logger.Errorw(msgErrorCreatingEntity, "err", err)
			respondError(rw, msgErrorCreatingEntity, nil)
			return
		}

		logger.Infow(msgCreatedEntity, "entity", entity)
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

		entity, err := db.ReadEntity(ctx, id)
		if err == mongo.ErrNoDocuments {
			logger.Errorw(msgErrorEntityNotFound, "err", err, "id", id)
			respondNotFound(rw)
			return
		}

		if err != nil {
			logger.Errorw(msgErrorReadingEntity, "err", err)
			respondError(rw, msgErrorReadingEntity, nil)
			return
		}

		logger.Infow(msgReadEntity, "entity", entity)
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

		partialEntity, err := db.EntityFromBytes(buf)
		if err != nil {
			logger.Errorw(msgErrorParsingBody, "err", err, "body", buf)
			respondBadRequest(rw, msgErrorParsingBody, nil)
			return
		}

		if partialEntity.Attributes == nil {
			logger.Errorw(msgErrorParsingBody, "err", err, "body", buf)
			respondBadRequest(rw, msgErrorAttributesRequired, nil)
			return
		}

		entity, err := db.UpdateEntity(ctx, id, partialEntity.Attributes)
		if err == mongo.ErrNoDocuments {
			logger.Errorw(msgErrorEntityNotFound, "err", err, "id", id)
			respondNotFound(rw)
			return
		}

		if err != nil {
			logger.Errorw(msgErrorUpdatingEntity, "err", err, "id", id)
			respondError(rw, msgErrorUpdatingEntity, nil)
		}

		logger.Infow(msgUpdatedEntity, "entity", entity)
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
			logger.Errorw(msgErrorEntityNotFound, "err", err, "id", id)
			respondNotFound(rw)
			return
		}

		if err != nil {
			logger.Errorw(msgErrorDeletingEntity, "err", err)
			respondError(rw, msgErrorDeletingEntity, nil)
			return
		}

		logger.Infow(msgDeletedEntity, "id", id)
		respondOk(rw, nil)
	}
}

func listEntities(ctx context.Context, logger util.Logger) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		id, ok := mux.Vars(r)["id"]
		var entities []*db.Entity
		var err error
		var detail string

		switch {
		case id == "root" && ok:
			entities, err = db.ListEntityRoot(ctx, r.URL.Query())
			detail = "error listing entity root"
		case ok:
			entities, err = db.ListEntityChildren(ctx, id, r.URL.Query())
			detail = "error listing entity children"
		default:
			entities, err = db.ListEntities(ctx, r.URL.Query())
			detail = "error listing entities"
		}

		if err == mongo.ErrNoDocuments || entities == nil {
			logger.Errorw(detail, "err", err)
			respondNotFound(rw)
			return
		}

		if err != nil {
			logger.Errorw(detail, "err", err)
			respondError(rw, detail, nil)
			return
		}

		logger.Infow(msgListedEntities, "count", len(entities))
		respondOk(rw, entities)
	}
}
