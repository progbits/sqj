#ifndef JSON_SCHEMA_H
#define JSON_SCHEMA_H

typedef struct JSONTableSchema {
    char* create_table_statement;
    char** columns;
    size_t n_columns;
} JSONTableSchema;

// Build a SQL 'CREATE TABLE ...' statement from a JSON AST.
void build_table_schema(JSONNode* ast, JSONTableSchema** table_schema);

void delete_table_schema(JSONTableSchema* schema);

// Extract the first matching AST node matching *target*.
void extract_column(JSONNode* ast, JSONNode** result, const char* target);

#endif // JSON_SCHEMA_H
