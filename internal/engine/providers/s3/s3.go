package s3

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/cyberactivelabz/kickbucket/internal/engine"
)

const defaultTimeout = 10 * time.Second

type Config struct {
	Endpoint   string
	ObjectPath string
	Client     *http.Client
	Timeout    time.Duration
}

type Provider struct {
	endpoint   string
	objectPath string
	client     *http.Client
}

func New(config Config) *Provider {
	timeout := config.Timeout
	if timeout <= 0 {
		timeout = defaultTimeout
	}
	client := config.Client
	if client == nil {
		client = &http.Client{Timeout: timeout}
	}
	return &Provider{endpoint: strings.TrimRight(config.Endpoint, "/"), objectPath: config.ObjectPath, client: client}
}

func (p *Provider) Name() string { return "s3" }

func (p *Provider) Candidates(keyword string, words []string) []string {
	return engine.Candidates(keyword, words)
}

func (p *Provider) Check(ctx context.Context, bucket string) (engine.Access, error) {
	if bucket == "" {
		return engine.AccessUnknown, errors.New("bucket is required")
	}
	if p.objectPath != "" {
		status, err := p.request(ctx, http.MethodHead, p.bucketURL(bucket, p.objectPath))
		if err != nil {
			return engine.AccessUnknown, err
		}
		if status >= 200 && status < 300 {
			return engine.AccessPublicRead, nil
		}
	}
	status, err := p.request(ctx, http.MethodGet, p.bucketURL(bucket, "?list-type=2&max-keys=0"))
	if err != nil {
		return engine.AccessUnknown, err
	}
	switch {
	case status >= 200 && status < 300:
		return engine.AccessListable, nil
	case status == http.StatusForbidden || status == http.StatusUnauthorized:
		return engine.AccessPrivate, nil
	case status == http.StatusNotFound:
		return engine.AccessUnknown, fmt.Errorf("bucket %q not found", bucket)
	default:
		return engine.AccessUnknown, fmt.Errorf("unexpected S3 status %d", status)
	}
}

func (p *Provider) bucketURL(bucket, suffix string) string {
	if p.endpoint != "" {
		return p.endpoint + "/" + url.PathEscape(bucket) + suffix
	}
	return "https://" + url.PathEscape(bucket) + ".s3.amazonaws.com" + suffix
}

func (p *Provider) request(ctx context.Context, method, address string) (int, error) {
	for attempt := 0; attempt < 3; attempt++ {
		req, err := http.NewRequestWithContext(ctx, method, address, nil)
		if err != nil {
			return 0, err
		}
		resp, err := p.client.Do(req)
		if err != nil {
			return 0, err
		}
		_, _ = io.Copy(io.Discard, resp.Body)
		_ = resp.Body.Close()
		if resp.StatusCode != http.StatusTooManyRequests && resp.StatusCode != http.StatusServiceUnavailable {
			return resp.StatusCode, nil
		}
		if attempt == 2 {
			return resp.StatusCode, nil
		}
		select {
		case <-ctx.Done():
			return 0, ctx.Err()
		case <-time.After(time.Duration(attempt+1) * 50 * time.Millisecond):
		}
	}
	return 0, errors.New("request retry exhausted")
}
