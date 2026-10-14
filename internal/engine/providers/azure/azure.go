package azure

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/cyberactivelabz/kickbucket/internal/engine"
)

type Config struct {
	Endpoint   string
	ObjectPath string
	Client     *http.Client
	Timeout    time.Duration
}

type Provider struct {
	endpoint, object string
	client           *http.Client
}

func New(config Config) *Provider {
	timeout := config.Timeout
	if timeout <= 0 {
		timeout = 10 * time.Second
	}
	client := config.Client
	if client == nil {
		client = &http.Client{Timeout: timeout}
	}
	return &Provider{strings.TrimRight(config.Endpoint, "/"), config.ObjectPath, client}
}

func (p *Provider) Name() string { return "azure" }
func (p *Provider) Candidates(keyword string, words []string) []string {
	return engine.Candidates(keyword, words)
}

func (p *Provider) Check(ctx context.Context, bucket string) (engine.Access, error) {
	if p.object != "" {
		status, err := p.request(ctx, http.MethodHead, p.containerURL(bucket)+"/"+strings.TrimLeft(p.object, "/"))
		if err != nil {
			return engine.AccessUnknown, err
		}
		if status >= 200 && status < 300 {
			return engine.AccessPublicRead, nil
		}
	}
	status, err := p.request(ctx, http.MethodGet, p.containerURL(bucket)+"?restype=container&comp=list&maxresults=0")
	if err != nil {
		return engine.AccessUnknown, err
	}
	switch {
	case status >= 200 && status < 300:
		return engine.AccessListable, nil
	case status == http.StatusForbidden || status == http.StatusUnauthorized:
		return engine.AccessPrivate, nil
	case status == http.StatusNotFound:
		return engine.AccessUnknown, fmt.Errorf("container %q not found", bucket)
	default:
		return engine.AccessUnknown, fmt.Errorf("unexpected Azure status %d", status)
	}
}

func (p *Provider) containerURL(bucket string) string {
	base := p.endpoint
	if base == "" {
		base = "https://" + url.PathEscape(bucket) + ".blob.core.windows.net"
		return base
	}
	return base + "/" + url.PathEscape(bucket)
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
	return 0, fmt.Errorf("request retry exhausted")
}
