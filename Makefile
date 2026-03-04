.PHONY: build build-all clean test install

BINARY = clappie
MODULE = github.com/brunojuliao/go-clappie

build:
	go build -o $(BINARY) .

install:
	go install .

test:
	go test ./...

clean:
	rm -f $(BINARY)
	rm -rf dist/

build-all: clean
	mkdir -p dist
	GOOS=linux   GOARCH=amd64 go build -o dist/$(BINARY)-linux-amd64 .
	GOOS=linux   GOARCH=arm64 go build -o dist/$(BINARY)-linux-arm64 .
	GOOS=darwin  GOARCH=amd64 go build -o dist/$(BINARY)-darwin-amd64 .
	GOOS=darwin  GOARCH=arm64 go build -o dist/$(BINARY)-darwin-arm64 .
	GOOS=windows GOARCH=amd64 go build -o dist/$(BINARY)-windows-amd64.exe .

lint:
	golangci-lint run ./...

fmt:
	go fmt ./...

vet:
	go vet ./...
