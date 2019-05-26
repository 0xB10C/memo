# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test


all: clean test build 
build: 
				cd api && $(GOBUILD) -o api -v api.go 
				cd memod && $(GOBUILD) -o memod -v memod.go 
test: 
				cd api && $(GOTEST) -v 
				cd memod && $(GOTEST) -v 
clean: 
				rm -f ./api/api
				rm -f ./memod/memod
