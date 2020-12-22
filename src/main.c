#include <stdio.h>
#include <sqlite3.h>
#include <string.h>

void debug_print(sqlite3_index_info* info) {
    printf("=== debug_print sqlite3_index_info ===\n");
    for (int i = 0; i < info->nConstraint; i++) {
        printf(
                "Column: %d\n"
                "Operator: %d\n",
                info->aConstraint[i].iColumn,
                info->aConstraint[i].op
        );
    }
    printf("=== debug_print sqlite3_index_info ===\n");
}

typedef struct json_vtab json_vtab;

struct json_vtab {
    sqlite3_vtab base;
    sqlite3 *db;
    int test;
};

typedef struct json_vtab_cursor json_vtab_cursor;

struct json_vtab_cursor {
    sqlite3_vtab_cursor *base;
    char *input;
    int count;
};

int xCreate(sqlite3 *db, void *pAux,
            int argc, const char *const *argv,
            sqlite3_vtab **ppVTab, char **pzErr) {
    // Log some debug info.
    printf("in xCreate(...)\n");

    // Declare the schema of our virtual table.
    sqlite3_declare_vtab(
            db,
            "CREATE TABLE json (name TEXT, value INT)"
    );

    // Allocate a new json_vtab instance.
    json_vtab *vtab = (json_vtab *) sqlite3_malloc(sizeof(json_vtab));
    *ppVTab = &vtab->base;

    return SQLITE_OK;
}

int xConnect(sqlite3 *db, void *pAux,
             int argc, const char *const *argv,
             sqlite3_vtab **ppVTab, char **pzErr) {
    // Log some debug info.
    printf("in xConnect(...)\n");
    // xConnect is just xCreate for ephemeral virtual tables.
    return xCreate(db, pAux, argc, argv, ppVTab, pzErr);
}

int xBestIndex(sqlite3_vtab *pVTab, sqlite3_index_info *pIndexInfo) {
    printf("in xBestIndex(...)\n");
    debug_print(pIndexInfo);
    return 0;
}

int xDisconnect(sqlite3_vtab *pVTab) {
    printf("in xDisconnect(...)\n");
    return 0;
}

int xDestroy(sqlite3_vtab *pVTab) {
    printf("in xDestroy(...)\n");
    return 0;
}

int xOpen(sqlite3_vtab *pVTab, sqlite3_vtab_cursor **ppCursor) {
    // Log some debug info.
    printf("in xOpen(...)\n");

    // Open a new cursor.
    json_vtab_cursor* cursor = sqlite3_malloc(sizeof(json_vtab_cursor));
    memset(cursor, 0, sizeof(json_vtab_cursor));
    *ppCursor = (sqlite3_vtab_cursor*)cursor;
    return SQLITE_OK;
}

int xClose(sqlite3_vtab_cursor *pVtabCursor) {
    printf("in xClose(...)\n");
    return 0;
}

int xFilter(sqlite3_vtab_cursor *pVtabCursor, int idxNum, const char *idxStr,
            int argc, sqlite3_value **argv) {
    printf("in xFilter(...)\n");
    return 0;
}

int xNext(sqlite3_vtab_cursor *pVtabCursor) {
    printf("in xNext(...)\n");
    json_vtab_cursor* cursor = (json_vtab_cursor*)pVtabCursor;
    (cursor->count)++;
    return SQLITE_OK;
}

int xEof(sqlite3_vtab_cursor *pVtabCursor) {
    printf("in xEof(...)\n");
    json_vtab_cursor* cursor = (json_vtab_cursor*)pVtabCursor;
    const int done = cursor->count > 10;
    return done;
}

int xColumn(sqlite3_vtab_cursor *pVtabCursor, sqlite3_context *pContext, int n) {
    printf("in xColumn(...)\n");
    printf("count: %d\n", ((json_vtab_cursor*)pVtabCursor)->count);
    return 0;
}

int xRowid(sqlite3_vtab_cursor *pVtabCursor, sqlite3_int64 *pRowid) {
    printf("in xRowid(...)\n");
    return 0;
}

int xUpdate(sqlite3_vtab *pVtabCursor, int argc, sqlite3_value **argv, sqlite3_int64 *pRowid) {
    printf("in xUpdate(...)\n");
    return 0;
}

int xBegin(sqlite3_vtab *pVTab) {
    printf("in xBegin(...)\n");
    return 0;
}

int xSync(sqlite3_vtab *pVTab) {
    printf("in xSync(...)\n");
    return 0;
}

int xCommit(sqlite3_vtab *pVTab) {
    printf("in xCommit(...)\n");
    return 0;
}

int xRollback(sqlite3_vtab *pVTab) {
    printf("in xRollback(...)\n");
    return 0;
}

int xFindFunction(sqlite3_vtab *pVtab, int nArg, const char *zName,
                  void (**pxFunc)(sqlite3_context *, int, sqlite3_value **),
                  void **ppArg) {
    printf("in xFindFunction(...)\n");
    return 0;
}

int xRename(sqlite3_vtab *pVtab, const char *zNew) {
    printf("in xRename(...)\n");
    return 0;
}

int xSavepoint(sqlite3_vtab *pVTab, int n) {
    printf("in xSavepoint(...)\n");
    return 0;
}

int xRelease(sqlite3_vtab *pVTab, int n) {
    printf("in xRelease(...)\n");
    return 0;
}

int xRollbackTo(sqlite3_vtab *pVTab, int n) {
    printf("in xRollbackTo(...)\n");
    return 0;
}

int xShadowName(const char *name) {
    printf("in xShadowName(...)\n");
    return 0;
}

int setup_sqlite3(sqlite3* db) {
    // Create a new virtual table and register methods.
    sqlite3_module json_vtab = {
            .iVersion = 1,
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
            .xShadowName = &xShadowName
    };
    json_vtab.xBegin = &xBegin;

    // Open a new database connection.
    printf("Opening Sqlite3 database...\n");
    int rc = sqlite3_open(":memory:", &db);
    if (rc != 0) {
        printf("Something went wrong!\n");
        sqlite3_close(db);
        return 1;
    }
    printf("Done!\n");

    // Register the virtual table.
    printf("Creating module...\n");
    rc = sqlite3_create_module(db, "sqjson", &json_vtab, NULL);
    if (rc != 0) {
        printf("Something went wrong!\n");
        sqlite3_close(db);
        return 1;
    }
    printf("Done!\n");

    printf("Create the virtual table\n");
    char *error_msg;
    rc = sqlite3_exec(
            db,
            "CREATE VIRTUAL TABLE sqjson USING sqjson",
            NULL,
            NULL,
            &error_msg
    );
    if (rc != 0) {
        printf("Something went wrong!\n");
        printf("%s\n", error_msg);
        sqlite3_close(db);
        return 1;
    }
    printf("Done!\n");

    printf("Querying virtual table\n");
    rc = sqlite3_exec(
            db,
            "SELECT * FROM sqjson WHERE value < 5",
            NULL,
            NULL,
            &error_msg
    );
    if (rc != 0) {
        printf("Something went wrong!\n");
        printf("%s\n", error_msg);
        sqlite3_close(db);
        return 1;
    }
    printf("Done!\n");

    return 0;
}

int main(int argc, char** argv) {
    // Setup a new database instance.
    sqlite3 *db = NULL;

    int rc = setup_sqlite3(db);
    if (rc != 0) {
        sqlite3_close(db);
        return rc;
    }

    // Time to wrap it up!.
    sqlite3_close(db);

    return 0;
}
