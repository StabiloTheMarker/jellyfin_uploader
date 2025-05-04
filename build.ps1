# build_for_linux.ps1
# =====================
# Builds a Go project for Linux ARM (Raspberry Pi)

# Configuration
$GoFile = "main.go"   # Name of your Go source file
$OutputName = "jellyfin_uploader"  # Name for the output binary
$GOOS = "linux"
$GOARCH = "arm64"

# Set environment variables for cross-compilation
$env:GOOS = $GOOS
$env:GOARCH = $GOARCH
$env:CGO_ENABLED=1
$env:GOARM = $GOARM

# Build the binary
Write-Host "Building $GoFile for $GOOS/$GOARCH (GOARM=$GOARM)..."
go build -o bin/$OutputName $GoFile

if ($LASTEXITCODE -eq 0) {
    Write-Host "✅ Build succeeded: $OutputName"
} else {
    Write-Host "❌ Build failed."
}

Set-Location webapp
npm run build

Write-Host "✅ Successfully built dist"

Set-Location ..
