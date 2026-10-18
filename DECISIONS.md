# Decisions

## Embedded default wordlist

The CLI embeds `cmd/kickbucket/wordlists/default.txt` with `go:embed`. A scan without `--wordlist` must work from any working directory, including when the binary is launched from `/tmp`; an external relative path would make the default behavior depend on the caller's directory.

An explicit `--wordlist` path still reads from the filesystem so users can bring their own lists.

## Standard flag parsing with flexible placement

The scan command uses Go's standard `flag.FlagSet` for validation and flag semantics. Go's parser stops at the first positional argument, while the documented command shape allows flags before or after the keyword, so recognized flags are normalized before passing the arguments to `FlagSet`.

## Cancellation and partial results

Scan requests use `signal.NotifyContext` for SIGINT and SIGTERM. Results collected before cancellation or a provider-level failure are rendered before the command exits. This avoids losing useful findings when a later provider or request fails.

## Read-only provider probes

Providers use anonymous `HEAD` or metadata/list requests only. They never download object contents or modify buckets. Public-read detection is opt-in through `--object-path`, because a bucket name alone does not identify an object that can be read.

## Local serve API

`serve` was added in this PR to give the web planner a local, programmatic entry point to the same engine used by the CLI. It exposes `POST /api/scan` and `GET /healthz` on loopback, rather than duplicating provider logic in the frontend.

The API requires a loopback `Host`, rejects non-loopback `Origin` values, caps request bodies at 1 MiB, and limits candidate work per request. It is intentionally not a remotely exposed service and has no authentication layer.

## Provider errors and summaries

A failed candidate is retained as a result with an `error` field so JSON and human output can explain what happened. Summary counts include only rows without errors; exposure counts remain limited to listable and public-read access levels.
