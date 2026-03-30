.PHONY: build run tidy test

build:
	go build -o _build/server ./cmd/server

run: build
	./_build/server configs/dev.yaml

tidy:
	go mod tidy

test:
	go test ./... -count=1
