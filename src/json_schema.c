#include <stddef.h>

#include "json_parse.h"
#include "json_schema.h"
#include "util.h"

typedef struct Columns {
    char** columns;
    size_t n_columns;
} Columns;

// Add a new column to a JSON based table schema.
//
// If a prefix is specified, the column is added as 'prefix$column'.
void add_schema_column(JSONTableSchema* schema, char* prefix, char* column) {
    ++schema->n_columns;
    schema->columns =
        realloc(schema->columns, schema->n_columns * sizeof(char*));
    const size_t column_size =
        strlen(prefix) + strlen(column) + strlen("$") + 1;
    schema->columns[schema->n_columns - 1] =
        calloc(column_size, sizeof(char));
    strcpy(schema->columns[schema->n_columns - 1], prefix);
    if (strlen(prefix) > 0) {
        strcat(schema->columns[schema->n_columns - 1], "$");
    }
    strcat(schema->columns[schema->n_columns - 1], column);
}

// Extract table column names from a JSON AST.
void collect_columns(JSONNode* ast, JSONTableSchema* schema, char* prefix) {
    if (ast == NULL) {
        return;
    }

    switch (ast->value) {
        case (JSON_VALUE_OBJECT): {
            // Objects can either be top level nodes or themselves object
            // members. If an object has a *name* value, it is a member and
            // contributes its name to the prefix of its children.
            char* new_prefix = strdup(prefix);
            if (ast->name != NULL) {
                const size_t new_size =
                    strlen(new_prefix) + strlen(ast->name) + strlen("$") + 1;
                new_prefix = realloc(new_prefix, new_size);
                if (strlen(new_prefix) > 0) {
                    strcat(new_prefix, "$");
                }
                strcat(new_prefix, ast->name);

                // We also register the object itself as a column.
                add_schema_column(schema, prefix, ast->name);
            }

            for (int i = 0; i < ast->n_members; i++) {
                collect_columns(&ast->members[i], schema, new_prefix);
            }
            free(new_prefix);
            break;
        }
        case (JSON_VALUE_ARRAY):
        case (JSON_VALUE_NUMBER):
        case (JSON_VALUE_STRING):
        case (JSON_VALUE_NULL):
        case (JSON_VALUE_TRUE):
        case (JSON_VALUE_FALSE): {
            add_schema_column(schema, prefix, ast->name);
            break;
        }
        default: {
            log_and_exit("unknown value\n");
            break;
        }
    }
}

void extract_column_impl(JSONNode* ast, JSONNode** result, char* prefix,
                         const char* target) {
    if (ast == NULL) {
        return;
    }

    switch (ast->value) {
        case (JSON_VALUE_OBJECT): {
            // Objects can either be top level nodes or themselves object
            // members. If an object has a *name* value, it is a member and
            // contributes its name to the prefix of its children.
            char* new_prefix = strdup(prefix);
            if (ast->name != NULL) {
                const size_t new_size =
                    strlen(new_prefix) + strlen(ast->name) + strlen("$") + 1;
                new_prefix = realloc(new_prefix, new_size);
                if (strlen(new_prefix) > 0) {
                    strcat(new_prefix, "$");
                }
                strcat(new_prefix, ast->name);
            }

            if (strcmp(new_prefix, target) == 0) {
                *result = ast;
                break;
            }

            for (int i = 0; i < ast->n_members; i++) {
                extract_column_impl(&ast->members[i], result, new_prefix,
                                    target);
            }
            free(new_prefix);
            break;
        }
        case (JSON_VALUE_ARRAY):
        case (JSON_VALUE_NUMBER):
        case (JSON_VALUE_STRING):
        case (JSON_VALUE_NULL):
        case (JSON_VALUE_TRUE):
        case (JSON_VALUE_FALSE): {
            const size_t column_size =
                strlen(prefix) + strlen(ast->name) + strlen("$") + 1;
            char* column_name = calloc(column_size, sizeof(char));
            strcpy(column_name, prefix);
            if (strlen(prefix) > 0) {
                strcat(column_name, "$");
            }

            strcat(column_name, ast->name);
            if (strcmp(column_name, target) == 0) {
                *result = ast;
            }
            free(column_name);
            break;
        }
        default: {
            log_and_exit("unknown value\n");
            break;
        }
    }
}

void extract_column(JSONNode* ast, JSONNode** result, const char* target) {
    extract_column_impl(ast, result, "", target);
}

// Build a table schema from a JSON AST.
void build_table_schema(JSONNode* ast, JSONTableSchema** schema) {
    // Currently we only support JSON_TYPE_ARRAY as our top level node.
    if (ast->value != JSON_VALUE_ARRAY) {
        log_and_exit("not a JSON array\n");
    }

    // Some sanity checks.
    if (ast->n_values == 0) {
        log_and_exit("empty array\n");
    }

    // Assume that the first object in the array is representative and use
    // that object to determine the columns of the schema.
    *schema = calloc(1, sizeof(JSONTableSchema));
    collect_columns(&ast->values[0], *schema, "");

    // Build our 'CREATE TABLE ...' statement from our columns.
    char* mem_stream_data = NULL;
    size_t mem_stream_size = 0;
    FILE* mem_stream = open_memstream(&mem_stream_data, &mem_stream_size);
    fputs( "CREATE TABLE [] (", mem_stream);
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
