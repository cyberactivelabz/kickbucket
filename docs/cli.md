# CLI usage

```sh
kickbucket scan <keyword>
```

Kickbucket builds candidate bucket names from the keyword and a wordlist, checks each cloud provider, and prints the access level for every bucket it finds.

Run `kickbucket --help` and `kickbucket scan --help` for the flags in your version.

## Flags

- `--wordlist <path>` reads candidate suffixes from a file. It defaults to `wordlists/default.txt`.
- `--providers <list>` selects a comma-separated subset of `s3`, `gcs`, and `azure`. It defaults to all three.
- `--object-path <path>` probes a known object path for anonymous public-read access.
- `--json` emits results as JSON on stdout.

The authorization notice and errors are written to stderr. Use Kickbucket only on cloud storage you own or are explicitly authorized to test.

## Local API

Run `kickbucket serve` to expose the read-only scan API on `127.0.0.1:8080`. The server only accepts loopback addresses so it cannot be exposed accidentally without adding an authentication layer.

- `--addr <address>` changes the listen address.
- `--wordlist <path>` selects the wordlist used by API requests.
- `--object-path <path>` probes that object path for public-read access.

Send `POST /api/scan` with `{"keyword":"acme","providers":["s3","gcs","azure"]}`. The `providers` field is optional and the response is the same JSON result shape as `kickbucket scan --json`.

Use `GET /healthz` for a lightweight readiness check. Request bodies are limited to 1 MiB and reject unknown JSON fields.
