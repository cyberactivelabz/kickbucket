# CLI usage

```sh
kickbucket scan <keyword>
```

Kickbucket builds candidate bucket names from the keyword and a wordlist, checks each cloud provider, and prints the access level for every bucket it finds.

Run `kickbucket --help` and `kickbucket scan --help` for the flags in your version.

## Flags

- `--wordlist <path>` reads candidate suffixes from a file. It defaults to `wordlists/default.txt`.
- `--providers <list>` selects a comma-separated subset of `s3`, `gcs`, and `azure`. It defaults to all three.
- `--json` emits results as JSON on stdout.

The authorization notice and errors are written to stderr. Use Kickbucket only on cloud storage you own or are explicitly authorized to test.
