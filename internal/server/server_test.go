package server

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/cyberactivelabz/kickbucket/internal/engine"
)

type testProvider struct {
	name   string
	access engine.Access
}

func (provider testProvider) Name() string { return provider.name }
func (provider testProvider) Candidates(keyword string, words []string) []string {
	return []string{keyword + "-candidate"}
}
func (provider testProvider) Check(context.Context, string) (engine.Access, error) {
	return provider.access, nil
}

func TestScanEndpoint(t *testing.T) {
	handler := NewHandler(Config{Providers: []engine.Provider{testProvider{"s3", engine.AccessListable}}, Words: []string{"dev"}})
	request := httptest.NewRequest(http.MethodPost, "/api/scan", strings.NewReader(`{"keyword":"acme","providers":["s3"]}`))
	recorder := httptest.NewRecorder()
	handler.ServeHTTP(recorder, request)
	if recorder.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200: %s", recorder.Code, recorder.Body)
	}
	var results []engine.Result
	if err := json.NewDecoder(recorder.Body).Decode(&results); err != nil {
		t.Fatal(err)
	}
	if len(results) != 1 || results[0].Name != "acme-candidate" || results[0].Access != engine.AccessListable {
		t.Fatalf("results = %#v", results)
	}
}

func TestScanEndpointRejectsUnknownProvider(t *testing.T) {
	handler := NewHandler(Config{Providers: []engine.Provider{testProvider{"s3", engine.AccessPrivate}}})
	request := httptest.NewRequest(http.MethodPost, "/api/scan", strings.NewReader(`{"keyword":"acme","providers":["gcs"]}`))
	recorder := httptest.NewRecorder()
	handler.ServeHTTP(recorder, request)
	if recorder.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want 400", recorder.Code)
	}
}
