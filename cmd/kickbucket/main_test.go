package main

import (
	"reflect"
	"testing"

	"github.com/cyberactivelabz/kickbucket/internal/engine"
)

func TestParseScanArgsDefaults(t *testing.T) {
	got, err := parseScanArgs([]string{"acme"})
	if err != nil {
		t.Fatal(err)
	}
	want := scanOptions{keyword: "acme", providers: []string{"s3", "gcs", "azure"}}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("parseScanArgs() = %#v, want %#v", got, want)
	}
}

func TestParseScanArgsFlags(t *testing.T) {
	got, err := parseScanArgs([]string{"acme", "--json", "--wordlist", "words.txt", "--providers", "GCS, azure", "--object-path", "public/index.html"})
	if err != nil {
		t.Fatal(err)
	}
	if got.keyword != "acme" || got.wordlistPath != "words.txt" || got.objectPath != "public/index.html" || !got.json || !reflect.DeepEqual(got.providers, []string{"gcs", "azure"}) {
		t.Fatalf("parseScanArgs() = %#v", got)
	}
}

func TestParseScanArgsFlagsBeforeKeyword(t *testing.T) {
	got, err := parseScanArgs([]string{"--json", "--providers", "gcs", "acme"})
	if err != nil {
		t.Fatal(err)
	}
	if got.keyword != "acme" || !got.json || !reflect.DeepEqual(got.providers, []string{"gcs"}) {
		t.Fatalf("parseScanArgs() = %#v", got)
	}
}

func TestParseScanArgsRejectsInvalidInput(t *testing.T) {
	for _, args := range [][]string{{}, {"acme", "--providers"}, {"acme", "--providers", "digitalocean"}, {"acme", "--wat"}} {
		if _, err := parseScanArgs(args); err == nil {
			t.Errorf("parseScanArgs(%q) returned nil error", args)
		}
	}
}

func TestParseServeArgs(t *testing.T) {
	got, err := parseServeArgs([]string{"--addr", "localhost:9090", "--wordlist", "custom.txt", "--object-path", "public/index.html"})
	if err != nil {
		t.Fatal(err)
	}
	want := serveOptions{address: "localhost:9090", wordlistPath: "custom.txt", objectPath: "public/index.html"}
	if got != want {
		t.Fatalf("parseServeArgs() = %#v, want %#v", got, want)
	}
}

func TestValidateServeAddress(t *testing.T) {
	for _, address := range []string{"127.0.0.1:8080", "localhost:9090", "[::1]:8080", ":8080"} {
		if err := validateServeAddress(address); err != nil {
			t.Errorf("validateServeAddress(%q) = %v", address, err)
		}
	}
	for _, address := range []string{"0.0.0.0:8080", "192.168.1.10:8080", "localhost"} {
		if err := validateServeAddress(address); err == nil {
			t.Errorf("validateServeAddress(%q) returned nil", address)
		}
	}
}

func TestCountFoundExcludesErrorRows(t *testing.T) {
	results := []engine.Result{{Name: "private"}, {Name: "missing", Error: "not found"}, {Name: "listable"}}
	if got := countFound(results); got != 2 {
		t.Fatalf("countFound() = %d, want 2", got)
	}
}
