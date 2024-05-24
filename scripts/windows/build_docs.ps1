
# Install swag first:
# go install github.com/swaggo/swag/cmd/swag@latest

swag init -d "./" -g ./cmd/server/main.go -o ./build/docs
