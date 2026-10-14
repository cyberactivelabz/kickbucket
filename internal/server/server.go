// Package server exposes the local, read-only scan API.
package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/cyberactivelabz/kickbucket/internal/engine"
)

type Config struct {
	Providers []engine.Provider
	Words     []string
	Workers   int
	Rate      float64
}

type Handler struct {
	providers map[string]engine.Provider
	words     []string
	workers   int
	rate      float64
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
	return &Handler{providers: providerMap, words: config.Words, workers: config.Workers, rate: config.Rate}
}

func (h *Handler) ServeHTTP(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
	response.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	if request.Method == http.MethodOptions {
		response.WriteHeader(http.StatusNoContent)
		return
	}
	if request.URL.Path != "/api/scan" || request.Method != http.MethodPost {
		writeError(response, http.StatusNotFound, "endpoint not found")
		return
	}
	var input scanRequest
	if err := json.NewDecoder(request.Body).Decode(&input); err != nil {
		writeError(response, http.StatusBadRequest, "invalid JSON body")
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
		limiter := engine.NewLimiter(h.rate)
		providerResults, scanErr := engine.Scan(request.Context(), provider, provider.Candidates(input.Keyword, h.words), h.workers, limiter)
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

func (h *Handler) selectProviders(names []string) ([]engine.Provider, error) {
	if len(names) == 0 {
		names = make([]string, 0, len(h.providers))
		for name := range h.providers {
			names = append(names, name)
		}
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
