run:
	go run cmd/api/*.go
build:
	go build cmd/api/*.go
test:
	go test ./...

.PHONY: run go test
