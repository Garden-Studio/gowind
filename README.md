# gowind

[![Version](https://img.shields.io/github/v/tag/Garden-Studio/gowind?sort=semver&label=version)](https://github.com/Garden-Studio/gowind/tags)
[![Go Reference](https://pkg.go.dev/badge/github.com/Garden-Studio/gowind/v4.svg)](https://pkg.go.dev/github.com/Garden-Studio/gowind/v4)
[![Test](https://github.com/Garden-Studio/gowind/actions/workflows/test.yml/badge.svg)](https://github.com/Garden-Studio/gowind/actions/workflows/test.yml)
[![Follow Tailwind](https://github.com/Garden-Studio/gowind/actions/workflows/follow-tailwind.yml/badge.svg)](https://github.com/Garden-Studio/gowind/actions/workflows/follow-tailwind.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/Garden-Studio/gowind/v4)](https://goreportcard.com/report/github.com/Garden-Studio/gowind/v4)
[![Go version](https://img.shields.io/github/go-mod/go-version/Garden-Studio/gowind)](go.mod)
[![License](https://img.shields.io/github/license/Garden-Studio/gowind)](LICENSE)

Run the [Tailwind CSS](https://tailwindcss.com) standalone CLI from a Go
project with `go tool`, pinned in `go.mod`, without Node.

```sh
go get -tool github.com/Garden-Studio/gowind/v4@latest
go tool gowind -i input.css -o output.css --minify
```

The module version is the Tailwind version it runs: `v4.3.3` runs Tailwind
CSS v4.3.3. The first run downloads that release's binary for your platform
from Tailwind's GitHub release, checks it against the SHA-256 the release
publishes (kept in `release.go`), and caches it under your user cache
directory (`GOWIND_CACHE` overrides it). Every argument goes to the CLI
unchanged.

A scheduled workflow tags a new version within a day of each Tailwind
release; update with `go get -tool github.com/Garden-Studio/gowind/v4@latest`.

Platforms: Linux (glibc and musl) and macOS on amd64 and arm64, Windows on
amd64.
