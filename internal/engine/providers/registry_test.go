package providers

import "testing"

func TestNewRegistry(t *testing.T) {
	got := NewRegistry(Config{})
	if len(got) != 3 || got[0].Name() != "s3" || got[1].Name() != "gcs" || got[2].Name() != "azure" {
		t.Fatalf("registry providers = %#v", got)
	}
}
