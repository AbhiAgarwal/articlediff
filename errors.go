package main

import (
	"encoding/json"
	"net/http"
)

type Errors struct {
	Errors []*Error `json:"errors"`
}

type Error struct {
	Id     string `json:"id"`
	Status int    `json:"status"`
	Title  string `json:"title"`
	Detail string `json:"detail"`
}

func WriteError(w http.ResponseWriter, err *Error) {
	w.Header().Set("Content-Type", "application/vnd.api+json")
	w.WriteHeader(err.Status)
	json.NewEncoder(w).Encode(Errors{[]*Error{err}})
}

var (
	ErrBadRequest           = &Error{"bad_request", http.StatusBadRequest, "Bad request", "Request body is not well-formed. It must be JSON."}
	ErrNotFound             = &Error{"not_found", http.StatusNotFound, "Not Found", "Resource or item you requested for has not been found."}
	ErrNotAcceptable        = &Error{"not_acceptable", http.StatusNotAcceptable, "Not Acceptable", "Accept header must be set to 'application/vnd.api+json'."}
	ErrUnsupportedMediaType = &Error{"unsupported_media_type", http.StatusUnsupportedMediaType, "Unsupported Media Type", "Content-Type header must be set to: 'application/vnd.api+json'."}
	ErrInternalServer       = &Error{"internal_server_error", http.StatusInternalServerError, "Internal Server Error", "Something went wrong."}
)
