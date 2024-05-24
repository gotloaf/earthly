
$env:GOOS="windows";
$env:GOARCH="";

go build -o build/server.exe ./cmd/server
