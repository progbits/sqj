# SQJson

Query JSON with SQL.

## Requirements

SQJson is build using the Go programming language. Instructions for downloading and installing Go can be found
at https://golang.org/dl/.

## Build and Installation

The project can be built and installed using the provided Makefile

```shell
make
make install
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

Extracting a single field

```shell
cat example.json | sqj 'SELECT id FROM [];' -

"6043c14205dfae1a521b819f"
```

Extracting multiple fields. Note that members clashing with SQL keywords must be quoted.

```shell
cat example.json | sqj 'SELECT "index", id, guid, about FROM [];' -

42
"6043c14205dfae1a521b819f"
"283bc66c-e5b3-4504-89c7-2df7e262cc49"
"Dolor id irure occaecat id do ea."
```

Queries can also contain arbitrary expressions.

```shell
cat example.json | sqj 'SELECT (5+4+3+2) % 6, (0 AND 1) AND NOT 0, guid about FROM [];'

2
0
"283bc66c-e5b3-4504-89c7-2df7e262cc49"
```
