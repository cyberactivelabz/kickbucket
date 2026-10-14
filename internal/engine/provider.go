// Package engine holds the scan engine and the provider interface.
package engine

import "context"

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
