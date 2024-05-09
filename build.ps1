
$env:GOOS="js";
$env:GOARCH="wasm";

# Copy wasm_exec interface to the assets directory
Copy-Item -Path "$(go env GOROOT)/misc/wasm/wasm_exec.js" -Destination assets/wasm_exec.js
# Build the go project
go build -o assets/earthly.wasm -ldflags "-s -w" ./cmd/wasm
