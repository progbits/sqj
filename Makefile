.PHONY: all build clean

all: sqjson

sqjson: src/main.c src/vector.c
	clang -g -O0 -lsqlite3 -o bin/sqjson src/main.c src/vector.c

clean:
	rm -rf bin
