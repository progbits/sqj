package sql

import (
	"strings"
	"unicode"
	"unicode/utf8"
)

// Token represents a token scanned from an input stream.
type Token uint8

const (
	INVALID Token = iota
	EOF

	// operators
	MINUS  // -
	LP     // (
	RP     // )
	SEMI   // ;
	PLUS   // +
	STAR   // *
	SLASH  // /
	REM    // %
	EQ     // =
	LE     // <=
	LSHIFT // <<
	LT     // <
	GE     // >=
	RSHIFT // >>
	GT     // >
	NE     // !=
	COMMA  // ,
	BITAND // &
	BITNOT // !
	BITOR  // |
	CONCAT // ||
	DOT    // .

	// literals
	IDENTIFIER      // Table, column name, alias etc...
	NUMERIC_LITERAL // Either INTEGER_LITERAL or FLOAT_LITERAL
	STRING_LITERAL  // TODO: Handle string literals.

	// keywords
	ABORT
	ACTION
	ADD
	AFTER
	ALL
	ALTER
	ALWAYS
	ANALYZE
	AND
	AS
	ASC
	ATTACH
	AUTOINCREMENT
	BEFORE
	BEGIN
	BETWEEN
	BY
	CASCADE
	CASE
	CAST
	CHECK
	COLLATE
	COLUMN
	COMMIT
	CONFLICT
	CONSTRAINT
	CREATE
	CROSS
	CURRENT
	CURRENT_DATE
	CURRENT_TIME
	CURRENT_TIMESTAMP
	DATABASE
	DEFAULT
	DEFERRABLE
	DEFERRED
	DELETE
	DESC
	DETACH
	DISTINCT
	DO
	DROP
	EACH
	ELSE
	END
	ESCAPE
	EXCEPT
	EXCLUDE
	EXCLUSIVE
	EXISTS
	EXPLAIN
	FAIL
	FILTER
	FIRST
	FOLLOWING
	FOR
	FOREIGN
	FROM
	FULL
	GENERATED
	GLOB
	GROUP
	GROUPS
	HAVING
	IF
	IGNORE
	IMMEDIATE
	IN
	INDEX
	INDEXED
	INITIALLY
	INNER
	INSERT
	INSTEAD
	INTERSECT
	INTO
	IS
	ISNULL
	JOIN
	KEY
	LAST
	LEFT
	LIKE
	LIMIT
	MATCH
	NATURAL
	NO
	NOT
	NOTHING
	NOTNULL
	NULL
	NULLS
	OF
	OFFSET
	ON
	OR
	ORDER
	OTHERS
	OUTER
	OVER
	PARTITION
	PLAN
	PRAGMA
	PRECEDING
	PRIMARY
	QUERY
	RAISE
	RANGE
	RECURSIVE
	REFERENCES
	REGEXP
	REINDEX
	RELEASE
	RENAME
	REPLACE
	RESTRICT
	RIGHT
	ROLLBACK
	ROW
	ROWS
	SAVEPOINT
	SELECT
	SET
	TABLE
	TEMP
	TEMPORARY
	THEN
	TIES
	TO
	TRANSACTION
	TRIGGER
	UNBOUNDED
	UNION
	UNIQUE
	UPDATE
	USING
	VACUUM
	VALUES
	VIEW
	VIRTUAL
	WHEN
	WHERE
	WINDOW
	WITH
	WITHOUT
)

var tokens = [...]string{
	INVALID: "INVALID",
	EOF:     "EOF",

	// operators
	MINUS:  "-",
	LP:     "(",
	RP:     ")",
	SEMI:   ";",
	PLUS:   "+",
	STAR:   "*",
	SLASH:  "/",
	REM:    "%",
	EQ:     "=",
	LE:     "<=",
	LSHIFT: "<<",
	LT:     "<",
	GE:     ">=",
	RSHIFT: ">>",
	GT:     ">",
	NE:     "!=",
	COMMA:  ",",
	BITAND: "&",
	BITNOT: "!",
	BITOR:  "|",
	CONCAT: "||",
	DOT:    ".",

	// literals
	IDENTIFIER:      "IDENTIFIER",
	NUMERIC_LITERAL: "NUMERIC_LITERAL",
	STRING_LITERAL:  "STRING_LITERAL",

	// keywords
	ABORT:             "ABORT",
	ACTION:            "ACTION",
	ADD:               "ADD",
	AFTER:             "AFTER",
	ALL:               "ALL",
	ALTER:             "ALTER",
	ALWAYS:            "ALWAYS",
	ANALYZE:           "ANALYZE",
	AND:               "AND",
	AS:                "AS",
	ASC:               "ASC",
	ATTACH:            "ATTACH",
	AUTOINCREMENT:     "AUTOINCREMENT",
	BEFORE:            "BEFORE",
	BEGIN:             "BEGIN",
	BETWEEN:           "BETWEEN",
	BY:                "BY",
	CASCADE:           "CASCADE",
	CASE:              "CASE",
	CAST:              "CAST",
	CHECK:             "CHECK",
	COLLATE:           "COLLATE",
	COLUMN:            "COLUMN",
	COMMIT:            "COMMIT",
	CONFLICT:          "CONFLICT",
	CONSTRAINT:        "CONSTRAINT",
	CREATE:            "CREATE",
	CROSS:             "CROSS",
	CURRENT:           "CURRENT",
	CURRENT_DATE:      "CURRENT_DATE",
	CURRENT_TIME:      "CURRENT_TIME",
	CURRENT_TIMESTAMP: "CURRENT_TIMESTAMP",
	DATABASE:          "DATABASE",
	DEFAULT:           "DEFAULT",
	DEFERRABLE:        "DEFERRABLE",
	DEFERRED:          "DEFERRED",
	DELETE:            "DELETE",
	DESC:              "DESC",
	DETACH:            "DETACH",
	DISTINCT:          "DISTINCT",
	DO:                "DO",
	DROP:              "DROP",
	EACH:              "EACH",
	ELSE:              "ELSE",
	END:               "END",
	ESCAPE:            "ESCAPE",
	EXCEPT:            "EXCEPT",
	EXCLUDE:           "EXCLUDE",
	EXCLUSIVE:         "EXCLUSIVE",
	EXISTS:            "EXISTS",
	EXPLAIN:           "EXPLAIN",
	FAIL:              "FAIL",
	FILTER:            "FILTER",
	FIRST:             "FIRST",
	FOLLOWING:         "FOLLOWING",
	FOR:               "FOR",
	FOREIGN:           "FOREIGN",
	FROM:              "FROM",
	FULL:              "FULL",
	GENERATED:         "GENERATED",
	GLOB:              "GLOB",
	GROUP:             "GROUP",
	GROUPS:            "GROUPS",
	HAVING:            "HAVING",
	IF:                "IF",
	IGNORE:            "IGNORE",
	IMMEDIATE:         "IMMEDIATE",
	IN:                "IN",
	INDEX:             "INDEX",
	INDEXED:           "INDEXED",
	INITIALLY:         "INITIALLY",
	INNER:             "INNER",
	INSERT:            "INSERT",
	INSTEAD:           "INSTEAD",
	INTERSECT:         "INTERSECT",
	INTO:              "INTO",
	IS:                "IS",
	ISNULL:            "ISNULL",
	JOIN:              "JOIN",
	KEY:               "KEY",
	LAST:              "LAST",
	LEFT:              "LEFT",
	LIKE:              "LIKE",
	LIMIT:             "LIMIT",
	MATCH:             "MATCH",
	NATURAL:           "NATURAL",
	NO:                "NO",
	NOT:               "NOT",
	NOTHING:           "NOTHING",
	NOTNULL:           "NOTNULL",
	NULL:              "NULL",
	NULLS:             "NULLS",
	OF:                "OF",
	OFFSET:            "OFFSET",
	ON:                "ON",
	OR:                "OR",
	ORDER:             "ORDER",
	OTHERS:            "OTHERS",
	OUTER:             "OUTER",
	OVER:              "OVER",
	PARTITION:         "PARTITION",
	PLAN:              "PLAN",
	PRAGMA:            "PRAGMA",
	PRECEDING:         "PRECEDING",
	PRIMARY:           "PRIMARY",
	QUERY:             "QUERY",
	RAISE:             "RAISE",
	RANGE:             "RANGE",
	RECURSIVE:         "RECURSIVE",
	REFERENCES:        "REFERENCES",
	REGEXP:            "REGEXP",
	REINDEX:           "REINDEX",
	RELEASE:           "RELEASE",
	RENAME:            "RENAME",
	REPLACE:           "REPLACE",
	RESTRICT:          "RESTRICT",
	RIGHT:             "RIGHT",
	ROLLBACK:          "ROLLBACK",
	ROW:               "ROW",
	ROWS:              "ROWS",
	SAVEPOINT:         "SAVEPOINT",
	SELECT:            "SELECT",
	SET:               "SET",
	TABLE:             "TABLE",
	TEMP:              "TEMP",
	TEMPORARY:         "TEMPORARY",
	THEN:              "THEN",
	TIES:              "TIES",
	TO:                "TO",
	TRANSACTION:       "TRANSACTION",
	TRIGGER:           "TRIGGER",
	UNBOUNDED:         "UNBOUNDED",
	UNION:             "UNION",
	UNIQUE:            "UNIQUE",
	UPDATE:            "UPDATE",
	USING:             "USING",
	VACUUM:            "VACUUM",
	VALUES:            "VALUES",
	VIEW:              "VIEW",
	VIRTUAL:           "VIRTUAL",
	WHEN:              "WHEN",
	WHERE:             "WHERE",
	WINDOW:            "WINDOW",
	WITH:              "WITH",
	WITHOUT:           "WITHOUT",
}

func (t Token) String() string {
	return tokens[t]
}

// Scanner represents a type that scans tokens.
type Scanner struct {
	input  []byte // Statement source text.
	cursor int    // Current position in `input`.
	char   rune   // Current character, -1 for EOF.
}

// NewScanner creates a new scanner from a statement.
func NewScanner(statement []byte) *Scanner {
	scanner := &Scanner{
		input: statement,
	}
	scanner.next()

	return scanner
}

// ScanToken scans the next token, returning the token and its value (if applicable).
func (s *Scanner) ScanToken() (Token, string) {
	s.skipWhitespace()

	// Greedily match identifiers, keywords and literals.
	switch c := s.char; {
	case c == 0:
		return EOF, ""
	case isLetter(c):
		value := string(c)
		for s.next(); isLetter(s.char) || unicode.IsDigit(s.char); s.next() {
			value += string(s.char)
		}

		for i := ABORT; i <= WITHOUT; i++ {
			if strings.ToUpper(value) == tokens[i] {
				return i, value
			}
		}
		return IDENTIFIER, value
	case isDigit(c) || c == '.' && isDigit(s.peek()):
		value := string(c)
		for s.next(); isLetter(s.char) || isDigit(s.char) || s.char == '.' || s.char == '-'; s.next() {
			value += string(s.char)
		}
		return NUMERIC_LITERAL, value
	case c == '\'' || c == '"':
		s.next()
		value := ""
		for {
			if s.cursor >= len(s.input)-1 || s.char == '\'' || s.char == '"' {
				s.next()
				break
			}
			value += string(s.char)
			s.next()
		}
		return STRING_LITERAL, value
	}

	// Match all other tokens.
	c := s.char
	s.next()
	var token Token

	switch c {
	case '-':
		token = MINUS
	case '(':
		token = LP
	case ')':
		token = RP
	case ';':
		token = SEMI
	case '+':
		token = PLUS
	case '*':
		token = STAR
	case '/':
		token = SLASH
	case '%':
		token = REM
	case ',':
		token = COMMA
	case '&':
		token = BITAND
	case '~':
		token = BITNOT
	case '.':
		token = DOT
	default:
		peek := s.char
		switch c {
		case '=':
			if peek == '=' {
				token = EQ
				s.next()
			} else {
				token = EQ
			}
		case '<':
			// Might either be <=, <>, << or <.
			if peek == '=' {
				token = LE
				s.next()
			} else if peek == '>' {
				token = NE
				s.next()
			} else if peek == '<' {
				token = LSHIFT
				s.next()
			} else {
				token = LT
			}
		case '>':
			// Might either be >=, >> or >.
			if peek == '=' {
				token = GE
				s.next()
			} else if peek == '>' {
				token = RSHIFT
				s.next()
			} else {
				token = GT
			}
		case '!':
			// Might either be ! or !=.
			if peek == '=' {
				token = NE
				s.next()
			} else {
				token = BITNOT
			}
		case '|':
			// Might either be ! or !=.
			if peek == '|' {
				token = CONCAT
				s.next()
			} else {
				token = BITOR
			}
		default:
			token = INVALID
		}
	}

	return token, ""
}

// next reads the next character from the input, or sets `s.char = -1`.
// TODO: Actually handle unicode.
func (s *Scanner) next() {
	if s.cursor < len(s.input) {
		if cur := rune(s.input[s.cursor]); cur == utf8.RuneError {
			panic("malformed input")
		} else if cur >= utf8.RuneSelf {
			char, size := utf8.DecodeRune(s.input[s.cursor:])
			s.char = char
			s.cursor += size
		} else {
			s.char = cur
			s.cursor += 1
		}
	} else {
		s.char = 0 // Unicode NULL.
	}
}

// peek returns the next character from the input without advancing the cursor.
// TODO: Actually handle unicode.
func (s *Scanner) peek() rune {
	if s.cursor < len(s.input) {
		if cur := rune(s.input[s.cursor]); cur == utf8.RuneError {
			panic("malformed input")
		} else if cur >= utf8.RuneSelf {
			char, _ := utf8.DecodeRune(s.input[s.cursor:])
			return char
		} else {
			return cur
		}
	}
	return 0 // Unicode NULL.
}

// skipWhitespace advances the cursor to the next non-whitespace character.
func (s *Scanner) skipWhitespace() {
	for unicode.IsSpace(s.char) {
		s.next()
	}
}

// isLetter returns true if char is and underscore or a letter.
func isLetter(char rune) bool {
	if unicode.IsLetter(char) {
		return true
	}

	switch char {
	case '_', '[', ']', '$':
		return true
	default:
		return false
	}
}

// isDigit returns true if char is a digit between 0 and 9.
func isDigit(char rune) bool {
	return char >= '0' && char <= '9'
}
