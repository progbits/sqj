#include <stddef.h>

#include "json_parse.h"
#include "json_schema.h"
#include "util.h"

typedef struct Columns {
    char** columns;
    size_t n_columns;
} Columns;

// Concatenate a prefix an a member name.
//
// JSON objects can either be top level nodes or themselves object members. For
// nested members, column names are a concatenation of the member names,
// separated by $. Top-level objects have no name, so their concatenated name is
// simply the prefix itself with no trailing $.
//
// The caller is responsible for freeing the returned string.
char* concat_prefix_name(const char* prefix, const char* name) {
    if (name == NULL) {
        return strdup(prefix);
    }

    // Allocate enough space for at most prefix + $ + name.
    const size_t length = strlen(prefix) + strlen(name) + 2;
    char* result = calloc(length, sizeof(char));
    if (strlen(prefix) > 0) {
        strcat(result, prefix);
        strcat(result, "$");
    }
    strcat(result, name);
    return result;
}

// Add a new column to JSONTableSchema.
void add_schema_column(JSONTableSchema* schema, char* column_name) {
    ++schema->n_columns;
    schema->columns =
        realloc(schema->columns, schema->n_columns * sizeof(char*));
    schema->columns[schema->n_columns - 1] = strdup(column_name);
}

// Extract table column names from a JSON AST.
void collect_columns(JSONNode* ast, JSONTableSchema* schema, char* prefix) {
    if (ast == NULL) {
        return;
    }

    char* column_name = concat_prefix_name(prefix, ast->name);
    if (ast->value == JSON_VALUE_OBJECT) {
        // Register named objects themselves as a column.
        if (ast->name != NULL) {
            add_schema_column(schema, column_name);
        }
        // Register the objects members as columns, prefixed with the
        // current column name.
        for (int i = 0; i < ast->n_members; i++) {
            collect_columns(&ast->members[i], schema, column_name);
        }
    } else {
        add_schema_column(schema, column_name);
    }
    free(column_name);
}

int extract_column_impl(JSONNode* ast, JSONNode** result, char* prefix,
                        const char* target) {
    if (ast == NULL) {
        return 0;
    }

    char* column_name = concat_prefix_name(prefix, ast->name);
    if (strcmp(column_name, target) == 0) {
        *result = ast;
        free(column_name);
        return 1;
    }

    if (ast->value == JSON_VALUE_OBJECT) {
        for (int i = 0; i < ast->n_members; i++) {
            JSONNode* member = &ast->members[i];
            if (extract_column_impl(member, result, column_name, target)) {
                free(column_name);
                return 1;
            }
        }
    }
    free(column_name);
    return 0;
}

void extract_column(JSONNode* ast, JSONNode** result, const char* target) {
    extract_column_impl(ast, result, "", target);
}

// Build a table schema from a JSON AST.
void build_table_schema(JSONNode* ast, JSONTableSchema** schema) {
    *schema = calloc(1, sizeof(JSONTableSchema));

    char* mem_stream_data = NULL;
    size_t mem_stream_size = 0;
    FILE* mem_stream = open_memstream(&mem_stream_data, &mem_stream_size);
    if (ast->value == JSON_VALUE_OBJECT) {
        collect_columns(ast, *schema, "");
        if (ast->n_members == 0) {
            add_schema_column(*schema, "INTERNAL_PLACEHOLDER");
        }
    } else if (ast->value == JSON_VALUE_ARRAY) {
        // Use the first entry of the array as the schema.
        // TODO: We should probably add a flag to relax this assumption and take
        //  the set of all values in the input.
        collect_columns(&ast->values[0], *schema, "");
        if (ast->n_values == 0) {
            add_schema_column(*schema, "INTERNAL_PLACEHOLDER");
        }
    }

    // Build our 'CREATE TABLE ...' statement from our columns.
    fputs("CREATE TABLE [] (", mem_stream);
    for (int i = 0; i < (*schema)->n_columns; i++) {
        fputs((*schema)->columns[i], mem_stream);
        if (i < (*schema)->n_columns - 1) {
            fputs(",", mem_stream);
        }
    }
    fputs(")", mem_stream);
    fflush(mem_stream);

    (*schema)->create_table_statement = strdup(mem_stream_data);

    const int rc = fclose(mem_stream);
    if (rc != 0) {
        log_and_exit("failed to close memstream\n");
    }
    free(mem_stream_data);
}

void delete_table_schema(JSONTableSchema* schema) {
    free(schema->create_table_statement);
    for (size_t i = 0; i < schema->n_columns; i++) {
        free(schema->columns[i]);
    }
    free(schema->columns);
    free(schema);
}
