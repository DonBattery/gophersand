$ErrorActionPreference = "Stop"

# Read version from VERSION file
$VERSION = (Get-Content -LiteralPath "VERSION" -Raw).Trim()

# Build stamp: yyyy.MM.dd#HH:mm:ss#git-hash
$hash = (git rev-parse --short HEAD).Trim()
$timestamp = (Get-Date).ToString("yyyy.MM.dd#HH:mm:ss")
$BUILD = "$timestamp#$hash"

# Set build environment for WebAssembly
$Env:GOOS = "js"
$Env:GOARCH = "wasm"

# Set build flags
$ldflags = "-X main.VERSION=$VERSION -X main.BUILD=$BUILD"

# Build
go build -ldflags $ldflags -o "docs/gophersand.wasm" .

# Clean up environment
Remove-Item Env:GOOS
Remove-Item Env:GOARCH
