# Basic go commands
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get

DEVICE_ID=f015ca13-7e3b-431e-ac61-8ec593db09ae
OUTPUT_DIR=/Users/tom/Pictures/Backgrounds

# This is a change to trigger the build

# Binary names
BINARY_NAME=backdrop

all: test build
build: 
		$(GOBUILD) -o ./dist/$(BINARY_NAME) -v ./cmd/main.go
test: 
		$(GOTEST) -v ./...
clean: 
		$(GOCLEAN)
		rm -f ./dist/$(BINARY_NAME)
run:
		$(GOCLEAN)
		$(GOBUILD) -o ./dist/$(BINARY_NAME) -v ./cmd/main.go
		./dist/$(BINARY_NAME) 
build-prod:
		$(GOBUILD) -ldflags "-s -w" -o ./dist/$(BINARY_NAME) -v ./cmd/main.go 
docker-build:
		docker build -t backdrop-go .
docker-run:
		docker build -t backdrop-go .
		docker run -it --rm --name backdrop-go -v $(OUTPUT_DIR):/output -e BG_DEVICE_ID=$(DEVICE_ID) backdrop-go

