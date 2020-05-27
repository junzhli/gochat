# Go commands
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOPATH=$(shell $(GOCMD) env GOPATH)

# Program specific
BINARY_NAME_SERVER=gochatd
BINARY_NAME_CLIENT=gochat

all: test build
build:
	$(GOBUILD) -o ./dist/$(BINARY_NAME_SERVER) ./cmd/server
	$(GOBUILD) -o ./dist/$(BINARY_NAME_CLIENT) ./cmd/client
install:
	$(GOBUILD) -o $(GOPATH)/bin/$(BINARY_NAME_SERVER) ./cmd/server
	$(GOBUILD) -o $(GOPATH)/bin/$(BINARY_NAME_CLIENT) ./cmd/client
test:
	$(GOTEST) -v ./...
clean:
	$(GOCLEAN)
	rm -rf ./dist
run-client:
	$(GOCMD) run ./cmd/client/main.go
run-server:
	$(GOCMD) run ./cmd/server/main.go