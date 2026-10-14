package providers

import (
	"github.com/cyberactivelabz/kickbucket/internal/engine"
	"github.com/cyberactivelabz/kickbucket/internal/engine/providers/azure"
	"github.com/cyberactivelabz/kickbucket/internal/engine/providers/gcs"
	"github.com/cyberactivelabz/kickbucket/internal/engine/providers/s3"
)

type Config struct {
	S3    s3.Config
	GCS   gcs.Config
	Azure azure.Config
}

func NewRegistry(config Config) []engine.Provider {
	return []engine.Provider{s3.New(config.S3), gcs.New(config.GCS), azure.New(config.Azure)}
}
