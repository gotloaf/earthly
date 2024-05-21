
$env:GOOS="js";
$env:GOARCH="wasm";

# Copy wasm_exec interface to the build directory
Copy-Item -Path "$(tinygo env TINYGOROOT)/targets/wasm_exec.js" -Destination build/wasm_exec.js
# Build the go project
tinygo build -o build/earthly.wasm -target wasm -no-debug ./cmd/wasm
# WASM optimization
# https://github.com/WebAssembly/binaryen
wasm-opt "build/earthly.wasm" -O2 --fast-math -o "build/earthly.wasm"
# https://github.com/WebAssembly/wabt
wasm-strip "build/earthly.wasm"
