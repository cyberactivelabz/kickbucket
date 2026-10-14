#!/usr/bin/env bash
# BucketMap bootstrap for GitHub Codespaces (or any Linux shell).
# Usage: bash setup.sh [--no-install] [--push]
#   --no-install  skip npm install in docs/ and website/
#   --push        create the GitHub repo with gh and push (needs `gh auth login`; Codespaces has it)
set -euo pipefail

INSTALL=1; PUSH=0
for a in "$@"; do
  case "$a" in
    --no-install) INSTALL=0 ;;
    --push) PUSH=1 ;;
    *) echo "unknown flag: $a" >&2; exit 2 ;;
  esac
done

if [ -n "${GH_USER:-}" ]; then :
elif command -v gh >/dev/null 2>&1 && gh auth status >/dev/null 2>&1; then GH_USER="$(gh api user -q .login)"
else echo "Set GH_USER first, e.g. GH_USER=yourname bash setup.sh" >&2; exit 1; fi

[ -e bucketmap ] && { echo "./bucketmap already exists; remove it or run elsewhere" >&2; exit 1; }
mkdir bucketmap && cd bucketmap
echo "Creating BucketMap for github.com/$GH_USER ..."

mkdir -p '.github'
cat > '.github/FUNDING.yml' <<'__BM_EOF__'
github: [YOUR_GH_USER]
buy_me_a_coffee: bucketmap
open_collective: bucketmap
__BM_EOF__
mkdir -p '.github/ISSUE_TEMPLATE'
cat > '.github/ISSUE_TEMPLATE/bug_report.yml' <<'__BM_EOF__'
name: Bug report
description: Something is broken or wrong
labels: [bug]
body:
  - type: input
    id: version
    attributes: { label: BucketMap version (bucketmap --version) }
    validations: { required: true }
  - type: input
    id: os
    attributes: { label: OS and architecture }
  - type: textarea
    id: steps
    attributes: { label: Command you ran and what happened, description: Paste output. Remove any target names you cannot share. }
    validations: { required: true }
  - type: textarea
    id: expected
    attributes: { label: What you expected }
__BM_EOF__
mkdir -p '.github/ISSUE_TEMPLATE'
cat > '.github/ISSUE_TEMPLATE/feature_request.yml' <<'__BM_EOF__'
name: Feature request
description: Suggest an improvement
labels: [enhancement]
body:
  - type: textarea
    id: problem
    attributes: { label: What are you trying to do? }
    validations: { required: true }
  - type: textarea
    id: idea
    attributes: { label: What would you like BucketMap to do? }
__BM_EOF__
mkdir -p '.github/ISSUE_TEMPLATE'
cat > '.github/ISSUE_TEMPLATE/new_provider.yml' <<'__BM_EOF__'
name: New provider
description: Propose support for another cloud storage provider
labels: [provider]
body:
  - type: input
    id: provider
    attributes: { label: Provider name }
    validations: { required: true }
  - type: textarea
    id: urls
    attributes: { label: Bucket URL formats and how public access is detected, description: Link the provider's docs. }
    validations: { required: true }
  - type: dropdown
    id: build
    attributes:
      label: Will you implement it?
      options: ["Yes", "No, I'm requesting it"]
__BM_EOF__
mkdir -p '.github'
cat > '.github/PULL_REQUEST_TEMPLATE.md' <<'__BM_EOF__'
## What changed and why

## How I tested it

## Checklist
- [ ] `gofmt`, `go vet ./...`, `go test ./...` pass
- [ ] Docs updated (if behavior changed)
- [ ] No tests hit real cloud endpoints
__BM_EOF__
mkdir -p '.github/workflows'
cat > '.github/workflows/docs-deploy.yml' <<'__BM_EOF__'
name: Deploy docs
on:
  push:
    branches: [main]
    paths: [docs/**, .github/workflows/docs-deploy.yml]
permissions:
  contents: read
  pages: write
  id-token: write
jobs:
  deploy:
    runs-on: ubuntu-latest
    environment:
      name: github-pages
      url: ${{ steps.deployment.outputs.page_url }}
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-node@v4
        with: { node-version: 20, cache: npm, cache-dependency-path: docs/package-lock.json }
      - run: npm ci
        working-directory: docs
      - run: npm run build
        working-directory: docs
      - uses: actions/configure-pages@v5
      - uses: actions/upload-pages-artifact@v3
        with: { path: docs/.vitepress/dist }
      - id: deployment
        uses: actions/deploy-pages@v4
__BM_EOF__
mkdir -p '.github/workflows'
cat > '.github/workflows/release.yml' <<'__BM_EOF__'
name: Release
on:
  push:
    tags: ["v*"]
permissions:
  contents: write
jobs:
  goreleaser:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
        with:
          fetch-depth: 0 # GoReleaser needs full history for the changelog
      - uses: actions/setup-go@v5
        with:
          go-version: stable
      - uses: goreleaser/goreleaser-action@v6
        with:
          version: "~> v2"
          args: release --clean
        env:
          GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}
__BM_EOF__
cat > '.gitignore' <<'__BM_EOF__'
node_modules/
.next/
docs/.vitepress/dist/
docs/.vitepress/cache/
__BM_EOF__
cat > '.goreleaser.yaml' <<'__BM_EOF__'
version: 2

before:
  hooks:
    - go mod tidy

builds:
  - id: bucketmap
    main: ./cmd/bucketmap
    binary: bucketmap
    env:
      - CGO_ENABLED=0
    goos: [linux, darwin, windows]
    goarch: [amd64, arm64]
    # Needs `var version = "dev"` in package main of cmd/bucketmap.
    ldflags:
      - -s -w -X main.version={{.Version}} -X main.commit={{.ShortCommit}} -X main.date={{.Date}}

archives:
  - formats: [tar.gz]
    name_template: "bucketmap_{{ .Version }}_{{ .Os }}_{{ .Arch }}"
    format_overrides:
      - goos: windows
        formats: [zip]
    files: [LICENSE, README.md]

checksum:
  name_template: checksums.txt

changelog:
  sort: asc
  filters:
    exclude: ["^docs:", "^test:", "^chore:"]
__BM_EOF__
cat > 'CODE_OF_CONDUCT.md' <<'__BM_EOF__'
# Code of Conduct

This project follows the [Contributor Covenant, version 2.1](https://www.contributor-covenant.org/version/2/1/code_of_conduct/).

Report violations to conduct@bucketmap.dev. Maintainers will review every report and keep the reporter's identity private.
__BM_EOF__
cat > 'CONTRIBUTING.md' <<'__BM_EOF__'
# Contributing to BucketMap

Thanks for helping. The fastest ways to contribute are a new provider, a wordlist, or a bug fix.

## Setup

```sh
git clone https://github.com/YOUR_GH_USER/bucketmap && cd bucketmap
go test ./...
go run ./cmd/bucketmap --help
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
__BM_EOF__
cat > 'LICENSE' <<'__BM_EOF__'
MIT License

Copyright (c) 2026 BucketMap contributors

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
__BM_EOF__
cat > 'README.md' <<'__BM_EOF__'
<p align="center">
  <a href="https://bucketmap.dev"><img src="logo.png" width="200" alt="BucketMap"></a>
</p>

<h3 align="center">The Nmap for cloud storage.</h3>
<p align="center">Find exposed S3, GCS and Azure buckets fast. Open source, written in Go.</p>

<p align="center">
  <a href="https://github.com/YOUR_GH_USER/bucketmap/stargazers"><img src="https://img.shields.io/github/stars/YOUR_GH_USER/bucketmap?style=flat" alt="Stars"></a>
  <a href="LICENSE"><img src="https://img.shields.io/github/license/YOUR_GH_USER/bucketmap" alt="MIT License"></a>
  <a href="https://github.com/YOUR_GH_USER/bucketmap/releases"><img src="https://img.shields.io/github/v/release/YOUR_GH_USER/bucketmap?include_prereleases" alt="Release"></a>
  <a href="https://github.com/sponsors/YOUR_GH_USER"><img src="https://img.shields.io/badge/Sponsor-%E2%9D%A4-red" alt="Sponsor"></a>
</p>

<p align="center">
  <a href="https://bucketmap.dev">Website</a> ·
  <a href="https://docs.bucketmap.dev">Docs</a> ·
  <a href="https://discord.gg/bucketmap">Discord</a> ·
  <a href="https://github.com/YOUR_GH_USER/bucketmap/discussions">Discussions</a>
</p>

---

> **Status: early development.** The engine and the S3, GCS and Azure providers are being built toward v0.1. Interfaces may change until v1.0. Watch or star the repo to follow along.

## What it does

Misconfigured cloud storage is one of the most common ways data leaks. BucketMap takes a keyword (a company, product or project name), generates likely bucket names from it and a wordlist, checks each one across cloud providers in parallel, and tells you how exposed each bucket is.

| Access level | Meaning |
|---|---|
| `public read` | Objects can be read without credentials |
| `listable` | The bucket's contents can be listed without credentials |
| `private` | The bucket exists and denies anonymous access |

## Features

- **Fast:** a concurrent Go engine with configurable rate limits and timeouts.
- **Multi-cloud:** S3, GCS and Azure, behind one provider interface. Adding a provider is a single package.
- **Use it your way:** CLI, local web GUI, or import it as a Go library.
- **Scriptable:** clean plain-text and JSON output for pipelines.
- **Bring your own wordlists:** or use the built-in defaults.

## Install

```sh
go install github.com/YOUR_GH_USER/bucketmap/cmd/bucketmap@latest
```

Prebuilt binaries for Linux, macOS and Windows (amd64 and arm64) are on the [releases page](https://github.com/YOUR_GH_USER/bucketmap/releases).

## Quick start

```sh
bucketmap scan acme
```

Example output (illustrative):

```
scanning 1,240 candidate names on s3, gcs, azure
s3     acme-backups    PUBLIC READ
azure  acme-logs       LISTABLE
gcs    acme-assets     private
3 buckets found, 2 exposed
```

More in the [docs](https://docs.bucketmap.dev): [CLI](https://docs.bucketmap.dev/cli) · [Web GUI](https://docs.bucketmap.dev/web-gui) · [Go library](https://docs.bucketmap.dev/api) · [Wordlists](https://docs.bucketmap.dev/wordlists)

## Providers

| Provider | Status |
|---|---|
| Amazon S3 | In progress |
| Google Cloud Storage | In progress |
| Azure Blob Storage | In progress |
| DigitalOcean Spaces, Backblaze B2, Wasabi | Wanted: [add one](CONTRIBUTING.md) |

## Responsible use

Use BucketMap only on assets you own or are explicitly authorized to test. Finding an exposed bucket does not give you permission to download or modify its contents. If you find someone else's exposed data, report it to the owner through their disclosure program.

## Contributing

The easiest ways in are adding a cloud provider or a wordlist. Read [CONTRIBUTING.md](CONTRIBUTING.md) and open a "New provider" issue to get started. Questions go in [Discussions](https://github.com/YOUR_GH_USER/bucketmap/discussions).

## Sponsors

BucketMap is MIT licensed and funded by sponsors. [Become a sponsor](https://github.com/sponsors/YOUR_GH_USER) to get your name or logo here and on [bucketmap.dev](https://bucketmap.dev).

## Contributors

<a href="https://github.com/YOUR_GH_USER/bucketmap/graphs/contributors">
  <img src="https://contrib.rocks/image?repo=YOUR_GH_USER/bucketmap" alt="Contributors">
</a>

## License

[MIT](LICENSE)
__BM_EOF__
cat > 'SECURITY.md' <<'__BM_EOF__'
# Security Policy

Report vulnerabilities privately through GitHub: Security tab, then "Report a vulnerability". You can also email security@bucketmap.dev.

Please include steps to reproduce and the affected version. You will get a first reply within 5 days. Do not open public issues for vulnerabilities.
__BM_EOF__
mkdir -p 'docs/.vitepress'
cat > 'docs/.vitepress/config.mts' <<'__BM_EOF__'
import { defineConfig } from 'vitepress'

export default defineConfig({
  title: 'BucketMap',
  description: 'The Nmap for cloud storage',
  cleanUrls: true,
  themeConfig: {
    nav: [
      { text: 'Website', link: 'https://bucketmap.dev' },
      { text: 'GitHub', link: 'https://github.com/YOUR_GH_USER/bucketmap' },
      { text: 'Sponsor', link: 'https://github.com/sponsors/YOUR_GH_USER' },
    ],
    sidebar: [
      { text: 'Getting started', link: '/getting-started' },
      {
        text: 'Usage',
        items: [
          { text: 'CLI', link: '/cli' },
          { text: 'Web GUI', link: '/web-gui' },
          { text: 'Desktop app', link: '/desktop' },
        ],
      },
      { text: 'Go library', link: '/api' },
      { text: 'Wordlists', link: '/wordlists' },
      { text: 'Contributing', link: '/contributing' },
    ],
    search: { provider: 'local' },
    editLink: {
      pattern: 'https://github.com/YOUR_GH_USER/bucketmap/edit/main/docs/:path',
    },
    socialLinks: [{ icon: 'github', link: 'https://github.com/YOUR_GH_USER/bucketmap' }],
  },
})
__BM_EOF__
mkdir -p 'docs'
cat > 'docs/api.md' <<'__BM_EOF__'
# Go library

```go
package main

import (
	"context"
	"fmt"
	"log"

	"github.com/YOUR_GH_USER/bucketmap/pkg/bucketmap"
)

func main() {
	results, err := bucketmap.Scan(context.Background(), "acme", bucketmap.Options{})
	if err != nil {
		log.Fatalf("scanning: %v", err)
	}
	for _, r := range results {
		fmt.Println(r.Provider, r.Name, r.Access)
	}
}
```

::: warning
This shows the intended API shape. Update it to match the real `pkg/bucketmap` package once it exists.
:::
__BM_EOF__
mkdir -p 'docs'
cat > 'docs/cli.md' <<'__BM_EOF__'
# CLI usage

```sh
bucketmap scan <keyword>
```

BucketMap builds candidate bucket names from the keyword and a wordlist, checks each cloud provider, and prints the access level for every bucket it finds.

Run `bucketmap --help` and `bucketmap scan --help` for the flags in your version.

::: info
This page is a starter. Document each flag here as the CLI stabilizes.
:::
__BM_EOF__
mkdir -p 'docs'
cat > 'docs/contributing.md' <<'__BM_EOF__'
# Contributing

The full guide lives in [CONTRIBUTING.md](https://github.com/YOUR_GH_USER/bucketmap/blob/main/CONTRIBUTING.md): setup, adding a provider, adding a wordlist, and the pull request checklist.
__BM_EOF__
mkdir -p 'docs'
cat > 'docs/desktop.md' <<'__BM_EOF__'
# Desktop app

Download the installer for your platform from the [releases page](https://github.com/YOUR_GH_USER/bucketmap/releases). The desktop app wraps the web GUI and the scan engine in one install.
__BM_EOF__
mkdir -p 'docs'
cat > 'docs/getting-started.md' <<'__BM_EOF__'
# Getting started

Use BucketMap only on assets you own or are authorized to test.

## Install

::: code-group
```sh [Go (any OS)]
go install github.com/YOUR_GH_USER/bucketmap/cmd/bucketmap@latest
```
```sh [Linux / macOS]
# Download the archive for your platform from the latest release,
# extract it, and move the binary onto your PATH:
sudo mv bucketmap /usr/local/bin/
```
```powershell [Windows]
# Download bucketmap_windows_amd64.zip from the latest release,
# extract it, and add its folder to your PATH.
```
:::

Releases: <https://github.com/YOUR_GH_USER/bucketmap/releases>

## First scan

```sh
bucketmap scan acme
```

Next: [CLI usage](/cli).
__BM_EOF__
mkdir -p 'docs'
cat > 'docs/index.md' <<'__BM_EOF__'
---
layout: home
hero:
  name: BucketMap
  tagline: The Nmap for cloud storage. Find exposed buckets across S3, GCS and Azure.
  actions:
    - { theme: brand, text: Get started, link: /getting-started }
    - { theme: alt, text: GitHub, link: 'https://github.com/YOUR_GH_USER/bucketmap' }
---
__BM_EOF__
mkdir -p 'docs'
cat > 'docs/package.json' <<'__BM_EOF__'
{
  "name": "bucketmap-docs",
  "private": true,
  "type": "module",
  "scripts": {
    "dev": "vitepress dev",
    "build": "vitepress build",
    "preview": "vitepress preview"
  },
  "devDependencies": {
    "vitepress": "^1.5.0"
  }
}
__BM_EOF__
mkdir -p 'docs/public'
cat > 'docs/public/CNAME' <<'__BM_EOF__'
docs.bucketmap.dev
__BM_EOF__
mkdir -p 'docs'
cat > 'docs/web-gui.md' <<'__BM_EOF__'
# Web GUI

Start the local web interface and open it in your browser:

```sh
bucketmap serve
```

The interface runs the same engine as the CLI. Document the port, flags and screens here.
__BM_EOF__
mkdir -p 'docs'
cat > 'docs/wordlists.md' <<'__BM_EOF__'
# Wordlists

A wordlist is a plain text file, one word per line, lowercase, no duplicates. BucketMap combines each word with your keyword to build candidate bucket names.

## Use your own

```sh
bucketmap scan acme --wordlist ./my-words.txt
```

## Share one

Add it to `wordlists/` and open a pull request. See [Contributing](/contributing).
__BM_EOF__
mkdir -p 'website/app'
cat > 'website/app/globals.css' <<'__BM_EOF__'
@tailwind base;
@tailwind components;
@tailwind utilities;

body { @apply bg-ink text-text font-sans antialiased; }
a:focus-visible, button:focus-visible { outline: 2px solid #FFB020; outline-offset: 3px; }

/* One orchestrated moment: scan output prints line by line on load. */
.scan-line { opacity: 0; animation: print 0.01s steps(1) forwards; animation-delay: calc(var(--i) * 260ms); }
@keyframes print { to { opacity: 1; } }
@media (prefers-reduced-motion: reduce) { .scan-line { opacity: 1; animation: none; } }
__BM_EOF__
mkdir -p 'website/app'
cat > 'website/app/layout.tsx' <<'__BM_EOF__'
import type { Metadata } from "next";
import { IBM_Plex_Mono, IBM_Plex_Sans } from "next/font/google";
import "./globals.css";

const sans = IBM_Plex_Sans({ subsets: ["latin"], weight: ["400", "600"], variable: "--font-sans" });
const mono = IBM_Plex_Mono({ subsets: ["latin"], weight: ["400", "500"], variable: "--font-mono" });

export const metadata: Metadata = {
  title: "BucketMap: the Nmap for cloud storage",
  description: "Find exposed S3, GCS and Azure buckets. Open source, written in Go.",
};

export default function RootLayout({ children }: { children: React.ReactNode }) {
  return (
    <html lang="en" className={`${sans.variable} ${mono.variable}`}>
      <body>{children}</body>
    </html>
  );
}
__BM_EOF__
mkdir -p 'website/app'
cat > 'website/app/page.tsx' <<'__BM_EOF__'
import Terminal from "@/components/Terminal";

const REPO = "YOUR_GH_USER/bucketmap";
const GITHUB = `https://github.com/${REPO}`;

type Release = { assets?: { download_count: number }[] };

async function getStats() {
  try {
    const opts = { next: { revalidate: 3600 } };
    const [repo, releases] = await Promise.all([
      fetch(`https://api.github.com/repos/${REPO}`, opts).then((r) => r.json()),
      fetch(`https://api.github.com/repos/${REPO}/releases`, opts).then((r) => r.json()),
    ]);
    const downloads = Array.isArray(releases)
      ? (releases as Release[])
          .flatMap((r) => r.assets ?? [])
          .reduce((sum, a) => sum + a.download_count, 0)
      : 0;
    return { stars: (repo.stargazers_count as number) ?? 0, downloads };
  } catch {
    return { stars: 0, downloads: 0 };
  }
}

const features = [
  { title: "Fast", body: "A concurrent Go engine checks thousands of candidate names in seconds." },
  { title: "Multi-cloud", body: "S3, GCS and Azure today. New providers are a single package to add." },
  { title: "CLI, web and desktop", body: "Same engine everywhere, plus a Go library for your own tools." },
];

export default async function Home() {
  const { stars, downloads } = await getStats();

  return (
    <main className="mx-auto max-w-5xl px-6">
      <header className="flex items-center justify-between py-6 font-mono text-sm">
        <span className="font-medium">bucketmap</span>
        <nav className="flex gap-6 text-mute">
          <a className="hover:text-text" href="https://docs.bucketmap.dev">Docs</a>
          <a className="hover:text-text" href={GITHUB}>GitHub</a>
          <a className="hover:text-text" href={`https://github.com/sponsors/${REPO.split("/")[0]}`}>Sponsor</a>
        </nav>
      </header>

      <section className="py-16">
        <h1 className="max-w-3xl text-4xl font-semibold leading-tight sm:text-6xl">
          The Nmap for cloud storage.
        </h1>
        <p className="mt-5 max-w-xl text-lg text-mute">
          Find exposed S3, GCS and Azure buckets before someone else does. Open source, written in Go.
        </p>
        <code className="mt-8 block w-fit max-w-full select-all overflow-x-auto rounded-md border border-line bg-panel px-4 py-3 font-mono text-sm text-amber">
          go install github.com/{REPO}/cmd/bucketmap@latest
        </code>
        <div className="mt-10"><Terminal /></div>
      </section>

      <section className="grid gap-10 border-t border-line py-14 sm:grid-cols-3">
        {features.map((f) => (
          <div key={f.title} className="border-l-2 border-amber pl-4">
            <h2 className="font-semibold">{f.title}</h2>
            <p className="mt-2 text-mute">{f.body}</p>
          </div>
        ))}
      </section>

      <section className="flex flex-wrap gap-x-12 gap-y-6 border-t border-line py-14 font-mono">
        <div><div className="text-3xl">{stars.toLocaleString()}</div><div className="text-sm text-mute">GitHub stars</div></div>
        <div><div className="text-3xl">{downloads.toLocaleString()}</div><div className="text-sm text-mute">release downloads</div></div>
      </section>

      <section className="flex flex-wrap gap-4 border-t border-line py-14">
        <a href={GITHUB} className="rounded-md bg-amber px-5 py-3 font-semibold text-ink">View on GitHub</a>
        <a href={`${GITHUB}/releases`} className="rounded-md border border-line px-5 py-3">Download</a>
        <a href="https://docs.bucketmap.dev" className="rounded-md border border-line px-5 py-3">Read the docs</a>
      </section>

      <footer className="border-t border-line py-8 text-sm text-mute">
        MIT licensed. Use only on assets you are authorized to test.
      </footer>
    </main>
  );
}
__BM_EOF__
mkdir -p 'website/components'
cat > 'website/components/Terminal.tsx' <<'__BM_EOF__'
// Illustrative output. Replace with an asciinema recording once the CLI exists.
type Line = { text: string; tone?: "amber" | "coral" | "mint" | "mute" };

const lines: Line[] = [
  { text: "$ bucketmap scan acme", tone: "mute" },
  { text: "scanning 1,240 candidate names on s3, gcs, azure", tone: "mute" },
  { text: "s3     acme-backups    PUBLIC READ", tone: "coral" },
  { text: "azure  acme-logs       LISTABLE", tone: "amber" },
  { text: "gcs    acme-assets     private", tone: "mint" },
  { text: "3 buckets found, 2 exposed", tone: "mute" },
];

const toneClass = {
  amber: "text-amber",
  coral: "text-coral",
  mint: "text-mint",
  mute: "text-mute",
} as const;

export default function Terminal() {
  return (
    <pre
      aria-label="Example BucketMap scan output"
      className="overflow-x-auto rounded-md border border-line bg-panel p-5 font-mono text-sm leading-7"
    >
      {lines.map((l, i) => (
        <div
          key={l.text}
          className={`scan-line ${l.tone ? toneClass[l.tone] : ""}`}
          style={{ "--i": i } as React.CSSProperties}
        >
          {l.text}
        </div>
      ))}
    </pre>
  );
}
__BM_EOF__
mkdir -p 'website'
cat > 'website/next.config.mjs' <<'__BM_EOF__'
const nextConfig = {};
export default nextConfig;
__BM_EOF__
mkdir -p 'website'
cat > 'website/package.json' <<'__BM_EOF__'
{
  "name": "bucketmap-website",
  "private": true,
  "scripts": {
    "dev": "next dev",
    "build": "next build",
    "start": "next start"
  },
  "dependencies": {
    "next": "^14.2.0",
    "react": "^18.3.0",
    "react-dom": "^18.3.0"
  },
  "devDependencies": {
    "@types/node": "^20.0.0",
    "@types/react": "^18.3.0",
    "autoprefixer": "^10.4.0",
    "postcss": "^8.4.0",
    "tailwindcss": "^3.4.0",
    "typescript": "^5.4.0"
  }
}
__BM_EOF__
mkdir -p 'website'
cat > 'website/postcss.config.js' <<'__BM_EOF__'
module.exports = { plugins: { tailwindcss: {}, autoprefixer: {} } };
__BM_EOF__
mkdir -p 'website'
cat > 'website/tailwind.config.ts' <<'__BM_EOF__'
import type { Config } from "tailwindcss";

// Palette: deep ink-blue base, with amber for "open" findings and coral for "public".
const config: Config = {
  content: ["./app/**/*.{ts,tsx}", "./components/**/*.{ts,tsx}"],
  theme: {
    extend: {
      colors: {
        ink: "#0E1520",
        panel: "#151F2E",
        line: "#26354A",
        text: "#D5DEE8",
        mute: "#8A9BB0",
        amber: "#FFB020",
        coral: "#FF6B5E",
        mint: "#5FD3A5",
      },
      fontFamily: {
        sans: ["var(--font-sans)", "system-ui", "sans-serif"],
        mono: ["var(--font-mono)", "ui-monospace", "monospace"],
      },
    },
  },
  plugins: [],
};
export default config;
__BM_EOF__
mkdir -p 'website'
cat > 'website/tsconfig.json' <<'__BM_EOF__'
{
  "compilerOptions": {
    "target": "es2017",
    "lib": ["dom", "dom.iterable", "esnext"],
    "strict": true,
    "noEmit": true,
    "esModuleInterop": true,
    "module": "esnext",
    "moduleResolution": "bundler",
    "resolveJsonModule": true,
    "isolatedModules": true,
    "jsx": "preserve",
    "incremental": true,
    "plugins": [{ "name": "next" }],
    "paths": { "@/*": ["./*"] }
  },
  "include": ["next-env.d.ts", "**/*.ts", "**/*.tsx", ".next/types/**/*.ts"],
  "exclude": ["node_modules"]
}
__BM_EOF__

# ---- Go skeleton ----
mkdir -p cmd/bucketmap internal/engine wordlists
cat > cmd/bucketmap/main.go <<'__BM_EOF__'
package main

import (
	"fmt"
	"os"
)

// Set at build time by GoReleaser (see .goreleaser.yaml).
var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func main() {
	if len(os.Args) > 1 && os.Args[1] == "--version" {
		fmt.Printf("bucketmap %s (%s, %s)\n", version, commit, date)
		return
	}
	fmt.Fprintln(os.Stderr, "usage: bucketmap scan <keyword> (not implemented yet)")
	os.Exit(2)
}
__BM_EOF__

cat > internal/engine/provider.go <<'__BM_EOF__'
// Package engine holds the scan engine and the provider interface.
package engine

import "context"

// Access describes how exposed a bucket is.
type Access int

const (
	AccessUnknown Access = iota
	AccessPrivate
	AccessListable
	AccessPublicRead
)

func (a Access) String() string {
	switch a {
	case AccessPrivate:
		return "private"
	case AccessListable:
		return "listable"
	case AccessPublicRead:
		return "public read"
	default:
		return "unknown"
	}
}

// Result is one bucket found by a scan.
type Result struct {
	Provider string
	Name     string
	Access   Access
}

// Provider is implemented once per cloud storage service.
type Provider interface {
	Name() string
	// Candidates builds bucket names to try from a keyword and wordlist.
	Candidates(keyword string, words []string) []string
	// Check reports the access level of one bucket.
	Check(ctx context.Context, bucket string) (Access, error)
}
__BM_EOF__

printf 'dev\nprod\nstaging\nbackup\nbackups\nassets\nlogs\ndata\npublic\nprivate\nuploads\nmedia\n' > wordlists/default.txt

go mod init "github.com/$GH_USER/bucketmap"
go vet ./... && go build -o /dev/null ./cmd/bucketmap

# ---- Fill in your GitHub username ----
grep -rl 'YOUR_GH_USER' . --exclude-dir=.git | xargs sed -i "s/YOUR_GH_USER/$GH_USER/g"

# ---- Docs and website dependencies (also creates the lockfiles CI needs) ----
if [ "$INSTALL" = 1 ]; then
  (cd docs && npm install --no-audit --no-fund)
  (cd website && npm install --no-audit --no-fund)
fi

# ---- Git ----
git init -q -b main
git add -A
git commit -q -m "Initial scaffold: engine skeleton, docs, website, community files"

if [ "$PUSH" = 1 ]; then
  gh repo create "$GH_USER/bucketmap" --public --source=. --remote=origin --push \
    --description "The Nmap for cloud storage"
  gh repo edit "$GH_USER/bucketmap" --enable-discussions
fi

cat <<'__BM_EOF__'

Done. Next:
  cd bucketmap
  go run ./cmd/bucketmap --version
  (cd docs && npm run dev)       # docs on :5173
  (cd website && npm run dev)    # landing page on :3000
Still to add: logo.png at the repo root.
Test the release build locally: goreleaser release --snapshot --clean
__BM_EOF__

