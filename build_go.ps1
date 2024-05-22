
$env:GOOS="js";
$env:GOARCH="wasm";

# Copy wasm_exec interface to the build directory
Copy-Item -Path "$(go env GOROOT)/misc/wasm/wasm_exec.js" -Destination build/wasm_exec_go.js
# Build the go project
go build -o build/earthly_go.wasm -ldflags "-s -w" ./cmd/wasm
