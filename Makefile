.PHONY: all build clean

all: sqjson

sqjson: src/main.c
	clang -o bin/sqjson src/main.c

clean:
	rm -rf bin
