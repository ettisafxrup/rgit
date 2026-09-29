#!/usr/bin/env sh
# Installs the latest rgit release on Linux or macOS:
#
#   curl -fsSL https://raw.githubusercontent.com/ettisafxrup/rgit/main/scripts/install.sh | sh
#
# Settings (environment variables):
#   RGIT_INSTALL_DIR  where to put rgit (default: /usr/local/bin if writable,
#                     otherwise ~/.local/bin)
#   RGIT_BASE_URL     where release files are downloaded from
set -eu

BASE_URL=${RGIT_BASE_URL:-https://github.com/ettisafxrup/rgit/releases/latest/download}

fail() { printf 'rgit install: %s\n' "$1" >&2; exit 1; }

case "$(uname -s)" in
    Linux)  os=linux ;;
    Darwin) os=macos ;;
    *)      fail "unsupported system $(uname -s); on Windows use the installer from the releases page" ;;
esac
case "$(uname -m)" in
    x86_64 | amd64)  arch=x64 ;;
    arm64 | aarch64) arch=arm64 ;;
    *)               fail "unsupported processor $(uname -m)" ;;
esac
asset="rgit-$os-$arch"

if [ -n "${RGIT_INSTALL_DIR:-}" ]; then
    dir=$RGIT_INSTALL_DIR
elif [ -w /usr/local/bin ]; then
    dir=/usr/local/bin
else
    dir="$HOME/.local/bin"
fi
mkdir -p "$dir"

tmp=$(mktemp -d)
trap 'rm -rf "$tmp"' EXIT

download() { # download <url> <file>
    if command -v curl >/dev/null 2>&1; then
        curl -fsSL "$1" -o "$2"
    elif command -v wget >/dev/null 2>&1; then
        wget -qO "$2" "$1"
    else
        fail "curl or wget is required"
    fi
}

echo "Downloading $asset ..."
download "$BASE_URL/$asset" "$tmp/$asset" || fail "download failed: $BASE_URL/$asset"

# Verify the checksum when the release publishes one.
if download "$BASE_URL/SHA256SUMS.txt" "$tmp/SHA256SUMS.txt" 2>/dev/null; then
    expected=$(grep " $asset\$" "$tmp/SHA256SUMS.txt" | cut -d' ' -f1)
    if command -v sha256sum >/dev/null 2>&1; then
        actual=$(sha256sum "$tmp/$asset" | cut -d' ' -f1)
    else
        actual=$(shasum -a 256 "$tmp/$asset" | cut -d' ' -f1)
    fi
    [ -n "$expected" ] && [ "$expected" = "$actual" ] || fail "checksum mismatch for $asset"
    echo "Checksum verified."
fi

chmod +x "$tmp/$asset"
# macOS marks downloaded files as quarantined; the binary is not notarized.
[ "$os" = macos ] && xattr -d com.apple.quarantine "$tmp/$asset" 2>/dev/null || true
mv "$tmp/$asset" "$dir/rgit"

echo "Installed rgit to $dir/rgit"
case ":$PATH:" in
    *":$dir:"*) "$dir/rgit" version ;;
    *) echo "Add $dir to your PATH, e.g.:  echo 'export PATH=\"$dir:\$PATH\"' >> ~/.profile" ;;
esac
