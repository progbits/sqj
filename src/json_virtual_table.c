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

int row_callback(void* pArg, int argc, char** argv, char** columnNames) {
    ClientData* client_data = (ClientData*)pArg;
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
    for (int i = 0; i < argc; i++) {
        ++result_object->n_members;
        result_object->members =
            realloc(result_object->members,
                    result_object->n_members * sizeof(struct JSONNode));

        JSONNode* member = NULL;
        extract_column(&client_data->ast->values[client_data->row], &member,
                       columnNames[i]);
        result_object->members[result_object->n_members - 1] = *member;
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
    json_vtab* vtab = (json_vtab*)sqlite3_malloc(sizeof(json_vtab));
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

int xDisconnect(sqlite3_vtab* pVTab) { return SQLITE_OK; }

int xDestroy(sqlite3_vtab* pVTab) { return SQLITE_OK; }

int xOpen(sqlite3_vtab* pVTab, sqlite3_vtab_cursor** ppCursor) {
    // Open a new cursor.
    json_vtab_cursor* cursor = sqlite3_malloc(sizeof(json_vtab_cursor));
    memset(cursor, 0, sizeof(json_vtab_cursor));
    cursor->client_data = ((json_vtab*)pVTab)->client_data;
    *ppCursor = (sqlite3_vtab_cursor*)cursor;
    return SQLITE_OK;
}

int xClose(sqlite3_vtab_cursor* pVtabCursor) { return SQLITE_OK; }

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
    return cursor->row >= cursor->client_data->ast->n_values;
}

int xColumn(sqlite3_vtab_cursor* pVtabCursor, sqlite3_context* pContext,
            int n) {
    json_vtab_cursor* cursor = (json_vtab_cursor*)pVtabCursor;

    // Get the value of the target column.
    char* target_column_name = cursor->client_data->schema->columns[n];
    JSONNode* ast_node;
    extract_column(&cursor->client_data->ast->values[cursor->row], &ast_node,
                   target_column_name);

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

    rc = sqlite3_exec(db, client_data->query, &row_callback, (void*)client_data,
                      NULL);
    if (rc != SQLITE_OK) {
        fprintf(stderr, "something went wrong\n");
        sqlite3_close(db);
        return rc;
    }

    sqlite3_close(db);
    return rc;
}
