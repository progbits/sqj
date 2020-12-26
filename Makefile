.PHONY: all build test clean

all: sqjson

sqjson: src/main.c src/vector.c src/json_tokenize.c
	clang -g -O0 -lsqlite3 -o bin/sqjson src/main.c src/vector.c src/json_tokenize.c

test: src/test_json_tokenize.c src/json_tokenize.c src/vector.c
	clang -g -O0 -lsqlite3 -o bin/test_sqjson src/test_json_tokenize.c src/json_tokenize.c src/vector.c

clean:
	rm -rf bin
