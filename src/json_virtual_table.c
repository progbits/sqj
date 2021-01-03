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
    if (client_data->columns_written == 0) {
        for (int i = 0; i < argc - 1; i++) {
            printf("%s,", columnNames[i]);
        }
        printf("%s\n", columnNames[argc - 1]);
        client_data->columns_written = 1;
    }

    for (int i = 0; i < argc - 1; i++) {
        JSONNode* result_node;
        printf("%s,", argv[i]);
    }
    printf("%s\n", argv[argc - 1]);
    return 0;
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
    (cursor->row)++;
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
            log_and_exit("unsupported column value\n");
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

int setup_virtual_table(sqlite3* db, ClientData* client_data) {
    // Register virtual table methods.
    sqlite3_module* json_vtab = calloc(1, sizeof(sqlite3_module));
    json_vtab->iVersion = 0;
    json_vtab->xCreate = &xCreate;
    json_vtab->xConnect = &xConnect;
    json_vtab->xBestIndex = &xBestIndex;
    json_vtab->xDisconnect = &xDisconnect;
    json_vtab->xDestroy = &xDestroy;
    json_vtab->xOpen = &xOpen;
    json_vtab->xClose = &xClose;
    json_vtab->xFilter = &xFilter;
    json_vtab->xNext = &xNext;
    json_vtab->xEof = &xEof;
    json_vtab->xColumn = &xColumn;
    json_vtab->xRowid = &xRowid;
    json_vtab->xUpdate = &xUpdate;
    json_vtab->xBegin = &xBegin;
    json_vtab->xSync = &xSync;
    json_vtab->xCommit = &xCommit;
    json_vtab->xRollback = &xRollback;
    json_vtab->xFindFunction = &xFindFunction;
    json_vtab->xRename = &xRename;
    json_vtab->xSavepoint = &xSavepoint;
    json_vtab->xRelease = &xRelease;
    json_vtab->xRollbackTo = &xRollbackTo;
    json_vtab->xShadowName = &xShadowName;

    int rc = SQLITE_OK;

    // Register the module.
    const char* module_name = "sqjson";
    rc = sqlite3_create_module(db, module_name, json_vtab, (void*)client_data);
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

int exec(sqlite3* db, ClientData* client_data) {
    return sqlite3_exec(db, client_data->query, &row_callback,
                        (void*)client_data, NULL);
}
