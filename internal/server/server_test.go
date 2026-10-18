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

type largeProvider struct{}

func (largeProvider) Name() string                         { return "s3" }
func (largeProvider) Candidates(string, []string) []string { return []string{"one", "two"} }
func (largeProvider) Check(context.Context, string) (engine.Access, error) {
	return engine.AccessPrivate, nil
}

func localRequest(method, path, body string) *http.Request {
	request := httptest.NewRequest(method, path, strings.NewReader(body))
	request.Host = "127.0.0.1:8080"
	return request
}

func TestScanEndpoint(t *testing.T) {
	handler := NewHandler(Config{Providers: []engine.Provider{testProvider{"s3", engine.AccessListable}}, Words: []string{"dev"}})
	request := localRequest(http.MethodPost, "/api/scan", `{"keyword":"acme","providers":["s3"]}`)
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
	request := localRequest(http.MethodPost, "/api/scan", `{"keyword":"acme","providers":["gcs"]}`)
	recorder := httptest.NewRecorder()
	handler.ServeHTTP(recorder, request)
	if recorder.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want 400", recorder.Code)
	}
}

func TestHealthEndpoint(t *testing.T) {
	handler := NewHandler(Config{})
	request := localRequest(http.MethodGet, "/healthz", "")
	recorder := httptest.NewRecorder()
	handler.ServeHTTP(recorder, request)
	if recorder.Code != http.StatusOK || recorder.Body.String() != "{\"status\":\"ok\"}\n" {
		t.Fatalf("health response = %d %q", recorder.Code, recorder.Body.String())
	}
}

func TestScanEndpointRejectsUnknownFieldsAndTrailingJSON(t *testing.T) {
	handler := NewHandler(Config{Providers: []engine.Provider{testProvider{"s3", engine.AccessPrivate}}})
	for _, body := range []string{`{"keyword":"acme","unexpected":true}`, `{"keyword":"acme"}{"keyword":"again"}`} {
		request := localRequest(http.MethodPost, "/api/scan", body)
		recorder := httptest.NewRecorder()
		handler.ServeHTTP(recorder, request)
		if recorder.Code != http.StatusBadRequest {
			t.Errorf("body %q returned status %d, want 400", body, recorder.Code)
		}
	}
}

func TestRejectsNonLoopbackHost(t *testing.T) {
	handler := NewHandler(Config{})
	request := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	request.Host = "example.com"
	recorder := httptest.NewRecorder()
	handler.ServeHTTP(recorder, request)
	if recorder.Code != http.StatusForbidden {
		t.Fatalf("status = %d, want 403", recorder.Code)
	}
}

func TestRejectsCrossOriginRequest(t *testing.T) {
	handler := NewHandler(Config{})
	request := localRequest(http.MethodGet, "/healthz", "")
	request.Header.Set("Origin", "https://example.com")
	recorder := httptest.NewRecorder()
	handler.ServeHTTP(recorder, request)
	if recorder.Code != http.StatusForbidden {
		t.Fatalf("status = %d, want 403", recorder.Code)
	}
}

func TestAllowsLoopbackOriginWithoutWildcard(t *testing.T) {
	handler := NewHandler(Config{})
	request := localRequest(http.MethodGet, "/healthz", "")
	request.Header.Set("Origin", "http://localhost:3000")
	recorder := httptest.NewRecorder()
	handler.ServeHTTP(recorder, request)
	if recorder.Code != http.StatusOK || recorder.Header().Get("Access-Control-Allow-Origin") != "http://localhost:3000" {
		t.Fatalf("response = %d, CORS = %q", recorder.Code, recorder.Header().Get("Access-Control-Allow-Origin"))
	}
}

func TestRejectsCandidateLimit(t *testing.T) {
	handler := NewHandler(Config{Providers: []engine.Provider{largeProvider{}}, MaxCandidates: 1})
	request := localRequest(http.MethodPost, "/api/scan", `{"keyword":"acme"}`)
	recorder := httptest.NewRecorder()
	handler.ServeHTTP(recorder, request)
	if recorder.Code != http.StatusRequestEntityTooLarge {
		t.Fatalf("status = %d, want 413", recorder.Code)
	}
}

func TestRejectsOversizedBody(t *testing.T) {
	handler := NewHandler(Config{})
	body := `{"keyword":"` + strings.Repeat("a", 1<<20) + `"}`
	request := localRequest(http.MethodPost, "/api/scan", body)
	recorder := httptest.NewRecorder()
	handler.ServeHTTP(recorder, request)
	if recorder.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want 400", recorder.Code)
	}
}
