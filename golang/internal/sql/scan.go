package sql

import (
	"strings"
	"unicode"
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
	scanner := &Scanner{statement, 0, rune(statement[0])}
	return scanner
}

// ScanToken scans the next token, returning the token and its value (if applicable).
func (s *Scanner) ScanToken() (Token, string) {
	if s.cursor == len(s.input) {
		return EOF, ""
	}
	s.skipWhitespace()

	// Greedily match identifiers, keywords and literals.
	switch c := s.char; {
	case isLetter(c):
		idStart := s.cursor
		for isLetter(s.char) || isDigit(s.char) {
			s.next()
		}
		value := string(s.input[idStart:s.cursor])

		for i := ABORT; i <= WITHOUT; i++ {
			if strings.ToUpper(value) == tokens[i] {
				return i, value
			}
		}
		return IDENTIFIER, value
	case isDigit(c) || c == '.' && isDigit(s.peek()):
		iStart := s.cursor
		for isLetter(s.char) || isDigit(s.char) || s.char == '.' || s.char == '-' {
			s.next()
		}
		value := string(s.input[iStart:s.cursor])
		return NUMERIC_LITERAL, value
	case c == '\'':
		// String literal, consume to peek '\'', ignoring nested '\'' for now.
		bytes := make([]byte, 0)
		for {
			s.cursor += 1
			if s.cursor >= len(s.input)-1 || s.input[s.cursor] == '\'' {
				s.cursor += 1
				break
			}
			bytes = append(bytes, s.input[s.cursor])
		}
		return STRING_LITERAL, ""
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
	if s.cursor == len(s.input)-1 {
		s.char = -1 // We've run out of input.
		return
	}
	s.cursor += 1
	s.char = rune(s.input[s.cursor])
}

// peek returns the next character from the input without advancing the cursor.
// TODO: Actually handle unicode.
func (s *Scanner) peek() rune {
	if s.cursor == len(s.input)-1 {
		return -1
	}
	return rune(s.input[s.cursor+1])
}

// skipWhitespace advances the cursor to the next non-whitespace character.
func (s *Scanner) skipWhitespace() {
	for unicode.IsSpace(rune(s.input[s.cursor])) {
		s.next()
	}
}

// isLetter returns true if char is and underscore or a letter.
func isLetter(char rune) bool {
	return char == '_' || char >= 'a' && char <= 'z' || char >= 'A' && char <= 'Z'
}

// isDigit returns true if char is a digit between 0 and 9.
func isDigit(char rune) bool {
	return char >= '0' && char <= '9'
}
