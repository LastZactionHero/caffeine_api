package main

// APIError field error message
import (
	"encoding/json"
	"net/http"
)

// APIError message for building API error
type APIError struct {
	Field   string
	Message string
}

func optionsHandler(w http.ResponseWriter, r *http.Request) {
	apiApplyCorsHeaders(w, r)
}

func apiErrorHandler(w http.ResponseWriter, r *http.Request, status int, errors []APIError) {
	w.WriteHeader(status)

	errorBody := map[string]interface{}{}
	for _, error := range errors {
		errorBody[error.Field] = error.Message
	}

	b, _ := json.Marshal(errorBody)
	w.Write(b)
}

func apiApplyCorsHeaders(w http.ResponseWriter, r *http.Request) {
	origin := r.Header.Get("Origin")
	w.Header().Set("Access-Control-Allow-Origin", origin)
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Origin, X-Auth-Token")
	w.Header().Set("Access-Control-Allow-Methods", "GET, PUT, POST, DELETE, OPTIONS")
}
