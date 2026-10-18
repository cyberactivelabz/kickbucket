# Google Cloud Storage

Provider name: `gcs`. Source: `internal/engine/providers/gcs/gcs.go`.

## URL formats

Kickbucket uses the GCS JSON API base, with the candidate name as the bucket:

| Probe | Method | URL |
|---|---|---|
| Object (with `--object-path`) | `HEAD` | `https://storage.googleapis.com/storage/v1/b/<bucket>/<object-path>` |
| Bucket metadata/list | `GET` | `https://storage.googleapis.com/storage/v1/b/<bucket>?maxResults=0` |

The bucket name is path-escaped. Leading `/` characters are removed from the object path.

When the provider is configured with a custom endpoint (the `Endpoint` field of `gcs.Config`), the endpoint replaces `https://storage.googleapis.com/storage/v1/b`: `<endpoint>/<bucket>/<object-path>` and `<endpoint>/<bucket>?maxResults=0`.

## How access is detected

| Step | Response | Result |
|---|---|---|
| Object `HEAD` | `2xx` | `public read` |
| Object `HEAD` | anything else | continue to the bucket probe |
| Bucket `GET` | `2xx` | `listable` |
| Bucket `GET` | `401` or `403` | `private` |
| Bucket `GET` | `404` | error `bucket "<name>" not found` |
| Bucket `GET` | any other status | error `unexpected GCS status <code>` |

Responses `429` and `503` are retried up to three attempts in total. See [Providers](/providers/) for the shared behaviour.
