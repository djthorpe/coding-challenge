package server

import (
	"context"
	"net/http"
	"regexp"
)

///////////////////////////////////////////////////////////////////////////////
// TYPES

type route struct {
	re      *regexp.Regexp
	handler http.HandlerFunc
}

// Router represents the requests that can be made to the server and routes them
// through to the right handler function
type Router struct {
	routes []route
}

///////////////////////////////////////////////////////////////////////////////
// CONSTANTS

const (
	KeyParams = "KeyParams"
)

///////////////////////////////////////////////////////////////////////////////
// LIFECYCLE

func NewRouter() *Router {
	return new(Router)
}

///////////////////////////////////////////////////////////////////////////////
// HANDLERS

// AddHandlerFunc adds a new backend handler for the server, with pattern matching. The
// groups within the pattern can be accessed using RequestParams() method
func (this *Router) AddHandlerFunc(re *regexp.Regexp, handler http.HandlerFunc) {
	this.routes = append(this.routes, route{re, handler})
}

// AddHandler adds a new backend handler for the server, with pattern matching. The
// groups within the pattern can be accessed using RequestParams() method
func (this *Router) AddHandler(re *regexp.Regexp, handler http.Handler) {
	this.routes = append(this.routes, route{re, handler.ServeHTTP})
}

// ServeHTTP provides the server routing to handlers based on regular expression
// patterns
func (this *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	for _, route := range this.routes {
		if args := route.re.FindStringSubmatch(req.URL.Path); len(args) >= 1 {
			route.handler(w, req.Clone(context.WithValue(req.Context(), KeyParams, args[1:])))
			return
		}
	}
	// No matched route
	ServeError(w, http.StatusNotFound)
}
