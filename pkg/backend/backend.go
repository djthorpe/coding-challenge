package backend

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"sync"

	// Modules
	schema "github.com/djthorpe/coding-challenge/pkg/schema"
	server "github.com/djthorpe/coding-challenge/pkg/server"
)

///////////////////////////////////////////////////////////////////////////////
// TYPES

// Backend represents the database and the methods needed to manipulate the
// data in the database.
type Backend struct {
	sync.Mutex

	Reports []*schema.Report `json:"elements"`
}

///////////////////////////////////////////////////////////////////////////////
// NEW

// Create a new backend with data from a file
func NewBackend(r io.Reader) (*Backend, error) {
	this := new(Backend)

	// Read in JSON data
	if err := json.NewDecoder(r).Decode(this); err != nil {
		return nil, err
	}

	// Return success
	return this, nil
}

///////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

// FindReport finds a report within the database. Returns nil if no report with
// that ID is found
func (this *Backend) FindReport(value string) *schema.Report {
	// Here's a linear (inefficient) search. Normally you wouldn't have an
	// in-memory database nor a linear search, but at least use a hashmap...
	for _, report := range this.Reports {
		if report.Id == value {
			return report
		}
	}
	// Return nil for not found
	return nil
}

///////////////////////////////////////////////////////////////////////////////
// HANDLERS

// ServeReports serves all reports in the backend
func (this *Backend) ServeReports(w http.ResponseWriter, req *http.Request) {
	// Ensure it's a GET method
	if req.Method != http.MethodGet {
		server.ServeError(w, http.StatusBadRequest)
		return
	}

	// Only serve reports which aren't RESOLVED
	var results []*schema.Report
	for _, report := range this.Reports {
		if report.State != "RESOLVED" {
			results = append(results, report)
		}
	}

	// Serve JSON response of all reports
	server.ServeJSON(w, results, http.StatusOK)
}

// ServeReport with GET will find and serve a single report from the backend,
// or with PUT will update an existing report
func (this *Backend) ServeReport(w http.ResponseWriter, req *http.Request) {
	// Ensure it's a GET or PUT method
	if req.Method != http.MethodGet && req.Method != http.MethodPut {
		server.ServeError(w, http.StatusBadRequest)
		return
	}

	// Get the report id
	args := server.RequestParams(req)
	if len(args) != 1 {
		server.ServeError(w, http.StatusBadRequest)
		return
	}
	report := this.FindReport(args[0])
	if report == nil {
		server.ServeError(w, http.StatusNotFound)
		return
	}

	// We lock requests so that the ticket can only be updated
	// in series
	this.Mutex.Lock()
	defer this.Mutex.Unlock()

	// If it's a PUT method then read in the body
	if req.Method == http.MethodPut {
		var ticket schema.Ticket

		// Read in the ticket
		defer req.Body.Close()
		if data, err := ioutil.ReadAll(req.Body); err != nil {
			server.ServeError(w, http.StatusBadRequest)
			return
		} else if err := json.Unmarshal(data, &ticket); err != nil {
			server.ServeError(w, http.StatusBadRequest)
			return
		}

		// Apply the new state
		if ticket.State != "" {
			report.State = ticket.State
		}
	}

	// Serve JSON response of the found reports
	server.ServeJSON(w, report, http.StatusOK)
}
