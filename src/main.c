#include <sqlite3.h>
#include <stdio.h>
#include <string.h>

#include "json_parse.h"
#include "json_schema.h"
#include "json_tokenize.h"
#include "json_virtual_table.h"
#include "util.h"

typedef struct ShellOptions {
    int compact;
    int nth;
} ShellOptions;

// Print usage information and exit.
void usage() {
    fprintf(stdout,
            "Usage: sqj [OPTION]... <SQL> [FILE]...\n"
            "Query JSON with SQL.\n"
            "\n"
            "\t--help    Display this message and exit\n"
            "\t--compact Format output without any extraneous whitespace\n");
    exit(EXIT_FAILURE);
}

int main(int argc, char** argv) {
    int rc = 0;

    // Fail early for obviously invalid usage.
    if (argc < 2) {
        usage();
    }

    // Parse command line options.
    ShellOptions shell_options = {.nth = -1};
    int i;
    for (i = 1; i < argc; i++) {
        char* z = argv[i];
        if (z[0] != '-') {
            break; // End of options.
        }
        if (z[1] == '-') {
            ++z; // Trim long options.
        }

        if (strcmp(z, "-help") == 0) {
            usage();
        } else if (strcmp(z, "-compact") == 0) {
            shell_options.compact = 1;
        } else if (strcmp(z, "-nth") == 0) {
            shell_options.nth = atoi(argv[++i]);
        }
    }

    // The query string should be the first argument after options.
    char* query = argv[i++];

    // Excess arguments after the query string are treated as files and mean we
    // do not read from stdin. A single file named "-" is  treated as an alias
    // for stdin.
    FILE* fin = stdin;
    if (i < argc && strcmp(argv[i], "-") != 0) {
        fin = fopen(argv[i], "r");
        if (fin == NULL) {
            log_and_exit("failed to open %s\n", argv[i]);
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
    fclose(mem_stream);

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

    // Query the table.
    if (exec(&client_data)) {
        fprintf(stderr, "failed to run query\n");
        rc = 1;
        goto clean_up;
    }

    // Output the results.
    if (shell_options.nth > 0) {
        // User requested a specified value.
        if (client_data.result_ast->value != JSON_VALUE_ARRAY) {
            fprintf(stderr, "result is not an array\n");
            rc = 1;
            goto clean_up;
        }
        if (shell_options.nth >= client_data.result_ast->n_values) {
            fprintf(stderr, "index out of range\n");
            rc = 1;
            goto clean_up;
        }
        pretty_print(&client_data.result_ast->values[shell_options.nth], stdout,
                     shell_options.compact);
    } else {
        pretty_print(client_data.result_ast, stdout, shell_options.compact);
    }

    // Time to wrap it up!.
clean_up:
    free(input_data);
    delete_tokens(tokens, n_tokens);
    delete_table_schema(schema);

    delete_ast(ast);
    free(ast);

    delete_ast(client_data.result_ast);
    free(client_data.result_ast);

    return rc;
}
