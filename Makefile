.PHONY: sqjson test-tokenize test-parse test clean

.DEFAULT_GOAL := sqjson

sqjson: src/main.c src/json_tokenize.c src/json_parse.c src/json_schema.c src/util.c
	clang -g -O0 -lsqlite3 -o bin/sqjson src/main.c src/json_tokenize.c src/json_parse.c src/json_schema.c src/util.c

test-tokenize: test/test_json_tokenize.c src/json_tokenize.c
	clang -g test/test_json_tokenize.c src/json_tokenize.c -o bin/test-tokenize -O0 -lsqlite3 -lcheck -lm -lpthread -lrt -lsubunit

test-parse: test/test_json_parse.c src/json_parse.c
	clang -g test/test_json_parse.c src/json_parse.c src/util.c -o bin/test-parse -O0 -lsqlite3 -lcheck -lm -lpthread -lrt -lsubunit

test: test-tokenize test-parse
	./bin/test-tokenize && ./bin/test-parse

clean:
	rm -rf bin
