.PHONY: all build test test-tokenize clean

all: sqjson

sqjson: src/main.c src/vector.c src/json_tokenize.c src/json_parse.c src/util.c
	clang -g -O0 -lsqlite3 -o bin/sqjson src/main.c src/vector.c src/json_tokenize.c src/json_parse.c src/util.c

test-tokenize: test/test_json_tokenize.c src/json_tokenize.c
	clang -g test/test_json_tokenize.c src/json_tokenize.c -o bin/test-tokenize -O0 -lsqlite3 -lcheck -lm -lpthread -lrt -lsubunit


test: test-tokenize
	./bin/test-tokenize

clean:
	rm -rf bin
