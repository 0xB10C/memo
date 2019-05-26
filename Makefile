# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test


all: clean build test
build: 
				GO111MODULE=on $(GOBUILD) -o ./api/api -v ./api/api.go 
				GO111MODULE=on $(GOBUILD) -o ./memod/memod -v ./memod/memod.go 
test: 
				$(GOTEST) -v ./api/
				$(GOTEST) -v ./memod/
clean: 
				rm -f ./api/api
				rm -f ./memod/memod
