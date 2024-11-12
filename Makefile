run:
	go run cmd/api/main.go
build:
	go build cmd/api/main.go
test:
	go test ./...

.PHONY: run go test
