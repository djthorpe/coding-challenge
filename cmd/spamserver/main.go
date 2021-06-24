package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"regexp"

	// Modules
	spamserver "github.com/djthorpe/coding-challenge"
	backend "github.com/djthorpe/coding-challenge/pkg/backend"
	server "github.com/djthorpe/coding-challenge/pkg/server"
)

///////////////////////////////////////////////////////////////////////////////
// GLOBALS

var (
	flagAddr = flag.String("addr", ":8080", "Address for server to listen on")
)

var (
	reRoot       = regexp.MustCompile(`^/.*$`)               // Matches any static file
	reAllReports = regexp.MustCompile(`^/reports/?$`)        // All reports
	reOneReport  = regexp.MustCompile(`^/reports/([\w-]+)$`) // Single report keyed by ID
)

///////////////////////////////////////////////////////////////////////////////
// MAIN

func main() {
	flag.Parse()

	// Create server and context for cancelling server
	server := server.NewServerWithConfig(*flagAddr)
	ctx := HandleSignal()

	// Open data for reading into the backend
	if flag.NArg() != 1 {
		fmt.Fprintln(os.Stderr, "Requires filename for backend data as argument")
		os.Exit(-1)
	}
	r, err := os.Open(flag.Arg(0))
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(-1)
	}
	defer r.Close()

	// Create backend, add handlers
	if backend, err := backend.NewBackend(r); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(-1)
	} else {
		server.AddHandlerFunc(reAllReports, backend.ServeReports)
		server.AddHandlerFunc(reOneReport, backend.ServeReport)
		server.AddHandler(reRoot, http.FileServer(http.FS(spamserver.Content)))
	}

	// Run the server until cancel
	log.Printf("Running server: %q", *flagAddr)
	if err := server.Run(ctx); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(-1)
	}
}

///////////////////////////////////////////////////////////////////////////////
// METHODS

// Handle signals - call cancel on returned context when interrupt received
func HandleSignal() context.Context {
	ctx, cancel := context.WithCancel(context.Background())
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)
	go func() {
		s := <-ch
		log.Println("Got signal:", s)
		cancel()
	}()
	return ctx
}
