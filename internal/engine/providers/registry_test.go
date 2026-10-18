package providers

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/cyberactivelabz/kickbucket/internal/engine"
	"github.com/cyberactivelabz/kickbucket/internal/engine/providers/azure"
	"github.com/cyberactivelabz/kickbucket/internal/engine/providers/gcs"
	"github.com/cyberactivelabz/kickbucket/internal/engine/providers/s3"
)

func TestNewRegistry(t *testing.T) {
	got := NewRegistry(Config{})
	if len(got) != 3 || got[0].Name() != "s3" || got[1].Name() != "gcs" || got[2].Name() != "azure" {
		t.Fatalf("registry providers = %#v", got)
	}
}

func TestNewRegistryForwardsObjectPath(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		if request.Method == http.MethodHead && request.URL.Path == "/acme/public.txt" {
			response.WriteHeader(http.StatusOK)
			return
		}
		response.WriteHeader(http.StatusForbidden)
	}))
	defer server.Close()
	registry := NewRegistry(Config{
		ObjectPath: "public.txt",
		S3:         s3.Config{Endpoint: server.URL},
		GCS:        gcs.Config{Endpoint: server.URL},
		Azure:      azure.Config{Endpoint: server.URL},
	})
	for _, provider := range registry {
		access, err := provider.Check(context.Background(), "acme")
		if err != nil || access != engine.AccessPublicRead {
			t.Errorf("%s Check() = %v, %v; want public read", provider.Name(), access, err)
		}
	}
}
