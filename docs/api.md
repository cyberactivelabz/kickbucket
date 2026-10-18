# Go API

Kickbucket does not expose a public Go library yet. The current programmatic interface is the local HTTP API started by `kickbucket serve`.

## Start the API

```sh
kickbucket serve --object-path public/index.html
```

## Scan

```sh
curl -sS http://127.0.0.1:8080/api/scan \
  -H 'Content-Type: application/json' \
  -d '{"keyword":"acme","providers":["s3","gcs","azure"]}'
```

The `providers` field is optional. Results contain `provider`, `name`, `access`, and an optional `error` field. Use `GET /healthz` to check readiness. The API is loopback-only and performs anonymous, read-only probes.
