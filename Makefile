# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test

BINARY_API=memo-api
BINARY_DAEMON=memo-daemon

all: clean test build 
build: 
	$(GOBUILD) -o $(BINARY_API) -v cmd/api/api.go 
	$(GOBUILD) -o $(BINARY_DAEMON) -v cmd/daemon/daemon.go 
test: 
	$(GOTEST) -v ./...
clean: 
	rm -f $(BINARY_API) 
	rm -f $(BINARY_DAEMON)
