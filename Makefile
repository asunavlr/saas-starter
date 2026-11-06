.PHONY: run test lint fmt tidy

run:
	go run ./cmd/api

test:
	go test ./...

lint:
	go vet ./...

fmt:
	gofmt -w ./cmd ./internal

tidy:
	go mod tidy
