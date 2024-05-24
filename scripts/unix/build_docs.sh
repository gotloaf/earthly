# Install swag first:
# go install github.com/swaggo/swag/cmd/swag@latest

PATH=${PATH}:`go env GOPATH`/bin swag init -d "./" -g ./cmd/server/main.go -o ./build/docs
