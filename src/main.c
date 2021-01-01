#include <sqlite3.h>
#include <stdio.h>
#include <string.h>

#include "json_parse.h"
#include "json_schema.h"
#include "json_tokenize.h"
#include "util.h"

void debug_print(sqlite3_index_info* info) {
    printf("=== debug_print sqlite3_index_info ===\n");
    for (int i = 0; i < info->nConstraint; i++) {
        printf("Column: %d\n"
               "Operator: %d\n",
               info->aConstraint[i].iColumn, info->aConstraint[i].op);
    }
    printf("=== debug_print sqlite3_index_info ===\n");
}

typedef struct ClientData {
    JSONNode* ast;
    JSONTableSchema* schema;
    char* query;
    int columns_written;
} ClientData;

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
    return 0;
}

int xDisconnect(sqlite3_vtab* pVTab) { return 0; }

int xDestroy(sqlite3_vtab* pVTab) { return 0; }

int xOpen(sqlite3_vtab* pVTab, sqlite3_vtab_cursor** ppCursor) {
    // Open a new cursor.
    json_vtab_cursor* cursor = sqlite3_malloc(sizeof(json_vtab_cursor));
    memset(cursor, 0, sizeof(json_vtab_cursor));
    cursor->client_data = ((json_vtab*)pVTab)->client_data;
    *ppCursor = (sqlite3_vtab_cursor*)cursor;
    return SQLITE_OK;
}

int xClose(sqlite3_vtab_cursor* pVtabCursor) { return 0; }

int xFilter(sqlite3_vtab_cursor* pVtabCursor, int idxNum, const char* idxStr,
            int argc, sqlite3_value** argv) {
    return 0;
}

int xNext(sqlite3_vtab_cursor* pVtabCursor) {
    json_vtab_cursor* cursor = (json_vtab_cursor*)pVtabCursor;
    (cursor->row)++;
    return SQLITE_OK;
}

int xEof(sqlite3_vtab_cursor* pVtabCursor) {
    json_vtab_cursor* cursor = (json_vtab_cursor*)pVtabCursor;
    const int done = cursor->row >= cursor->client_data->ast->n_values;
    return done;
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
    return 0;
}

int xRowid(sqlite3_vtab_cursor* pVtabCursor, sqlite3_int64* pRowid) {
    return 0;
}

int xUpdate(sqlite3_vtab* pVtabCursor, int argc, sqlite3_value** argv,
            sqlite3_int64* pRowid) {
    return 0;
}

int xBegin(sqlite3_vtab* pVTab) { return 0; }

int xSync(sqlite3_vtab* pVTab) { return 0; }

int xCommit(sqlite3_vtab* pVTab) { return 0; }

int xRollback(sqlite3_vtab* pVTab) { return 0; }

int xFindFunction(sqlite3_vtab* pVtab, int nArg, const char* zName,
                  void (**pxFunc)(sqlite3_context*, int, sqlite3_value**),
                  void** ppArg) {
    return 0;
}

int xRename(sqlite3_vtab* pVtab, const char* zNew) { return 0; }

int xSavepoint(sqlite3_vtab* pVTab, int n) { return 0; }

int xRelease(sqlite3_vtab* pVTab, int n) { return 0; }

int xRollbackTo(sqlite3_vtab* pVTab, int n) { return 0; }

int xShadowName(const char* name) { return 0; }

int setup_sqlite3(sqlite3* db, ClientData* client_data) {
    // Create a new virtual table and register methods.
    sqlite3_module json_vtab = {.iVersion = 1,
                                .xCreate = &xCreate,
                                .xConnect = &xConnect,
                                .xBestIndex = &xBestIndex,
                                .xDisconnect = &xDisconnect,
                                .xDestroy = &xDestroy,
                                .xOpen = &xOpen,
                                .xClose = &xClose,
                                .xFilter = &xFilter,
                                .xNext = &xNext,
                                .xEof = &xEof,
                                .xColumn = &xColumn,
                                .xRowid = &xRowid,
                                .xUpdate = &xUpdate,
                                .xBegin = &xBegin,
                                .xSync = &xSync,
                                .xCommit = &xCommit,
                                .xRollback = &xRollback,
                                .xFindFunction = &xFindFunction,
                                .xRename = &xRename,
                                .xSavepoint = &xSavepoint,
                                .xRelease = &xRelease,
                                .xRollbackTo = &xRollbackTo,
                                .xShadowName = &xShadowName};
    json_vtab.xBegin = &xBegin;

    // Open a new database connection.
    int rc = sqlite3_open(":memory:", &db);
    if (rc != 0) {
        fprintf(stderr, "Something went wrong!\n");
        sqlite3_close(db);
        return 1;
    }

    // Register the virtual table.
    rc = sqlite3_create_module(db, "sqjson", &json_vtab, (void*)client_data);
    if (rc != 0) {
        printf("Something went wrong!\n");
        sqlite3_close(db);
        return 1;
    }

    char* error_msg;
    rc = sqlite3_exec(db, "CREATE VIRTUAL TABLE [] USING sqjson", NULL, NULL,
                      &error_msg);
    if (rc != 0) {
        fprintf(stderr, "%s\n", error_msg);
        sqlite3_close(db);
        return 1;
    }

    rc = sqlite3_exec(db, client_data->query, &row_callback, (void*)client_data,
                      &error_msg);
    if (rc != 0) {
        fprintf(stderr, "%s\n", error_msg);
        sqlite3_close(db);
        return 1;
    }

    return 0;
}

int main(int argc, char** argv) {
    // First argument should be a filename.
    if (argc < 3) {
        log_and_exit("usage: sqjson INPUT_FILE QUERY\n");
    }

    // Get the size of our input file.
    FILE* f = fopen(argv[1], "r");
    fseek(f, 0, SEEK_END);
    const size_t input_size = ftell(f);
    fseek(f, 0, SEEK_SET);

    // Read our input file.
    char* input_data = calloc(input_size + 1, sizeof(char));
    const size_t bytes_read =
        fread(input_data, sizeof(*input_data), input_size, f);
    if (bytes_read != input_size) {
        log_and_exit("failed to read input");
    }

    // Read our query.
    char* query = argv[2];

    // Tokenize the input data.
    Token* tokens = NULL;
    size_t n_tokens = 0;
    tokenize(input_data, &tokens, &n_tokens);

    // Parse our input.
    JSONNode* ast;
    parse(tokens, &ast);

    // Build the 'CREATE TABLE ...' statement.
    JSONTableSchema* schema;
    build_table_schema(ast, &schema);

    // Create a new ClientData instance to register with our module.
    ClientData client_data = {.ast = ast, .schema = schema, .query = query};

    // Setup a new database instance.
    sqlite3* db = NULL;

    int rc = setup_sqlite3(db, &client_data);
    if (rc != 0) {
        sqlite3_close(db);
        return rc;
    }

    // Time to wrap it up!.
    sqlite3_close(db);

    return 0;
}
