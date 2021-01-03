# sqjson

Query JSON with SQL.

## Requirements

The project aims to keep external dependencies to a minimum.

Current dependencies are:
 - [Check Unit Testing Framework](https://libcheck.github.io/check/)
 - [Bats: Bash Automated Testing System](https://github.com/sstephenson/bats)

## Build and Installation

The project can be built and installed using Make

```shell
make
make install
```

## Running the Tests

Unit and integration tests can be run using Make

```shell
make test
```

## Limitations

Because names are registered as table columns, names cannot conflict with any
reserved [SQLite keywords](https://sqlite.org/lang_keywords.html). This
limitation will be addressed in a future release.

