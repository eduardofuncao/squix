#!/usr/bin/env bash

set -e

# Configuration
APP_NAME="squix"
VERSION=${1:-"dev"}
BUILD_DIR="./dist"
NATIVE_OS=$(go env GOOS)
NATIVE_ARCH=$(go env GOARCH)

echo "Building $APP_NAME version $VERSION"
echo "=================================="

# Clean previous builds
rm -rf "$BUILD_DIR"
mkdir -p "$BUILD_DIR"

FULL_DRIVERS=()
LIMITED_DRIVERS=()

# Build with appropriate CGO setting based on target arch
build_platform() {
    local platform=$1
    local goos=$2
    local goarch=$3
    local output=$4

    echo ""
    echo "Building $platform..."

    if [ "$goos" = "$NATIVE_OS" ] && [ "$goarch" = "$NATIVE_ARCH" ]; then
        CGO_ENABLED=1 GOOS=$goos GOARCH=$goarch go build -ldflags='-s -w' -o "$output" ./cmd/squix
        FULL_DRIVERS+=("$platform")
    else
        CGO_ENABLED=0 GOOS=$goos GOARCH=$goarch go build -ldflags='-s -w' -o "$output" ./cmd/squix
        LIMITED_DRIVERS+=("$platform")
    fi

    echo "✓ Built $platform successfully"
}

# Build Linux AMD64
build_platform "Linux AMD64" "linux" "amd64" "$BUILD_DIR/${APP_NAME}-linux-amd64"

# Build Linux ARM64
build_platform "Linux ARM64" "linux" "arm64" "$BUILD_DIR/${APP_NAME}-linux-arm64"

# Build Windows AMD64
build_platform "Windows AMD64" "windows" "amd64" "$BUILD_DIR/${APP_NAME}-windows-amd64.exe"

# Build macOS AMD64
build_platform "macOS AMD64" "darwin" "amd64" "$BUILD_DIR/${APP_NAME}-darwin-amd64"

# Build macOS ARM64
build_platform "macOS ARM64" "darwin" "arm64" "$BUILD_DIR/${APP_NAME}-darwin-arm64"

echo ""
echo "=================================="
echo "Build complete!"
echo "=================================="
echo ""
echo "Full driver support (CGO_ENABLED=1):"
for p in "${FULL_DRIVERS[@]}"; do echo "  ✓ $p"; done
echo ""
echo "Limited driver support - no DuckDB (CGO_ENABLED=0):"
for p in "${LIMITED_DRIVERS[@]}"; do echo "  ✓ $p"; done
echo ""

# Create source code archives
echo "Creating source code archives..."
mkdir -p "$BUILD_DIR/source"
git archive --format=tar.gz --prefix="${APP_NAME}-${VERSION}/" -o "$BUILD_DIR/source/${APP_NAME}_${VERSION}_source.tar.gz" HEAD
git archive --format=zip --prefix="${APP_NAME}-${VERSION}/" -o "$BUILD_DIR/source/${APP_NAME}_${VERSION}_source.zip" HEAD
echo "✓ Source archives created"

echo ""
echo "Generating checksums for all artifacts..."
find "$BUILD_DIR" -type f -exec sha256sum {} + | sed "s|$BUILD_DIR/||" > "$BUILD_DIR/checksums.txt"
echo "✓ Checksums generated in $BUILD_DIR/checksums.txt"

echo ""
echo "=================================="
echo "All tasks complete!"
echo "=================================="
echo ""
echo "Binaries location: $BUILD_DIR/"
echo "Source archives location: $BUILD_DIR/source/"
echo "Checksum file: $BUILD_DIR/checksums.txt"
echo ""
ls -lh "$BUILD_DIR/"
echo ""
ls -lh "$BUILD_DIR/source/"
echo ""
