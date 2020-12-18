.PHONY: all build clean

all: sqjson

sqjson: src/main.c
	clang -g -O0 -lsqlite3 -o bin/sqjson src/main.c

clean:
	rm -rf bin
