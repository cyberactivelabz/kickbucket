// Package server exposes the local, read-only scan API.
package server

import (
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"sort"
	"strings"

	"github.com/cyberactivelabz/kickbucket/internal/engine"
)

type Config struct {
	Providers     []engine.Provider
	Words         []string
	Workers       int
	Rate          float64
	MaxCandidates int
}

type Handler struct {
	providers     map[string]engine.Provider
	words         []string
	workers       int
	rate          float64
	maxCandidates int
}

type scanRequest struct {
	Keyword   string   `json:"keyword"`
	Providers []string `json:"providers"`
}

func NewHandler(config Config) http.Handler {
	providerMap := make(map[string]engine.Provider, len(config.Providers))
	for _, provider := range config.Providers {
		if provider != nil {
			providerMap[provider.Name()] = provider
		}
	}
	if config.Workers < 1 {
		config.Workers = 8
	}
	if config.Rate <= 0 {
		config.Rate = 10
	}
	if config.MaxCandidates <= 0 {
		config.MaxCandidates = 5000
	}
	return &Handler{providers: providerMap, words: config.Words, workers: config.Workers, rate: config.Rate, maxCandidates: config.MaxCandidates}
}

func (h *Handler) ServeHTTP(response http.ResponseWriter, request *http.Request) {
	if !isLoopbackHost(request.Host) {
		writeError(response, http.StatusForbidden, "host must be loopback")
		return
	}
	origin := request.Header.Get("Origin")
	if origin != "" {
		if !isLoopbackOrigin(origin) {
			writeError(response, http.StatusForbidden, "cross-origin request denied")
			return
		}
		response.Header().Set("Access-Control-Allow-Origin", origin)
		response.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		response.Header().Add("Vary", "Origin")
	}
	if request.Method == http.MethodOptions {
		response.WriteHeader(http.StatusNoContent)
		return
	}
	if request.URL.Path == "/healthz" && request.Method == http.MethodGet {
		response.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(response).Encode(map[string]string{"status": "ok"})
		return
	}
	if request.URL.Path != "/api/scan" || request.Method != http.MethodPost {
		writeError(response, http.StatusNotFound, "endpoint not found")
		return
	}
	request.Body = http.MaxBytesReader(response, request.Body, 1<<20)
	var input scanRequest
	decoder := json.NewDecoder(request.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&input); err != nil {
		writeError(response, http.StatusBadRequest, "invalid JSON body")
		return
	}
	var extra any
	if err := decoder.Decode(&extra); err != io.EOF {
		writeError(response, http.StatusBadRequest, "request body must contain one JSON object")
		return
	}
	input.Keyword = strings.TrimSpace(input.Keyword)
	if input.Keyword == "" {
		writeError(response, http.StatusBadRequest, "keyword is required")
		return
	}

	selected, err := h.selectProviders(input.Providers)
	if err != nil {
		writeError(response, http.StatusBadRequest, err.Error())
		return
	}
	results := make([]engine.Result, 0)
	for _, provider := range selected {
		candidates := provider.Candidates(input.Keyword, h.words)
		if len(candidates) > h.maxCandidates {
			writeError(response, http.StatusRequestEntityTooLarge, "candidate limit exceeded")
			return
		}
		limiter := engine.NewLimiter(h.rate)
		providerResults, scanErr := engine.Scan(request.Context(), provider, candidates, h.workers, limiter)
		limiter.Stop()
		if scanErr != nil {
			writeError(response, http.StatusRequestTimeout, scanErr.Error())
			return
		}
		results = append(results, providerResults...)
	}
	response.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(response).Encode(results)
}

func isLoopbackHost(hostport string) bool {
	host := hostport
	if parsedHost, _, err := net.SplitHostPort(hostport); err == nil {
		host = parsedHost
	}
	if strings.EqualFold(host, "localhost") {
		return true
	}
	ip := net.ParseIP(host)
	return ip != nil && ip.IsLoopback()
}

func isLoopbackOrigin(origin string) bool {
	parsed, err := url.Parse(origin)
	if err != nil || parsed.User != nil || parsed.Path != "" || parsed.RawQuery != "" || parsed.Fragment != "" {
		return false
	}
	return parsed.Scheme == "http" && isLoopbackHost(parsed.Host)
}

func (h *Handler) selectProviders(names []string) ([]engine.Provider, error) {
	if len(names) == 0 {
		names = make([]string, 0, len(h.providers))
		for name := range h.providers {
			names = append(names, name)
		}
		sort.Strings(names)
	}
	selected := make([]engine.Provider, 0, len(names))
	seen := make(map[string]bool, len(names))
	for _, name := range names {
		name = strings.ToLower(strings.TrimSpace(name))
		provider, ok := h.providers[name]
		if !ok {
			return nil, fmt.Errorf("unknown provider %q", name)
		}
		if !seen[name] {
			selected = append(selected, provider)
			seen[name] = true
		}
	}
	return selected, nil
}

func writeError(response http.ResponseWriter, status int, message string) {
	response.Header().Set("Content-Type", "application/json")
	response.WriteHeader(status)
	_ = json.NewEncoder(response).Encode(map[string]string{"error": message})
}
