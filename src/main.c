#include <sqlite3.h>
#include <stdio.h>
#include <string.h>

#include "json_parse.h"
#include "json_schema.h"
#include "json_tokenize.h"
#include "json_virtual_table.h"
#include "util.h"

int main(int argc, char** argv) {
    // Our input file. This could be stdin or a named file. Initially assume
    // we are working with stdin.
    FILE* fin = stdin;

    if (argc < 2 || argc > 3) {
        fprintf(stdout, "sqj - Query JSON with SQL\n"
                        "Usage: sqj <SQL> [FILE]\n");
        exit(EXIT_FAILURE);
    }

    // We have a named file as our input. This might be an actual path or '-'
    // as an alias for stdin.
    if (argc > 2 && strcmp(argv[2], "-") != 0) {
        fin = fopen(argv[2], "r");
        if (fin == NULL) {
            log_and_exit("failed to open %s\n", argv[2]);
        }
    }

    // Read the input file to a buffer.
    char* input_data;
    size_t input_data_size;
    FILE* mem_stream = open_memstream(&input_data, &input_data_size);

    char buffer[1024];
    while (fgets(buffer, sizeof(buffer), fin)) {
        fputs(buffer, mem_stream);
    }
    fflush(mem_stream);
    fclose(mem_stream);

    // SQL query string should be the first argument.
    char* query = argv[1];

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
    if (sqlite3_open(":memory:", &db) != SQLITE_OK) {
        fprintf(stderr, "failed to open in-memory database\n");
        sqlite3_close(db);
        return EXIT_FAILURE;
    }

    if (setup_virtual_table(db, &client_data) != SQLITE_OK) {
        fprintf(stderr, "failed to setup database\n");
        sqlite3_close(db);
        return EXIT_FAILURE;
    }

    if (exec(db, &client_data) != SQLITE_OK) {
        fprintf(stderr, "failed to run query\n");
        sqlite3_close(db);
        return EXIT_FAILURE;
    }

    // Time to wrap it up!. Close our database connection and free our input
    // buffer.
    sqlite3_close(db);
    free(input_data);

    return EXIT_SUCCESS;
}
