package gcs

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/cyberactivelabz/kickbucket/internal/engine"
)

func TestAccessLevels(t *testing.T) {
	for _, test := range []struct {
		name   string
		status int
		want   engine.Access
	}{
		{"listable", http.StatusOK, engine.AccessListable}, {"private", http.StatusForbidden, engine.AccessPrivate},
	} {
		t.Run(test.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(test.status) }))
			defer server.Close()
			got, err := New(Config{Endpoint: server.URL}).Check(context.Background(), "acme")
			if err != nil || got != test.want {
				t.Fatalf("Check() = %v, %v; want %v", got, err, test.want)
			}
		})
	}
}

func TestPublicRead(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodHead {
			w.WriteHeader(http.StatusOK)
		} else {
			w.WriteHeader(http.StatusForbidden)
		}
	}))
	defer server.Close()
	got, err := New(Config{Endpoint: server.URL, ObjectPath: "/object"}).Check(context.Background(), "acme")
	if err != nil || got != engine.AccessPublicRead {
		t.Fatalf("Check() = %v, %v; want public read", got, err)
	}
}
