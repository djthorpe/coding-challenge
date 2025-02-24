package server

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
)

///////////////////////////////////////////////////////////////////////////////
// CONSTANTS

const (
	ContentTypeJSON = "application/json"
	KeyParams       = "KeyParams"
)

///////////////////////////////////////////////////////////////////////////////
// METHODS

// RequestWithParams returns a request with set parameter
func RequestWithParams(req *http.Request, params []string) *http.Request {
	return req.Clone(context.WithValue(req.Context(), KeyParams, params))
}

// RequestParams returns parameters for a path or nil otherwise
func RequestParams(req *http.Request) []string {
	if args, ok := req.Context().Value(KeyParams).([]string); ok {
		return args
	} else {
		return nil
	}
}

// ServeError is a utility function to serve an error code as plaintext on HTTP
func ServeError(w http.ResponseWriter, code int) {
	http.Error(w, http.StatusText(code), code)
}

// ServeJSON is a utility function to serve an arbitary object as JSON
func ServeJSON(w http.ResponseWriter, v interface{}, code int) {
	w.Header().Add("Content-Type", ContentTypeJSON)
	w.WriteHeader(code)
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	if err := enc.Encode(v); err != nil {
		log.Println(err)
	}
}
