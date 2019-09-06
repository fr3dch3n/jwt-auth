MAKEFLAGS += --silent

all: test

test:
	go test -v ./...

coverage:
	go test -v ./...

fmt:
	gofmt -s -w .
