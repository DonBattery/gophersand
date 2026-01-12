# Read version from VERSION file
$VERSION = (Get-Content -LiteralPath "VERSION" -Raw).Trim()

# Get short git commit hash
try {
    $BUILD = (git rev-parse --short HEAD).Trim()
}
catch {
    $BUILD = "dev"
}

# Set build environment
$Env:GOOS = "js"
$Env:GOARCH = "wasm"

# Build with ldflags (pass as a separate argument)
$ldflags = "-X main.VERSION=$VERSION -X main.BUILD=$BUILD"
Write-Host "ldflags = $ldflags"

go build -ldflags $ldflags -o "docs/gophersand.wasm" .

# Clean up environment
Remove-Item Env:GOOS, Env:GOARCH
