package server

import (
	"context"
	"errors"
	"net/http"
	"time"

	// Modules
	multierror "github.com/hashicorp/go-multierror"
)

///////////////////////////////////////////////////////////////////////////////
// TYPES

// Server is the in-built golang server and adds a router on top
type Server struct {
	*http.Server
	*Router
}

///////////////////////////////////////////////////////////////////////////////
// GLOBALS

const (
	DefaultTimeout = 10 * time.Second
)

///////////////////////////////////////////////////////////////////////////////
// LIFECYCLE

func NewServerWithConfig(addr string) *Server {
	this := new(Server)
	this.Server = &http.Server{}
	this.Router = NewRouter()

	// Set server parameters
	this.Addr = addr
	this.Handler = this.Router
	this.ReadHeaderTimeout = DefaultTimeout
	this.IdleTimeout = DefaultTimeout

	// Return success
	return this
}

///////////////////////////////////////////////////////////////////////////////
// METHODS

// Run starts the server in the foreground and returns when context is cancelled
func (this *Server) Run(ctx context.Context) error {
	var result error
	go func() {
		<-ctx.Done()
		if err := this.Close(); err != nil {
			result = multierror.Append(result, err)
		}
	}()
	if err := this.ListenAndServe(); err != nil && errors.Is(err, http.ErrServerClosed) == false {
		result = multierror.Append(result, err)
	}
	return result
}
