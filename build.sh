#!/bin/bash

# build_for_linux.sh
# =====================
# Builds a Go project for Linux ARM (Raspberry Pi)

# Configuration
GO_FILE="main.go"                 # Name of your Go source file
OUTPUT_NAME="jellyfin_uploader"  # Output binary name
GOOS="linux"
GOARCH="arm64"
CGO_ENABLED=1

# Optional: Set GOARM if targeting 32-bit ARM (not needed for arm64)
# GOARM=7
# export GOARM

# Export environment variables
export GOOS=$GOOS
export GOARCH=$GOARCH
export CGO_ENABLED=$CGO_ENABLED

# Build the binary
echo "Building $GO_FILE for $GOOS/$GOARCH..."
go build -o bin/$OUTPUT_NAME $GO_FILE

if [ $? -eq 0 ]; then
    echo "✅ Build succeeded: $OUTPUT_NAME"
else
    echo "❌ Build failed."
    exit 1
fi

# Build frontend
cd webapp || exit 1
npm run build

if [ $? -eq 0 ]; then
    echo "✅ Successfully built dist"
else
    echo "❌ Web build failed."
    exit 1
fi
