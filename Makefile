.PHONY: build run tidy test

build:
	go build -o _build/server ./cmd/server

run: build
	ACTIVE_ENV=DEV ./_build/server

tidy:
	go mod tidy

test:
	go test ./... -count=1
