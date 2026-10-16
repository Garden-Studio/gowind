// Command gowind runs the standalone Tailwind CSS CLI of the release this
// module version names, so a Go project pins Tailwind in go.mod and runs it
// with `go tool gowind` without Node.
//
// The first run downloads the release's binary for this platform from
// Tailwind's GitHub release, checks it against the SHA-256 the release
// published, and keeps it in the user cache; later runs reuse it. Every
// argument goes to the CLI unchanged.
package main

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"time"
)

func main() {
	binary, err := install()
	if err != nil {
		fmt.Fprintln(os.Stderr, "gowind:", err)
		os.Exit(1)
	}
	if err := run(binary, os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, "gowind:", err)
		os.Exit(1)
	}
}

// install returns the path of this release's binary, downloading and
// checking it first when the cache does not hold it yet.
func install() (string, error) {
	asset, err := assetName(runtime.GOOS, runtime.GOARCH, musl())
	if err != nil {
		return "", err
	}
	sum, ok := checksums[asset]
	if !ok {
		return "", fmt.Errorf("Tailwind %s publishes no %s", Version, asset)
	}
	dir, err := cacheDir()
	if err != nil {
		return "", err
	}
	binary := filepath.Join(dir, Version, asset)
	if verified(binary, sum) {
		return binary, nil
	}
	if err := download(binary, releaseURL(asset), sum); err != nil {
		return "", err
	}
	return binary, nil
}

// assetName is the release asset for a platform, as Tailwind names it.
func assetName(goos, goarch string, isMusl bool) (string, error) {
	arch := map[string]string{"amd64": "x64", "arm64": "arm64"}[goarch]
	if arch == "" {
		return "", fmt.Errorf("Tailwind publishes no binary for %s/%s", goos, goarch)
	}
	switch goos {
	case "linux":
		if isMusl {
			return "tailwindcss-linux-" + arch + "-musl", nil
		}
		return "tailwindcss-linux-" + arch, nil
	case "darwin":
		return "tailwindcss-macos-" + arch, nil
	case "windows":
		if arch == "x64" {
			return "tailwindcss-windows-x64.exe", nil
		}
	}
	return "", fmt.Errorf("Tailwind publishes no binary for %s/%s", goos, goarch)
}

// musl reports a Linux whose C library is musl (Alpine), which needs the
// -musl build.
func musl() bool {
	if runtime.GOOS != "linux" {
		return false
	}
	matches, _ := filepath.Glob("/lib/ld-musl-*")
	return len(matches) > 0
}

// cacheDir is GOWIND_CACHE, or gowind under the user cache.
func cacheDir() (string, error) {
	if dir := os.Getenv("GOWIND_CACHE"); dir != "" {
		return dir, nil
	}
	base, err := os.UserCacheDir()
	if err != nil {
		return "", fmt.Errorf("find a cache directory (set GOWIND_CACHE): %w", err)
	}
	return filepath.Join(base, "gowind"), nil
}

func releaseURL(asset string) string {
	return "https://github.com/tailwindlabs/tailwindcss/releases/download/" + Version + "/" + asset
}

// verified reports a file whose SHA-256 is sum.
func verified(path, sum string) bool {
	file, err := os.Open(path)
	if err != nil {
		return false
	}
	defer file.Close()
	hash := sha256.New()
	if _, err := io.Copy(hash, file); err != nil {
		return false
	}
	return hex.EncodeToString(hash.Sum(nil)) == sum
}

// download fetches url into path through a temporary file beside it, and
// keeps it only when its SHA-256 is sum.
func download(path, url, sum string) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	client := &http.Client{Timeout: 5 * time.Minute}
	response, err := client.Get(url)
	if err != nil {
		return fmt.Errorf("download %s: %w", url, err)
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("download %s: HTTP %d", url, response.StatusCode)
	}
	temp, err := os.CreateTemp(filepath.Dir(path), ".download-*")
	if err != nil {
		return err
	}
	defer os.Remove(temp.Name())
	hash := sha256.New()
	if _, err := io.Copy(io.MultiWriter(temp, hash), response.Body); err != nil {
		temp.Close()
		return fmt.Errorf("download %s: %w", url, err)
	}
	if err := temp.Close(); err != nil {
		return err
	}
	if got := hex.EncodeToString(hash.Sum(nil)); got != sum {
		return fmt.Errorf("download %s: SHA-256 %s, the release published %s", url, got, sum)
	}
	if err := os.Chmod(temp.Name(), 0o755); err != nil {
		return err
	}
	if err := os.Rename(temp.Name(), path); err != nil {
		return errors.Join(err, os.Remove(temp.Name()))
	}
	return nil
}
