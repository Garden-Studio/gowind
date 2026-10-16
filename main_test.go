package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestEveryPlatformMapsToAPublishedBinary(t *testing.T) {
	for _, platform := range []struct {
		goos, goarch string
		musl         bool
		want         string
	}{
		{"linux", "amd64", false, "tailwindcss-linux-x64"},
		{"linux", "arm64", true, "tailwindcss-linux-arm64-musl"},
		{"darwin", "arm64", false, "tailwindcss-macos-arm64"},
		{"darwin", "amd64", false, "tailwindcss-macos-x64"},
		{"windows", "amd64", false, "tailwindcss-windows-x64.exe"},
	} {
		got, err := assetName(platform.goos, platform.goarch, platform.musl)
		if err != nil || got != platform.want {
			t.Errorf("%s/%s musl=%v = %q, %v; want %q", platform.goos, platform.goarch, platform.musl, got, err, platform.want)
		}
		if _, ok := checksums[got]; !ok {
			t.Errorf("Tailwind %s publishes no checksum for %s", Version, got)
		}
	}
	if _, err := assetName("windows", "arm64", false); err == nil {
		t.Error("windows/arm64 has no binary")
	}
}

func TestADownloadThatDoesNotMatchItsChecksumIsRefused(t *testing.T) {
	path := filepath.Join(t.TempDir(), "tailwindcss")
	if err := os.WriteFile(path, []byte("not tailwind"), 0o755); err != nil {
		t.Fatal(err)
	}
	if verified(path, checksums["tailwindcss-linux-x64"]) {
		t.Fatal("a file with another hash passed")
	}
}
