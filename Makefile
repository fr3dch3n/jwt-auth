all: dep test

test:
	go test -v ./...

dep:
	dep ensure

fmt:
	gofmt -s -w .
