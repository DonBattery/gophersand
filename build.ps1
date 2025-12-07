$Env:GOOS = 'js'
$Env:GOARCH = 'wasm'

go build -o docs/gophersand.wasm .

Remove-Item Env:GOOS
Remove-Item Env:GOARCH
