# Kickbucket: instructions for Copilot

Kickbucket is an open-source (MIT) scanner for exposed cloud storage buckets, "the Nmap for cloud storage". Given a keyword, it generates candidate bucket names from a wordlist, checks them concurrently across S3, GCS and Azure, and reports each bucket as private, listable or public read.

## Layout
- cmd/kickbucket: CLI entry point (version vars are injected by GoReleaser)
- internal/engine: scan engine and the Provider interface (provider.go)
- internal/engine/providers/<name>: one package per cloud provider
- pkg/kickbucket: public Go library (to be extracted in month 3)
- web/: React frontend, served by `kickbucket serve` (not started)
- docs/: VitePress docs; website/: Next.js + Tailwind landing page
- wordlists/: plain text, one lowercase word per line

## Conventions
- Idiomatic Go: small interfaces, context.Context first parameter, errors wrapped with %w, no panics in library code.
- Run gofmt, go vet ./... and go test ./... before finishing any change.
- Tests use mocked HTTP servers (httptest). Never hit real cloud endpoints in tests.
- Default to safe behavior: rate limits and timeouts on by default.
- Keep stdout clean in JSON or piped mode; logs and notices go to stderr.
- Commit messages: short, imperative. GoReleaser filters out docs:, test: and chore: prefixes.

## Current status
Skeleton only: main.go prints --version, and engine/provider.go defines Access, Result and the Provider interface. Nothing else is implemented.

## Month 1 goal (a CLI worth starring)
1. Engine core: candidate-name generator, concurrent worker pool with context cancellation, rate limiter, timeouts.
2. S3 provider: detect public read, listable and private; handle redirects, region errors and throttling.
3. GCS and Azure providers behind the same interface.
4. `kickbucket scan <keyword>` with plain and JSON output, --wordlist flag, and an authorization notice.
5. Keep docs/ in sync with the real flags.

## Rules
- Use only on assets the user owns or is authorized to test. Never add code that downloads or modifies bucket contents.
- Do not invent CLI flags or APIs in the docs; document what exists.

## Names and paths
- Tool and command name: kickbucket. Go module: github.com/cyberactivelabz/kickbucket. Entry point: cmd/kickbucket.
- Repo: github.com/cyberactivelabz/kickbucket. Docs are served under the /kickbucket/ base path on GitHub Pages.
- Work on a branch per task and open a PR. Keep commits small.
