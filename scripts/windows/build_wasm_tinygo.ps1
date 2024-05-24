
$env:GOOS="js";
$env:GOARCH="wasm";

# Copy wasm_exec interface to the build directory
Copy-Item -Path "$(tinygo env TINYGOROOT)/targets/wasm_exec.js" -Destination build/wasm_exec_tinygo.js
# Build the go project
tinygo build -o build/earthly_tinygo.wasm -target wasm -opt=2 -no-debug ./cmd/wasm
# WASM optimization
# https://github.com/WebAssembly/binaryen
wasm-opt "build/earthly_tinygo.wasm" -O2 --fast-math -o "build/earthly_tinygo.wasm"
# https://github.com/WebAssembly/wabt
wasm-strip "build/earthly_tinygo.wasm"
