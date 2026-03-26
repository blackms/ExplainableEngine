.PHONY: build test lint run clean fmt

build:
	go build -o bin/server ./cmd/server

test:
	go test ./... -v -count=1

test-cover:
	go test ./... -coverprofile=coverage.out -covermode=atomic
	go tool cover -func=coverage.out

lint:
	go vet ./...
	@if command -v golangci-lint >/dev/null 2>&1; then golangci-lint run; fi

fmt:
	gofmt -w .
	goimports -w . 2>/dev/null || true

run:
	go run ./cmd/server

clean:
	rm -rf bin/ coverage.out *.db
