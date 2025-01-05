build:
	@go build -o bin/lucy cmd/lucy/*

run: build
	@./bin/lucy

test:
	@go test -v ./..
