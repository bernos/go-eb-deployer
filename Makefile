OUT_DIR=./bin
OUTPUT=$(OUT_DIR)/ebdeploy
PACKAGE=./cli/ebdeploy.go

ifeq ($(OS),Windows_NT)
	OUTPUT=$(OUT_DIR)/ebdeploy.exe
endif

all: test

build: test
	go build -v -o $(OUTPUT) $(PACKAGE)

test: deps
	go test -v ./...

deps:
	go get -v -t ./...

clean:
	go clean -i -x ./...
	-rm -rf $(OUT_DIR)

.PHONY: all build test deps clean
