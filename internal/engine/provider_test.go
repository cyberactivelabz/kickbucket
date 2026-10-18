package engine

import (
	"encoding/json"
	"testing"
)

func TestAccessJSONUsesStrings(t *testing.T) {
	for access, want := range map[Access]string{
		AccessUnknown:    `"unknown"`,
		AccessPrivate:    `"private"`,
		AccessListable:   `"listable"`,
		AccessPublicRead: `"public read"`,
	} {
		got, err := json.Marshal(access)
		if err != nil {
			t.Fatalf("json.Marshal(%v): %v", access, err)
		}
		if string(got) != want {
			t.Errorf("json.Marshal(%v) = %s, want %s", access, got, want)
		}
	}
}

func TestResultJSONUsesAccessString(t *testing.T) {
	got, err := json.Marshal(Result{Provider: "s3", Name: "acme-private", Access: AccessPrivate})
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != `{"provider":"s3","name":"acme-private","access":"private"}` {
		t.Fatalf("json.Marshal(Result) = %s", got)
	}
}

func TestResultJSONRoundTripsAccessString(t *testing.T) {
	input := []byte(`{"provider":"gcs","name":"acme-assets","access":"listable"}`)
	var result Result
	if err := json.Unmarshal(input, &result); err != nil {
		t.Fatal(err)
	}
	if result.Access != AccessListable {
		t.Fatalf("decoded access = %v, want listable", result.Access)
	}
}
