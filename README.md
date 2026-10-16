# gotailwind

Run the [Tailwind CSS](https://tailwindcss.com) standalone CLI from a Go
project with `go tool`, pinned in `go.mod`, without Node.

```sh
go get -tool github.com/Garden-Studio/gotailwind/v4@v4.3.3
go tool gotailwind -i input.css -o output.css --minify
```

The module version is the Tailwind version it runs: `v4.3.3` runs Tailwind
CSS v4.3.3. The first run downloads that release's binary for your platform
from Tailwind's GitHub release, checks it against the SHA-256 the release
publishes (kept in `release.go`), and caches it under your user cache
directory (`GOTAILWIND_CACHE` overrides it). Every argument goes to the CLI
unchanged.

A scheduled workflow tags a new version within a day of each Tailwind release.

Platforms: Linux (glibc and musl) and macOS on amd64 and arm64, Windows on
amd64.
