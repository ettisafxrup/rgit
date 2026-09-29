#!/usr/bin/env sh
# Builds rgit on Linux and macOS. (The Windows installer is built with
# scripts/build.ps1 on Windows.)
#
#   ./scripts/build.sh            # test and build for this machine -> dist/rgit
#   ./scripts/build.sh --all      # build every platform with release names
set -eu

cd "$(dirname "$0")/.."
VERSION=${VERSION:-$(sed -n 's/^var Version = "\(.*\)"/\1/p' internal/version/version.go)}
LDFLAGS="-s -w -X github.com/ettisafxrup/rgit/internal/version.Version=$VERSION"

echo "Building rgit $VERSION"
go vet ./...
go test ./...
mkdir -p dist

build() { # build <GOOS> <GOARCH> <output>
    echo "> $1/$2 -> dist/$3"
    CGO_ENABLED=0 GOOS=$1 GOARCH=$2 go build -trimpath -ldflags "$LDFLAGS" -o "dist/$3" ./cmd/rgit
}

if [ "${1:-}" = "--all" ]; then
    build linux amd64 rgit-linux-x64
    build linux arm64 rgit-linux-arm64
    build darwin amd64 rgit-macos-x64
    build darwin arm64 rgit-macos-arm64
    build windows amd64 rgit-windows-x64.exe
    (cd dist && if command -v sha256sum >/dev/null; then sha256sum rgit-*; else shasum -a 256 rgit-*; fi > SHA256SUMS.txt)
    echo "Checksums written to dist/SHA256SUMS.txt"
else
    build "$(go env GOOS)" "$(go env GOARCH)" rgit
    echo "Install it with: sudo install dist/rgit /usr/local/bin/rgit"
fi
