package api

import (
	"encoding/json"
	"net/http"
)

type responseDetail = string

const (
	msgMustSpecifyId    responseDetail = "must specify id"
	msgErrorReadingBody responseDetail = "error reading body"
	msgErrorParsingBody responseDetail = "error parsing body"
)

type Response struct {
	Message string      `json:"message"`
	Detail  *string     `json:"detail"`
	Data    interface{} `json:"data"`
}

func mustMarshal(data interface{}) []byte {
	result, _ := json.Marshal(data)
	return result
}

func respond(rw http.ResponseWriter, code int, detail string, data interface{}) {
	response := Response{
		Message: http.StatusText(code),
		Data:    data,
	}

	if detail != "" {
		response.Detail = &detail
	}

	rw.WriteHeader(code)
	rw.Write(mustMarshal(response))
}

func respondNotFound(rw http.ResponseWriter) {
	respond(rw, http.StatusNotFound, "", nil)
}

func respondError(rw http.ResponseWriter, detail string, data interface{}) {
	respond(rw, http.StatusInternalServerError, detail, data)
}

func respondBadRequest(rw http.ResponseWriter, detail string, data interface{}) {
	respond(rw, http.StatusBadRequest, detail, data)
}

func respondOk(rw http.ResponseWriter, data interface{}) {
	respond(rw, http.StatusOK, "", data)
}
