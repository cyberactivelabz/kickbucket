# Providers

Kickbucket checks each candidate bucket name against every selected provider. Each provider lives in its own package under `internal/engine/providers/<name>/` and implements the `Provider` interface in `internal/engine/provider.go`.

| Provider | `--providers` name | Page |
|---|---|---|
| Amazon S3 | `s3` | [Amazon S3](/providers/s3) |
| Google Cloud Storage | `gcs` | [Google Cloud Storage](/providers/gcs) |
| Azure Blob Storage | `azure` | [Azure Blob Storage](/providers/azure) |

## Candidate names

All providers use the same candidate generator. The keyword is lowercased and trimmed, then for each wordlist entry Kickbucket tries `<keyword>-<word>` and `<word>-<keyword>`. The bare keyword is always tried first, leading and trailing hyphens are stripped, and duplicates are dropped.

## How access is detected

Every provider sends anonymous, read-only requests. It never downloads object contents or changes anything.

1. **Object probe (optional).** With `--object-path <path>`, the provider first sends `HEAD` for that object. A `2xx` response means `public read` and the check stops there.
2. **Bucket probe.** Otherwise the provider sends an anonymous `GET` metadata/list request that asks for zero results:

| Response | Result |
|---|---|
| `2xx` | `listable` |
| `401` or `403` | `private` |
| `404` | not found error |
| anything else | error with the status code |

Responses `429` and `503` are retried up to three attempts in total, with a short backoff (50 ms, then 100 ms). Network errors (including DNS failures) and timeouts (10 seconds by default) are reported as errors.

The CLI hides results whose error contains "not found" unless `--verbose` is set. Other errors, such as DNS failures, are always shown. The HTTP API (`kickbucket serve`) returns every result, including not-found errors.
