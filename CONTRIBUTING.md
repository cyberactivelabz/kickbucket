# Contributing to Kickbucket

Thanks for helping. The fastest ways to contribute are a new provider, a wordlist, or a bug fix.

## Setup

```sh
git clone https://github.com/cyberactivelabz/kickbucket && cd kickbucket
go test ./...
go run ./cmd/kickbucket --help
```

## Add a cloud provider

1. Open a "New provider" issue so we agree on scope first.
2. Create `internal/engine/providers/<name>/` and implement the `Provider` interface in `internal/engine/provider.go` (URL candidates for a name, and a check that returns the access level).
3. Register it in `internal/engine/providers/registry.go`.
4. Add tests with a mocked HTTP server. Tests must not hit real cloud endpoints.
5. Add a docs page under `docs/providers/<name>.md`.

## Add a wordlist

Put a plain text file (one word per line, lowercase, no duplicates) in `wordlists/`. Say in the PR where the list came from.

## Pull requests

- One change per PR; keep it small.
- Run `gofmt` and `go vet ./...` and `go test ./...` before pushing.
- Explain what was wrong and why the change fixes it.

## Contributor ladder

Contributor, then Provider author (merged a provider), then Maintainer (invited after sustained contributions). Maintainers share in sponsor funds.

Security issues: see [SECURITY.md](SECURITY.md), not public issues.
