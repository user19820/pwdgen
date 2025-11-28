.PHONY: lint, run, clean

LEN ?= 20

lint:
	@golangci-lint run

run:
	@go run . $(LEN)

clean:
	@rm -rf dist/
