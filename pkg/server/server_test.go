package server_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"
	"time"

	// Modules
	spamserver "github.com/djthorpe/coding-challenge"
	server "github.com/djthorpe/coding-challenge/pkg/server"
)

func TestServer_001(t *testing.T) {
	if server := server.NewServerWithConfig(":8001"); server == nil {
		t.Fatal("Unexpected return from NewServerWithConfig")
	}
}

func TestServer_002(t *testing.T) {
	if server := server.NewServerWithConfig(":8001"); server == nil {
		t.Fatal("Unexpected return from NewServerWithConfig")
	} else {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		if err := server.Run(ctx); err != nil {
			t.Error(err)
		}
	}
}

func TestServer_003(t *testing.T) {
	server := server.NewServerWithConfig(":8001")
	if server == nil {
		t.Fatal("Unexpected return from NewServerWithConfig")
	}

	// Add static file handling
	server.AddHandler(regexp.MustCompile(`^/.*$`), http.FileServer(http.FS(spamserver.Content)))

	// Serve a static file
	w := httptest.NewRecorder()
	if req, err := http.NewRequest(http.MethodGet, "/html/", nil); err != nil {
		t.Fatal(err)
	} else {
		server.ServeHTTP(w, req)
	}

	// Check response
	if w.Result().StatusCode != http.StatusOK {
		t.Fatal("Unexpected status", w.Result().Status)
	}
}
