package engine

import (
	"reflect"
	"testing"
)

func TestCandidates(t *testing.T) {
	got := Candidates(" Acme ", []string{"prod", "prod", "", " backup "})
	want := []string{"acme", "acme-prod", "prod-acme", "acme-backup", "backup-acme"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("Candidates() = %#v, want %#v", got, want)
	}
}
