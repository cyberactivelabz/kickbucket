# Getting started

Use Kickbucket only on assets you own or are authorized to test.

## Install

::: code-group
```sh [Go (any OS)]
go install github.com/cyberactivelabz/kickbucket/cmd/kickbucket@latest
```
```sh [Linux / macOS]
# Download the archive for your platform from the latest release,
# extract it, and move the binary onto your PATH:
sudo mv kickbucket /usr/local/bin/
```
```powershell [Windows]
# Download kickbucket_windows_amd64.zip from the latest release,
# extract it, and add its folder to your PATH.
```
:::

Releases: <https://github.com/cyberactivelabz/kickbucket/releases>

## First scan

```sh
kickbucket scan acme
```

Next: [CLI usage](/cli).
