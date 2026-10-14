package s3

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/cyberactivelabz/kickbucket/internal/engine"
)

func TestProviderAccessLevels(t *testing.T) {
	tests := []struct {
		name   string
		status int
		want   engine.Access
	}{
		{"listable", http.StatusOK, engine.AccessListable},
		{"private", http.StatusForbidden, engine.AccessPrivate},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != http.MethodGet || r.URL.Query().Get("max-keys") != "0" {
					t.Errorf("unexpected request: %s %s", r.Method, r.URL)
				}
				w.WriteHeader(test.status)
			}))
			defer server.Close()
			access, err := New(Config{Endpoint: server.URL}).Check(context.Background(), "acme")
			if err != nil || access != test.want {
				t.Fatalf("Check() = %v, %v; want %v", access, err, test.want)
			}
		})
	}
}

func TestProviderPublicRead(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodHead && r.URL.Path == "/acme/object" {
			w.WriteHeader(http.StatusOK)
			return
		}
		w.WriteHeader(http.StatusForbidden)
	}))
	defer server.Close()
	access, err := New(Config{Endpoint: server.URL, ObjectPath: "/object"}).Check(context.Background(), "acme")
	if err != nil || access != engine.AccessPublicRead {
		t.Fatalf("Check() = %v, %v; want public read", access, err)
	}
}

func TestProviderRetriesThrottling(t *testing.T) {
	attempts := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		attempts++
		if attempts < 3 {
			w.WriteHeader(http.StatusServiceUnavailable)
			return
		}
		w.WriteHeader(http.StatusForbidden)
	}))
	defer server.Close()
	access, err := New(Config{Endpoint: server.URL}).Check(context.Background(), "acme")
	if err != nil || access != engine.AccessPrivate || attempts != 3 {
		t.Fatalf("Check() = %v, %v after %d attempts", access, err, attempts)
	}
}
