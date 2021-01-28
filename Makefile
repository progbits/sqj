.PHONY: all build test fmt clean run

TARGET_NAME=sqj

all: test build
build:
	go build -tags=sqlite_vtable -o bin/$(TARGET_NAME) -v ./cmd/main
test:
	go test -tags=sqlite_vtable -v ./...
fmt:
	go fmt ./...
clean:
	go clean
	@rm -r bin
run: build
	./$(TARGET_NAME)
