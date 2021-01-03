#ifndef JSON_VIRTUAL_TABLE_H
#define JSON_VIRTUAL_TABLE_H

#include <sqlite3.h>

#include "json_parse.h"
#include "json_schema.h"

typedef struct ClientData {
    JSONNode* ast;
    JSONTableSchema* schema;
    char* query;
    int columns_written;
} ClientData;

int setup_virtual_table(sqlite3* db, ClientData* client_data);

int exec(sqlite3* db, ClientData* client_data);

#endif /// JSON_VIRTUAL_TABLE_H
