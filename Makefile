.PHONY: lint, run

LEN ?= 20

lint:
	@golangci-lint run

run:
	@go run . $(LEN)
