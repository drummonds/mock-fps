package jsonapi

import (
	"encoding/json"
	"net/http"
	"strconv"
)

func WriteError(w http.ResponseWriter, status int, title, detail string) {
	w.Header().Set("Content-Type", "application/vnd.api+json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(ErrorResponse{
		Errors: []Error{
			{
				Status: strconv.Itoa(status),
				Title:  title,
				Detail: detail,
			},
		},
	})
}

func NotFound(w http.ResponseWriter, resourceType, id string) {
	WriteError(w, http.StatusNotFound, "Not Found", resourceType+" "+id+" not found")
}

func BadRequest(w http.ResponseWriter, detail string) {
	WriteError(w, http.StatusBadRequest, "Bad Request", detail)
}

func Conflict(w http.ResponseWriter, detail string) {
	WriteError(w, http.StatusConflict, "Conflict", detail)
}

func InternalError(w http.ResponseWriter) {
	WriteError(w, http.StatusInternalServerError, "Internal Server Error", "")
}
