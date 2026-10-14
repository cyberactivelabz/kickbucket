package main

import (
	"reflect"
	"testing"
)

func TestParseScanArgsDefaults(t *testing.T) {
	got, err := parseScanArgs([]string{"acme"})
	if err != nil {
		t.Fatal(err)
	}
	want := scanOptions{keyword: "acme", wordlistPath: "wordlists/default.txt", providers: []string{"s3", "gcs", "azure"}}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("parseScanArgs() = %#v, want %#v", got, want)
	}
}

func TestParseScanArgsFlags(t *testing.T) {
	got, err := parseScanArgs([]string{"acme", "--json", "--wordlist", "words.txt", "--providers", "GCS, azure"})
	if err != nil {
		t.Fatal(err)
	}
	if got.keyword != "acme" || got.wordlistPath != "words.txt" || !got.json || !reflect.DeepEqual(got.providers, []string{"gcs", "azure"}) {
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
