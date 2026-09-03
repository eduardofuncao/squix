#!/usr/bin/env bash

set -e

# Configuration
APP_NAME="squix"
VERSION=${1:-$(git describe --tags --always --dirty 2>/dev/null || echo dev)}
BUILD_DIR="./dist"
NATIVE_OS=$(go env GOOS)
NATIVE_ARCH=$(go env GOARCH)

LITE_TAGS="noduckdb,nosnowflake,noracle"
MINIMAL_TAGS="noduckdb,nosnowflake,noracle,noclickhouse,nofirebird,nosqlserver"

echo "Building $APP_NAME version $VERSION"
echo "=================================="

# Clean previous builds
rm -rf "$BUILD_DIR"
mkdir -p "$BUILD_DIR"

FULL_DRIVERS=()

# Build with appropriate CGO setting based on target arch.
# Variants:
#   full    - every driver (DuckDB needs CGO; only the native platform gets CGO=1)
#   lite    - no DuckDB/Snowflake/Oracle, works without CGO
#   minimal - lite minus ClickHouse/Firebird/SQL Server (postgres, mysql, sqlite)
build_platform() {
    local platform=$1
    local goos=$2
    local goarch=$3
    local name_suffix=$4

    echo ""
    echo "Building $platform ($name_suffix)..."

    local tags=""
    if [ "$name_suffix" = "lite" ]; then tags=$LITE_TAGS; fi
    if [ "$name_suffix" = "minimal" ]; then tags=$MINIMAL_TAGS; fi

    local output_name="${APP_NAME}-${platform}"
    if [ "$name_suffix" != "full" ]; then
        output_name="${APP_NAME}-${name_suffix}-${platform}"
    fi
    if [ "$goos" = "windows" ]; then
        output_name="${output_name}.exe"
    fi

    local cgo=0
    if [ "$name_suffix" = "full" ] && [ "$goos" = "$NATIVE_OS" ] && [ "$goarch" = "$NATIVE_ARCH" ]; then
        cgo=1
    fi

    local extra_args=()
    if [ -n "$tags" ]; then
        extra_args=(-tags "$tags")
    fi

    CGO_ENABLED=$cgo GOOS=$goos GOARCH=$goarch go build \
        -ldflags="-s -w -X main.Version=$VERSION" \
        "${extra_args[@]}" \
        -o "$BUILD_DIR/$output_name" ./cmd/squix

    if [ "$cgo" = "1" ]; then
        FULL_DRIVERS+=("$platform")
    fi

    echo "✓ Built $output_name successfully"
}

for suffix in full lite minimal; do
    # Build Linux AMD64
    build_platform "linux-amd64" "linux" "amd64" "$suffix"

    # Build Linux ARM64
    build_platform "linux-arm64" "linux" "arm64" "$suffix"

    # Build Windows AMD64
    build_platform "windows-amd64" "windows" "amd64" "$suffix"

    # Build macOS AMD64
    build_platform "darwin-amd64" "darwin" "amd64" "$suffix"

    # Build macOS ARM64
    build_platform "darwin-arm64" "darwin" "arm64" "$suffix"
done

echo ""
echo "=================================="
echo "Build complete!"
echo ""
echo "Full driver support (CGO_ENABLED=1, native platform only):"
for p in "${FULL_DRIVERS[@]}"; do echo "  ✓ $p"; done
echo ""
echo "DuckDB is not included in non-native full builds (CGO_ENABLED=0):"
echo "  use the lite/minimal variants for those platforms instead"
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
ls -lh "$BUILD_DIR/"
echo ""
ls -lh "$BUILD_DIR/source/"
