# Azure Blob Storage

Provider name: `azure`. Source: `internal/engine/providers/azure/azure.go`.

## URL formats

By default the candidate name is used as the **storage account** name, and requests go to that account's blob endpoint with no container segment in the path:

| Probe | Method | URL |
|---|---|---|
| Object (with `--object-path`) | `HEAD` | `https://<name>.blob.core.windows.net/<object-path>` |
| List | `GET` | `https://<name>.blob.core.windows.net?restype=container&comp=list&maxresults=0` |

The name is path-escaped. Leading `/` characters are removed from the object path. Because the name becomes a hostname, the object path should include the container, for example `--object-path public/index.html` probes the blob `index.html` in the container `public`.

When the provider is configured with a custom endpoint (the `Endpoint` field of `azure.Config`), the candidate name is used as the **container** under that endpoint instead: `<endpoint>/<name>/<object-path>` and `<endpoint>/<name>?restype=container&comp=list&maxresults=0`.

## How access is detected

| Step | Response | Result |
|---|---|---|
| Object `HEAD` | `2xx` | `public read` |
| Object `HEAD` | anything else | continue to the list probe |
| List `GET` (List Blobs, zero results) | `2xx` | `listable` |
| List `GET` | `401` or `403` | `private` |
| List `GET` | `404` | error `container "<name>" not found` |
| List `GET` | any other status | error `unexpected Azure status <code>` |

In default mode, a name with no matching storage account usually fails DNS resolution instead of returning `404`. That shows up as a network error, which the CLI displays even without `--verbose`.

Responses `429` and `503` are retried up to three attempts in total. See [Providers](/providers/) for the shared behaviour.
