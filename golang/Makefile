.PHONY: all build test fmt clean run

TARGET_NAME=sqj

all: test build
build:
	go build -tags=sqlite_vtable -o $(TARGET_NAME) -v ./...
test:
	go test -tags=sqlite_vtable -v ./...
fmt:
	go fmt ./...
clean:
	go clean
	rm -f $(TARGET_NAME)
	rm -f $(BINARY_UNIX)
run: build
	./$(TARGET_NAME)

