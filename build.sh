#!/bin/bash

# build_for_linux.sh
# =====================
# Builds a Go project for Linux ARM64 (e.g., Raspberry Pi)
# Moves the binary and frontend dist folder to /jellyfin
# Sets ownership to jellyfin

# Configuration
GO_FILE="main.go"
OUTPUT_NAME="jellyfin_uploader"
GOOS="linux"
GOARCH="arm64"
CGO_ENABLED=1

DEST_DIR="/jellyfin"
BIN_DIR="$DEST_DIR/bin"
DIST_DIR="$DEST_DIR/dist"

# Export Go environment
export GOOS=$GOOS
export GOARCH=$GOARCH
export CGO_ENABLED=$CGO_ENABLED

# Build Go binary
echo "🔧 Building $GO_FILE for $GOOS/$GOARCH..."
go build -o bin/$OUTPUT_NAME .

if [ $? -ne 0 ]; then
    echo "❌ Go build failed."
    exit 1
fi

echo "✅ Go build succeeded."

# Build frontend
cd webapp || { echo "❌ Could not enter webapp directory"; exit 1; }
npm run build

if [ $? -ne 0 ]; then
    echo "❌ Web build failed."
    exit 1
fi

echo "✅ Frontend build succeeded."
cd ..

# Create destination directories
sudo mkdir -p "$BIN_DIR" "$DIST_DIR"

# Move binary and dist to /jellyfin
sudo cp "bin/$OUTPUT_NAME" "$BIN_DIR/"
sudo cp -r "webapp/dist/"* "$DIST_DIR/"

# Change ownership
sudo chown -R jellyfin:jellyfin "$DEST_DIR"

echo "🚚 Moved binary and dist to $DEST_DIR"
echo "✅ Ownership set to user: jellyfin"

