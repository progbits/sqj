ifeq ($(PREFIX),)
    PREFIX := /usr/local/bin
endif

BINARY = sqj

.PHONY: sqjson test-tokenize test-parse test-integration test install clean

.DEFAULT_GOAL := sqjson

sqjson: src/main.c src/json_tokenize.c src/json_parse.c src/json_schema.c src/util.c
	@mkdir -p bin
	clang -g -O3 -lsqlite3 -o bin/$(BINARY) src/main.c src/json_tokenize.c src/json_parse.c src/json_schema.c src/util.c

test-tokenize: test/test_json_tokenize.c src/json_tokenize.c
	clang -g test/test_json_tokenize.c src/json_tokenize.c -o bin/test-tokenize -O0 -lsqlite3 -lcheck -lm -lpthread -lrt -lsubunit

test-parse: test/test_json_parse.c src/json_parse.c
	clang -g test/test_json_parse.c src/json_parse.c src/util.c -o bin/test-parse -O0 -lsqlite3 -lcheck -lm -lpthread -lrt -lsubunit

test-integration: sqjson
	./test/integration.sh

test: test-tokenize test-parse test-integration
	./bin/test-tokenize && ./bin/test-parse

install:
	install -d $(DESTDIR)$(PREFIX)
	install -m 755 ./bin/sqj $(DESTDIR)$(PREFIX)

clean:
	rm -rf bin
