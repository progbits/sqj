#ifndef UTIL_H
#define UTIL_H

void log_and_exit(const char* format, ...);

// Escape a string such that it forms a valid column identifier.
char* escape_string(const char* str);

// Return an escaped string in its original form.
char* unescape_string(const char* str);

#endif // UTIL_H
