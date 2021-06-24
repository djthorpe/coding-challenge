package backend_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	// Modules
	backend "github.com/djthorpe/coding-challenge/pkg/backend"
	schema "github.com/djthorpe/coding-challenge/pkg/schema"
	server "github.com/djthorpe/coding-challenge/pkg/server"
)

const (
	TestData = "../../data/reports.json"
)

func TestBackend_001(t *testing.T) {
	f, err := os.Open(TestData)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	backend, err := backend.NewBackend(f)
	if err != nil {
		t.Fatal(err)
	}
	for _, report := range backend.Reports {
		if report := backend.FindReport(report.Id); report == nil {
			t.Fatalf("Unexpected return from FindReport for report=%q", report.Id)
		}
	}
}

func TestBackend_002(t *testing.T) {
	// Create backend
	f, err := os.Open(TestData)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	backend, err := backend.NewBackend(f)
	if err != nil {
		t.Fatal(err)
	}

	// Serve all reports
	w := httptest.NewRecorder()
	if req, err := http.NewRequest(http.MethodGet, "/reports", nil); err != nil {
		t.Fatal(err)
	} else {
		backend.ServeReports(w, req)
	}

	// Check response
	if w.Result().StatusCode != http.StatusOK {
		t.Fatal("Unexpected status", w.Result().Status)
	}
	if w.Result().Header.Get("Content-Type") != "application/json" {
		t.Fatal("Unexpected content type", w.Result().Header.Get("Content-Type"))
	}
}

func TestBackend_003(t *testing.T) {
	// Create backend
	f, err := os.Open(TestData)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	backend, err := backend.NewBackend(f)
	if err != nil {
		t.Fatal(err)
	}

	for _, report := range backend.Reports {
		// Serve single report
		w := httptest.NewRecorder()
		if req, err := http.NewRequest(http.MethodGet, "/reports/"+report.Id, nil); err != nil {
			t.Fatal(err)
		} else {
			backend.ServeReports(w, req)
		}

		// Check response
		if w.Result().StatusCode != http.StatusOK {
			t.Fatal("Unexpected status", w.Result().Status)
		}
		if w.Result().Header.Get("Content-Type") != "application/json" {
			t.Fatal("Unexpected content type", w.Result().Header.Get("Content-Type"))
		}
	}

}

func TestBackend_004(t *testing.T) {
	// Create backend
	f, err := os.Open(TestData)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	backend, err := backend.NewBackend(f)
	if err != nil {
		t.Fatal(err)
	}

	for _, report := range backend.Reports {
		payload, err := json.Marshal(schema.Ticket{
			State: "TEST",
		})
		if err != nil {
			t.Fatal(err)
		}
		// Update single report
		w := httptest.NewRecorder()
		if req, err := http.NewRequest(http.MethodPut, "/reports/"+report.Id, bytes.NewReader(payload)); err != nil {
			t.Fatal(err)
		} else {
			backend.ServeReport(w, server.RequestWithParams(req, []string{report.Id}))
		}

		// Check response
		if w.Result().StatusCode != http.StatusOK {
			t.Fatal("Unexpected status", w.Result().Status)
		}
		if w.Result().Header.Get("Content-Type") != "application/json" {
			t.Fatal("Unexpected content type", w.Result().Header.Get("Content-Type"))
		}
	}

}
