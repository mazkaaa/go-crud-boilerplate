.PHONY: run build test fmt vet

run:
	go run ./cmd/server

build:
	go build -o bin/server ./cmd/server

test:
	go test ./... -v

fmt:
	go fmt ./...

vet:
	go vet ./...
