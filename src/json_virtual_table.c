#include <stddef.h>

#include "json_virtual_table.h"
#include "util.h"

typedef struct json_vtab {
    sqlite3_vtab base;
    sqlite3* db;
    ClientData* client_data;
} json_vtab;

typedef struct json_vtab_cursor {
    sqlite3_vtab_cursor* base;
    ClientData* client_data;
    int row;
} json_vtab_cursor;

void debug_print(sqlite3_index_info* info) {
    printf("=== debug_print sqlite3_index_info ===\n");
    for (int i = 0; i < info->nConstraint; i++) {
        printf("Column: %d\n"
               "Operator: %d\n",
               info->aConstraint[i].iColumn, info->aConstraint[i].op);
    }
    printf("=== debug_print sqlite3_index_info ===\n");
}

int row_callback(ClientData* client_data) {
    int column_count = sqlite3_column_count(client_data->stmt);

    // We should probably handle this allocation when we construct the
    // ClientData instance.
    if (!client_data->result_ast) {
        client_data->result_ast = calloc(1, sizeof(JSONNode));
        client_data->result_ast->value = JSON_VALUE_ARRAY;
    }

    // Create a new object to hold the current record.
    ++client_data->result_ast->n_values;
    client_data->result_ast->values =
        realloc(client_data->result_ast->values,
                client_data->result_ast->n_values * sizeof(JSONNode));
    memset(
        &client_data->result_ast->values[client_data->result_ast->n_values - 1],
        0, sizeof(JSONNode));

    JSONNode* result_object =
        &client_data->result_ast->values[client_data->result_ast->n_values - 1];
    result_object->value = JSON_VALUE_OBJECT;
    for (int i = 0; i < column_count; i++) {
        ++result_object->n_members;
        result_object->members =
            realloc(result_object->members,
                    result_object->n_members * sizeof(struct JSONNode));

        const char* column_name = sqlite3_column_name(client_data->stmt, i);

        JSONNode* source_node = NULL;
        if (client_data->ast->value == JSON_VALUE_OBJECT) {
            extract_column(client_data->ast, &source_node, column_name);
        } else if (client_data->ast->value == JSON_VALUE_ARRAY) {
            extract_column(&client_data->ast->values[client_data->row],
                           &source_node, column_name);
        }

        if (!source_node) {
            // At the moment, we assume that this is the result of an aliased
            // expression. However, it could be a column that hasn't been
            // registered in our schema. We should use sqlite3_column_type(...)
            // to get the actual type of the column.
            const int column_type = sqlite3_column_type(client_data->stmt, i);
            switch (column_type) {
                case SQLITE_INTEGER:
                case SQLITE_FLOAT: {
                    result_object->members[result_object->n_members - 1].value =
                        JSON_VALUE_NUMBER;
                    result_object->members[result_object->n_members - 1]
                        .number_value =
                        sqlite3_column_double(client_data->stmt, i);
                    break;
                }
                case SQLITE_TEXT: {
                    result_object->members[result_object->n_members - 1].value =
                        JSON_VALUE_STRING;
                    result_object->members[result_object->n_members - 1].name =
                        strdup(column_name);
                    result_object->members[result_object->n_members - 1]
                        .string_value =
                        strdup(sqlite3_column_text(client_data->stmt, i));
                    break;
                }
                case SQLITE_BLOB: {
                    log_and_exit("handle values of type SQLITE_BLOB\n");
                }
                case SQLITE_NULL: {
                    result_object->members[result_object->n_members - 1].value =
                        JSON_VALUE_NULL;
                    break;
                }
                default: {
                    return SQLITE_FAIL;
                }
            }
        } else {
            deep_clone(source_node,
                       &result_object->members[result_object->n_members - 1]);
        }
    }

    return SQLITE_OK;
}

int xCreate(sqlite3* db, void* pAux, int argc, const char* const* argv,
            sqlite3_vtab** ppVTab, char** pzErr) {
    // Declare the schema of our virtual table.
    ClientData* client_data = pAux;
    int rc =
        sqlite3_declare_vtab(db, client_data->schema->create_table_statement);
    char** message;
    if (rc != 0) {
        fprintf(stderr, "%s!\n", sqlite3_errmsg(db));
        sqlite3_close(db);
        return 1;
    }

    // Allocate a new json_vtab instance.
    json_vtab* vtab = calloc(1, sizeof(json_vtab));
    vtab->client_data = client_data;
    *ppVTab = &vtab->base;

    return SQLITE_OK;
}

int xConnect(sqlite3* db, void* pAux, int argc, const char* const* argv,
             sqlite3_vtab** ppVTab, char** pzErr) {
    // xConnect is just xCreate for ephemeral virtual tables.
    return xCreate(db, pAux, argc, argv, ppVTab, pzErr);
}

int xBestIndex(sqlite3_vtab* pVTab, sqlite3_index_info* pIndexInfo) {
    return SQLITE_OK;
}

int xDisconnect(sqlite3_vtab* pVTab) {
    free(pVTab);
    return SQLITE_OK;
}

int xDestroy(sqlite3_vtab* pVTab) {
    free(pVTab);
    return SQLITE_OK;
}

int xOpen(sqlite3_vtab* pVTab, sqlite3_vtab_cursor** ppCursor) {
    // Open a new cursor.
    json_vtab_cursor* cursor = calloc(1, sizeof(json_vtab_cursor));
    memset(cursor, 0, sizeof(json_vtab_cursor));
    cursor->client_data = ((json_vtab*)pVTab)->client_data;
    *ppCursor = (sqlite3_vtab_cursor*)cursor;
    return SQLITE_OK;
}

int xClose(sqlite3_vtab_cursor* pVtabCursor) {
    free(pVtabCursor);
    return SQLITE_OK;
}

int xFilter(sqlite3_vtab_cursor* pVtabCursor, int idxNum, const char* idxStr,
            int argc, sqlite3_value** argv) {
    return SQLITE_OK;
}

int xNext(sqlite3_vtab_cursor* pVtabCursor) {
    json_vtab_cursor* cursor = (json_vtab_cursor*)pVtabCursor;
    ++(cursor->row);
    ++(cursor->client_data->row);
    return SQLITE_OK;
}

int xEof(sqlite3_vtab_cursor* pVtabCursor) {
    json_vtab_cursor* cursor = (json_vtab_cursor*)pVtabCursor;

    if (cursor->client_data->ast->value == JSON_VALUE_OBJECT) {
        return cursor->row > 0;
    }
    return cursor->row >= cursor->client_data->ast->n_values;
}

int xColumn(sqlite3_vtab_cursor* pVtabCursor, sqlite3_context* pContext,
            int n) {
    json_vtab_cursor* cursor = (json_vtab_cursor*)pVtabCursor;

    // Get the value of the target column.
    char* target_column_name = cursor->client_data->schema->columns[n];
    JSONNode* ast_node = NULL;
    if (cursor->client_data->ast->value == JSON_VALUE_OBJECT) {
        extract_column(cursor->client_data->ast, &ast_node, target_column_name);
    } else if (cursor->client_data->ast->value == JSON_VALUE_ARRAY) {
        extract_column(&cursor->client_data->ast->values[cursor->row],
                       &ast_node, target_column_name);
    }

    // Record the value.
    switch (ast_node->value) {
        case JSON_VALUE_OBJECT:
        case JSON_VALUE_ARRAY: {
            // Stringify object or array.
            char* buffer;
            size_t buffer_size;
            FILE* memstream = open_memstream(&buffer, &buffer_size);
            pretty_print(ast_node, memstream, 0);
            fclose(memstream);

            // Copy the result text to an SQLite owned block and free our copy.
            sqlite3_result_text(pContext, buffer, -1, SQLITE_TRANSIENT);
            free(buffer);
            break;
        }
        case JSON_VALUE_NUMBER: {
            sqlite3_result_double(pContext, ast_node->number_value);
            break;
        }
        case JSON_VALUE_STRING: {
            sqlite3_result_text(pContext, ast_node->string_value, -1, NULL);
            break;
        }
        case JSON_VALUE_NULL: {
            sqlite3_result_null(pContext);
            break;
        }
        case JSON_VALUE_TRUE: {
            sqlite3_result_int(pContext, 1);
            break;
        }
        case JSON_VALUE_FALSE: {
            sqlite3_result_int(pContext, 0);
            break;
        }
    }
    return SQLITE_OK;
}

int xRowid(sqlite3_vtab_cursor* pVtabCursor, sqlite3_int64* pRowid) {
    return SQLITE_OK;
}

int xUpdate(sqlite3_vtab* pVtabCursor, int argc, sqlite3_value** argv,
            sqlite3_int64* pRowid) {
    return SQLITE_OK;
}

int xBegin(sqlite3_vtab* pVTab) { return SQLITE_OK; }

int xSync(sqlite3_vtab* pVTab) { return SQLITE_OK; }

int xCommit(sqlite3_vtab* pVTab) { return SQLITE_OK; }

int xRollback(sqlite3_vtab* pVTab) { return SQLITE_OK; }

int xFindFunction(sqlite3_vtab* pVtab, int nArg, const char* zName,
                  void (**pxFunc)(sqlite3_context*, int, sqlite3_value**),
                  void** ppArg) {
    return SQLITE_OK;
}

int xRename(sqlite3_vtab* pVtab, const char* zNew) { return SQLITE_OK; }

int xSavepoint(sqlite3_vtab* pVTab, int n) { return SQLITE_OK; }

int xRelease(sqlite3_vtab* pVTab, int n) { return SQLITE_OK; }

int xRollbackTo(sqlite3_vtab* pVTab, int n) { return SQLITE_OK; }

int xShadowName(const char* name) { return SQLITE_OK; }

static sqlite3_module module = {
    0, // iVersion
    xCreate,
    xConnect,
    xBestIndex,
    xDisconnect,
    xDestroy,
    xOpen,
    xClose,
    xFilter,
    xNext,
    xEof,
    xColumn,
    xRowid,
    NULL,          // xUpdate
    NULL,          // xBegin
    NULL,          // xSync
    NULL,          // xCommit
    NULL,          // xRollback
    xFindFunction, // xFindFunction
    NULL,          // xRename
    NULL,          // xSavepoint
    NULL,          // xRelease
    NULL           // xRollbackto
};

int setup_virtual_table(sqlite3* db, ClientData* client_data) {
    int rc = SQLITE_OK;

    // Register the module.
    const char* module_name = "sqjson";
    rc = sqlite3_create_module(db, module_name, &module, (void*)client_data);
    if (rc != SQLITE_OK) {
        return rc;
    }

    // Create the virtual table.
    const char* statement = "CREATE VIRTUAL TABLE [] USING sqjson";
    rc = sqlite3_exec(db, statement, NULL, NULL, NULL);
    if (rc != SQLITE_OK) {
        return rc;
    }

    return SQLITE_OK;
}

int exec(ClientData* client_data) {
    int rc = SQLITE_OK;

    sqlite3* db = NULL;
    rc = sqlite3_open(":memory:", &db);
    if (rc != SQLITE_OK) {
        fprintf(stderr, "failed to open in-memory database\n");
        sqlite3_close(db);
        return rc;
    }

    rc = setup_virtual_table(db, client_data);
    if (rc != SQLITE_OK) {
        fprintf(stderr, "something went wrong\n");
        sqlite3_close(db);
        return rc;
    }

    rc = sqlite3_prepare_v2(db, client_data->query, strlen(client_data->query),
                            &client_data->stmt, NULL);
    if (rc != SQLITE_OK) {
        fprintf(stderr, "something went wrong\n");
        sqlite3_close(db);
        return rc;
    }

    // If we have an empty array or an empty object and our query doesn't
    // contain any columns, except our internal placeholder required to have a
    // valid CREATE TABLE statement for empty objects/arrays. We can just pretty
    // print the AST as it stands.
    const int n_columns = sqlite3_column_count(client_data->stmt);
    if (n_columns == 1) {
        const char* column_name = sqlite3_column_name(client_data->stmt, 0);
        if (strcmp(column_name, "INTERNAL_PLACEHOLDER") == 0) {
            client_data->result_ast = calloc(1, sizeof(JSONNode));
            client_data->result_ast->value = client_data->ast->value;
            sqlite3_finalize(client_data->stmt);
            sqlite3_close(db);
            return rc;
        }
    }

    while (sqlite3_step(client_data->stmt) == SQLITE_ROW) {
        rc = row_callback(client_data);
        if (rc != SQLITE_OK) {
            break;
        }
    }

    sqlite3_finalize(client_data->stmt);
    sqlite3_close(db);
    return rc;
}
