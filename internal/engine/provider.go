// Package engine holds the scan engine and the provider interface.
package engine

import (
	"context"
	"encoding/json"
	"fmt"
)

// Access describes how exposed a bucket is.
type Access int

const (
	AccessUnknown Access = iota
	AccessPrivate
	AccessListable
	AccessPublicRead
)

func (a Access) String() string {
	switch a {
	case AccessPrivate:
		return "private"
	case AccessListable:
		return "listable"
	case AccessPublicRead:
		return "public read"
	default:
		return "unknown"
	}
}

func (a Access) MarshalJSON() ([]byte, error) {
	return json.Marshal(a.String())
}

func (a *Access) UnmarshalJSON(data []byte) error {
	var value string
	if err := json.Unmarshal(data, &value); err != nil {
		return err
	}
	switch value {
	case "unknown":
		*a = AccessUnknown
	case "private":
		*a = AccessPrivate
	case "listable":
		*a = AccessListable
	case "public read":
		*a = AccessPublicRead
	default:
		return fmt.Errorf("unknown access level %q", value)
	}
	return nil
}

// Result is one bucket found by a scan.
type Result struct {
	Provider string `json:"provider"`
	Name     string `json:"name"`
	Access   Access `json:"access"`
	Error    string `json:"error,omitempty"`
}

// Provider is implemented once per cloud storage service.
type Provider interface {
	Name() string
	// Candidates builds bucket names to try from a keyword and wordlist.
	Candidates(keyword string, words []string) []string
	// Check reports the access level of one bucket.
	Check(ctx context.Context, bucket string) (Access, error)
}
