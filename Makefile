TARGET_NAME=sqjson

.PHONY: all build clean test format

all: build

build:
	go build -v -o bin/$(TARGET_NAME) src/*.go

clean:
	rm -rf bin

test:
	go test -v

fmt:
	gofmt -w *.go
