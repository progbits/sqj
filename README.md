# sqj

[![Build](https://github.com/progbits/sqj/actions/workflows/build.yaml/badge.svg?branch=main)](https://github.com/progbits/sqj/actions/workflows/build.yaml)

Query JSON with SQL.

## Requirements

The sqj project is build using the Go programming language. Instructions for
downloading and installing Go can be found on the
[Go project website](https://golang.org/dl/).

## Build and Installation

The project can be built and installed using the provided Makefile:

```shell
make
```

Once build, the project can be installed:

```shell
cp bin/sql ~/.local/bin
```

## Running the Tests

Unit and integration tests can also be run using Make

```shell
make test
```

## Examples

### Querying objects

```json
{
  "id": "6043c14205dfae1a521b819f",
  "index": 42,
  "guid": "283bc66c-e5b3-4504-89c7-2df7e262cc49",
  "isActive": false,
  "about": "Dolor id irure occaecat id do ea."
}
```

Extracting a single field.

```shell
sqj 'SELECT id FROM [];' -

"6043c14205dfae1a521b819f"
```

Extracting multiple fields. Note, field names clashing with SQL keywords must be quoted.

```shell
sqj 'SELECT "index", id, guid, about FROM [];' -

42
"6043c14205dfae1a521b819f"
"283bc66c-e5b3-4504-89c7-2df7e262cc49"
"Dolor id irure occaecat id do ea."
```

Queries can also contain arbitrary expressions.

```shell
sqj 'SELECT (5+4+3+2) % 6, (0 AND 1) AND NOT 0, guid about FROM [];'

2
0
"283bc66c-e5b3-4504-89c7-2df7e262cc49"
```

Nested objects and arrays can also be extracted.

```json
{
  "id": "6043f3419f51a307278d160f",
  "index": 0,
  "guid": "2eb51437-51d9-458f-b805-877dcf2ef908",
  "isActive": false,
  "content": [
    {
      "id": 0,
      "word": "velit"
    },
    {
      "id": 1,
      "word": "culpa"
    },
    {
      "id": 2,
      "word": "pariatur"
    }
  ]
}
```

```shell
sqj 'SELECT content FROM [];' -

"[{"id": 0,"word": "velit"},{"id": 1,"word": "culpa"},{"id": 2,"word": "pariatur"}]"
```
