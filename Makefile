.PHONY: lint, build

lint:
	@golangci-lint run ./...

build:
	@go build -o temp cmd/pwdgen/main.go

