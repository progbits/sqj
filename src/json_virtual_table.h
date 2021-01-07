#ifndef JSON_VIRTUAL_TABLE_H
#define JSON_VIRTUAL_TABLE_H

#include <sqlite3.h>

#include "json_parse.h"
#include "json_schema.h"

typedef struct ClientData {
    JSONNode* ast;
    JSONNode* result_ast;
    JSONTableSchema* schema;
    char* query;
    int columns_written;
    int row;
} ClientData;

int exec(ClientData* client_data);

#endif /// JSON_VIRTUAL_TABLE_H
