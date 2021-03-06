package sql

import (
	"testing"
)

// A TokenValuePair represents a token and its value, if applicable.
type TokenValuePair struct {
	token Token
	value string
}

// A TestCase is a statement and a list of expected (token,value) pairs.
type TestCase struct {
	statement string
	expected  []TokenValuePair
}

// checkEquality tests to sets of TokenValuePairs for equality.
func checkEquality(t *testing.T, actual []TokenValuePair, expected []TokenValuePair) {
	if len(actual) != len(expected) {
		t.Fatalf("unexpected number of values: got %d, expected %d", len(actual), len(expected))
	}

	for i := 0; i < len(actual); i++ {
		if actual[i].token != expected[i].token {
			t.Fatalf("invalid token: expected %d, got %d", expected[i].token, actual[i].token)
		}
		if actual[i].value != expected[i].value {
			t.Fatalf("invalid token: expected %s, got %s", expected[i].value, actual[i].value)
		}
	}
}

// scanAll extracts all tokens and their values from a Scanner.
func scanAll(s *Scanner) (tokens []TokenValuePair) {
	tokens = make([]TokenValuePair, 0)
	for {
		token, value := s.ScanToken()
		if token == INVALID || token == EOF {
			break
		}
		tokens = append(tokens, TokenValuePair{token, value})
	}
	return
}

func TestIdentifiers(t *testing.T) {
	var cases = []TestCase{
		{"[];", []TokenValuePair{{IDENTIFIER, "[]"}, {SEMI, ""}}},
		{"α;", []TokenValuePair{
			{IDENTIFIER, "α"},
			{SEMI, ""},
		}},
		{"α, γ, δ, ϵ, ζ, θ, μ, ψ;", []TokenValuePair{
			{IDENTIFIER, "α"},
			{COMMA, ""},
			{IDENTIFIER, "γ"},
			{COMMA, ""},
			{IDENTIFIER, "δ"},
			{COMMA, ""},
			{IDENTIFIER, "ϵ"},
			{COMMA, ""},
			{IDENTIFIER, "ζ"},
			{COMMA, ""},
			{IDENTIFIER, "θ"},
			{COMMA, ""},
			{IDENTIFIER, "μ"},
			{COMMA, ""},
			{IDENTIFIER, "ψ"},
			{SEMI, ""},
		}},
		{"test;", []TokenValuePair{{IDENTIFIER, "test"}, {SEMI, ""}}},
		{"test_table;", []TokenValuePair{{IDENTIFIER, "test_table"}, {SEMI, ""}}},
		{"_test_table;", []TokenValuePair{{IDENTIFIER, "_test_table"}, {SEMI, ""}}},
		{"_t1e2s3t_4t5a6l7e;", []TokenValuePair{{IDENTIFIER, "_t1e2s3t_4t5a6l7e"}, {SEMI, ""}}},
		{"test$;", []TokenValuePair{{IDENTIFIER, "test$"}, {SEMI, ""}}},
		{"test$table;", []TokenValuePair{{IDENTIFIER, "test$table"}, {SEMI, ""}}},
		{"test_database.test_table.test_column;",
			[]TokenValuePair{
				{IDENTIFIER, "test_database"},
				{DOT, ""},
				{IDENTIFIER, "test_table"},
				{DOT, ""},
				{IDENTIFIER, "test_column"},
				{SEMI, ""},
			},
		},
		{"_test_database._test_table._test_column;",
			[]TokenValuePair{
				{IDENTIFIER, "_test_database"},
				{DOT, ""},
				{IDENTIFIER, "_test_table"},
				{DOT, ""},
				{IDENTIFIER, "_test_column"},
				{SEMI, ""},
			},
		},
	}

	for _, _case := range cases {
		scanner := NewScanner([]byte(_case.statement))
		tokens := scanAll(scanner)
		checkEquality(t, tokens, _case.expected)
	}
}

func TestStringLiterals(t *testing.T) {
	var cases = []TestCase{
		{"\"hello, world\";", []TokenValuePair{{STRING_LITERAL, "hello, world"}, {SEMI, ""}}},
	}

	for _, _case := range cases {
		scanner := NewScanner([]byte(_case.statement))
		tokens := scanAll(scanner)
		checkEquality(t, tokens, _case.expected)
	}
}

func TestRealNumbers(t *testing.T) {
	var cases = []TestCase{
		//{"42.;", []TokenValuePair{{NUMERIC_LITERAL, "42."}, {SEMI, ""}}},
		//{"4e2;", []TokenValuePair{{NUMERIC_LITERAL, "4e2"}, {SEMI, ""}}},
		//{"123.456e-78;", []TokenValuePair{{NUMERIC_LITERAL, "123.456e-78"}, {SEMI, ""}}},
		{".1E2;", []TokenValuePair{{NUMERIC_LITERAL, ".1E2"}, {SEMI, ""}}},
	}

	for _, _case := range cases {
		scanner := NewScanner([]byte(_case.statement))
		tokens := scanAll(scanner)
		checkEquality(t, tokens, _case.expected)
	}
}

func TestIntegers(t *testing.T) {
	var cases = []TestCase{
		{"0;", []TokenValuePair{{NUMERIC_LITERAL, "0"}, {SEMI, ""}}},
		{"1;", []TokenValuePair{{NUMERIC_LITERAL, "1"}, {SEMI, ""}}},
		{"1234;", []TokenValuePair{{NUMERIC_LITERAL, "1234"}, {SEMI, ""}}},
		{"0x0;", []TokenValuePair{{NUMERIC_LITERAL, "0x0"}, {SEMI, ""}}},
		{"0x1;", []TokenValuePair{{NUMERIC_LITERAL, "0x1"}, {SEMI, ""}}},
		{"0x321;", []TokenValuePair{{NUMERIC_LITERAL, "0x321"}, {SEMI, ""}}},
		{"0xa;", []TokenValuePair{{NUMERIC_LITERAL, "0xa"}, {SEMI, ""}}},
		{"0xA;", []TokenValuePair{{NUMERIC_LITERAL, "0xA"}, {SEMI, ""}}},
		{"0xCDE;", []TokenValuePair{{NUMERIC_LITERAL, "0xCDE"}, {SEMI, ""}}},
		{"0xCdE;", []TokenValuePair{{NUMERIC_LITERAL, "0xCdE"}, {SEMI, ""}}},
		{"0xabcde012345;", []TokenValuePair{{NUMERIC_LITERAL, "0xabcde012345"}, {SEMI, ""}}},
		{"0xABCDE012345;", []TokenValuePair{{NUMERIC_LITERAL, "0xABCDE012345"}, {SEMI, ""}}},
		{"0xaBcDe012345;", []TokenValuePair{{NUMERIC_LITERAL, "0xaBcDe012345"}, {SEMI, ""}}},
		{"0xAbCdE012345;", []TokenValuePair{{NUMERIC_LITERAL, "0xAbCdE012345"}, {SEMI, ""}}},
		{"0xa1b2c3d4e;", []TokenValuePair{{NUMERIC_LITERAL, "0xa1b2c3d4e"}, {SEMI, ""}}},
		{"0xA1B2C3D4E;", []TokenValuePair{{NUMERIC_LITERAL, "0xA1B2C3D4E"}, {SEMI, ""}}},
		{"0xA1b2C3d4E;", []TokenValuePair{{NUMERIC_LITERAL, "0xA1b2C3d4E"}, {SEMI, ""}}},
	}

	for _, _case := range cases {
		scanner := NewScanner([]byte(_case.statement))
		tokens := scanAll(scanner)
		checkEquality(t, tokens, _case.expected)
	}
}

func TestKeywords(t *testing.T) {
	var cases = []TestCase{
		{"ABORT  abort;", []TokenValuePair{{ABORT, "ABORT"}, {ABORT, "abort"}, {SEMI, ""}}},
		{"ACTION  action;", []TokenValuePair{{ACTION, "ACTION"}, {ACTION, "action"}, {SEMI, ""}}},
		{"ADD  add;", []TokenValuePair{{ADD, "ADD"}, {ADD, "add"}, {SEMI, ""}}},
		{"AFTER  after;", []TokenValuePair{{AFTER, "AFTER"}, {AFTER, "after"}, {SEMI, ""}}},
		{"ALL  all;", []TokenValuePair{{ALL, "ALL"}, {ALL, "all"}, {SEMI, ""}}},
		{"ALTER  alter;", []TokenValuePair{{ALTER, "ALTER"}, {ALTER, "alter"}, {SEMI, ""}}},
		{"ALWAYS  always;", []TokenValuePair{{ALWAYS, "ALWAYS"}, {ALWAYS, "always"}, {SEMI, ""}}},
		{"ANALYZE  analyze;", []TokenValuePair{{ANALYZE, "ANALYZE"}, {ANALYZE, "analyze"}, {SEMI, ""}}},
		{"AND  and;", []TokenValuePair{{AND, "AND"}, {AND, "and"}, {SEMI, ""}}},
		{"AS  as;", []TokenValuePair{{AS, "AS"}, {AS, "as"}, {SEMI, ""}}},
		{"ASC  asc;", []TokenValuePair{{ASC, "ASC"}, {ASC, "asc"}, {SEMI, ""}}},
		{"ATTACH  attach;", []TokenValuePair{{ATTACH, "ATTACH"}, {ATTACH, "attach"}, {SEMI, ""}}},
		{"AUTOINCREMENT  autoincrement;", []TokenValuePair{{AUTOINCREMENT, "AUTOINCREMENT"}, {AUTOINCREMENT, "autoincrement"}, {SEMI, ""}}},
		{"BEFORE  before;", []TokenValuePair{{BEFORE, "BEFORE"}, {BEFORE, "before"}, {SEMI, ""}}},
		{"BEGIN  begin;", []TokenValuePair{{BEGIN, "BEGIN"}, {BEGIN, "begin"}, {SEMI, ""}}},
		{"BETWEEN  between;", []TokenValuePair{{BETWEEN, "BETWEEN"}, {BETWEEN, "between"}, {SEMI, ""}}},
		{"BY  by;", []TokenValuePair{{BY, "BY"}, {BY, "by"}, {SEMI, ""}}},
		{"CASCADE  cascade;", []TokenValuePair{{CASCADE, "CASCADE"}, {CASCADE, "cascade"}, {SEMI, ""}}},
		{"CASE  case;", []TokenValuePair{{CASE, "CASE"}, {CASE, "case"}, {SEMI, ""}}},
		{"CAST  cast;", []TokenValuePair{{CAST, "CAST"}, {CAST, "cast"}, {SEMI, ""}}},
		{"CHECK  check;", []TokenValuePair{{CHECK, "CHECK"}, {CHECK, "check"}, {SEMI, ""}}},
		{"COLLATE  collate;", []TokenValuePair{{COLLATE, "COLLATE"}, {COLLATE, "collate"}, {SEMI, ""}}},
		{"COLUMN  column;", []TokenValuePair{{COLUMN, "COLUMN"}, {COLUMN, "column"}, {SEMI, ""}}},
		{"COMMIT  commit;", []TokenValuePair{{COMMIT, "COMMIT"}, {COMMIT, "commit"}, {SEMI, ""}}},
		{"CONFLICT  conflict;", []TokenValuePair{{CONFLICT, "CONFLICT"}, {CONFLICT, "conflict"}, {SEMI, ""}}},
		{"CONSTRAINT  constraint;", []TokenValuePair{{CONSTRAINT, "CONSTRAINT"}, {CONSTRAINT, "constraint"}, {SEMI, ""}}},
		{"CREATE  create;", []TokenValuePair{{CREATE, "CREATE"}, {CREATE, "create"}, {SEMI, ""}}},
		{"CROSS  cross;", []TokenValuePair{{CROSS, "CROSS"}, {CROSS, "cross"}, {SEMI, ""}}},
		{"CURRENT  current;", []TokenValuePair{{CURRENT, "CURRENT"}, {CURRENT, "current"}, {SEMI, ""}}},
		{"CURRENT_DATE  current_date;", []TokenValuePair{{CURRENT_DATE, "CURRENT_DATE"}, {CURRENT_DATE, "current_date"}, {SEMI, ""}}},
		{"CURRENT_TIME  current_time;", []TokenValuePair{{CURRENT_TIME, "CURRENT_TIME"}, {CURRENT_TIME, "current_time"}, {SEMI, ""}}},
		{"CURRENT_TIMESTAMP  current_timestamp;", []TokenValuePair{{CURRENT_TIMESTAMP, "CURRENT_TIMESTAMP"}, {CURRENT_TIMESTAMP, "current_timestamp"}, {SEMI, ""}}},
		{"DATABASE  database;", []TokenValuePair{{DATABASE, "DATABASE"}, {DATABASE, "database"}, {SEMI, ""}}},
		{"DEFAULT  default;", []TokenValuePair{{DEFAULT, "DEFAULT"}, {DEFAULT, "default"}, {SEMI, ""}}},
		{"DEFERRABLE  deferrable;", []TokenValuePair{{DEFERRABLE, "DEFERRABLE"}, {DEFERRABLE, "deferrable"}, {SEMI, ""}}},
		{"DEFERRED  deferred;", []TokenValuePair{{DEFERRED, "DEFERRED"}, {DEFERRED, "deferred"}, {SEMI, ""}}},
		{"DELETE  delete;", []TokenValuePair{{DELETE, "DELETE"}, {DELETE, "delete"}, {SEMI, ""}}},
		{"DESC  desc;", []TokenValuePair{{DESC, "DESC"}, {DESC, "desc"}, {SEMI, ""}}},
		{"DETACH  detach;", []TokenValuePair{{DETACH, "DETACH"}, {DETACH, "detach"}, {SEMI, ""}}},
		{"DISTINCT  distinct;", []TokenValuePair{{DISTINCT, "DISTINCT"}, {DISTINCT, "distinct"}, {SEMI, ""}}},
		{"DO  do;", []TokenValuePair{{DO, "DO"}, {DO, "do"}, {SEMI, ""}}},
		{"DROP  drop;", []TokenValuePair{{DROP, "DROP"}, {DROP, "drop"}, {SEMI, ""}}},
		{"EACH  each;", []TokenValuePair{{EACH, "EACH"}, {EACH, "each"}, {SEMI, ""}}},
		{"ELSE  else;", []TokenValuePair{{ELSE, "ELSE"}, {ELSE, "else"}, {SEMI, ""}}},
		{"END  end;", []TokenValuePair{{END, "END"}, {END, "end"}, {SEMI, ""}}},
		{"ESCAPE  escape;", []TokenValuePair{{ESCAPE, "ESCAPE"}, {ESCAPE, "escape"}, {SEMI, ""}}},
		{"EXCEPT  except;", []TokenValuePair{{EXCEPT, "EXCEPT"}, {EXCEPT, "except"}, {SEMI, ""}}},
		{"EXCLUDE  exclude;", []TokenValuePair{{EXCLUDE, "EXCLUDE"}, {EXCLUDE, "exclude"}, {SEMI, ""}}},
		{"EXCLUSIVE  exclusive;", []TokenValuePair{{EXCLUSIVE, "EXCLUSIVE"}, {EXCLUSIVE, "exclusive"}, {SEMI, ""}}},
		{"EXISTS  exists;", []TokenValuePair{{EXISTS, "EXISTS"}, {EXISTS, "exists"}, {SEMI, ""}}},
		{"EXPLAIN  explain;", []TokenValuePair{{EXPLAIN, "EXPLAIN"}, {EXPLAIN, "explain"}, {SEMI, ""}}},
		{"FAIL  fail;", []TokenValuePair{{FAIL, "FAIL"}, {FAIL, "fail"}, {SEMI, ""}}},
		{"FILTER  filter;", []TokenValuePair{{FILTER, "FILTER"}, {FILTER, "filter"}, {SEMI, ""}}},
		{"FIRST  first;", []TokenValuePair{{FIRST, "FIRST"}, {FIRST, "first"}, {SEMI, ""}}},
		{"FOLLOWING  following;", []TokenValuePair{{FOLLOWING, "FOLLOWING"}, {FOLLOWING, "following"}, {SEMI, ""}}},
		{"FOR  for;", []TokenValuePair{{FOR, "FOR"}, {FOR, "for"}, {SEMI, ""}}},
		{"FOREIGN  foreign;", []TokenValuePair{{FOREIGN, "FOREIGN"}, {FOREIGN, "foreign"}, {SEMI, ""}}},
		{"FROM  from;", []TokenValuePair{{FROM, "FROM"}, {FROM, "from"}, {SEMI, ""}}},
		{"FULL  full;", []TokenValuePair{{FULL, "FULL"}, {FULL, "full"}, {SEMI, ""}}},
		{"GENERATED  generated;", []TokenValuePair{{GENERATED, "GENERATED"}, {GENERATED, "generated"}, {SEMI, ""}}},
		{"GLOB  glob;", []TokenValuePair{{GLOB, "GLOB"}, {GLOB, "glob"}, {SEMI, ""}}},
		{"GROUP  group;", []TokenValuePair{{GROUP, "GROUP"}, {GROUP, "group"}, {SEMI, ""}}},
		{"GROUPS  groups;", []TokenValuePair{{GROUPS, "GROUPS"}, {GROUPS, "groups"}, {SEMI, ""}}},
		{"HAVING  having;", []TokenValuePair{{HAVING, "HAVING"}, {HAVING, "having"}, {SEMI, ""}}},
		{"IF  if;", []TokenValuePair{{IF, "IF"}, {IF, "if"}, {SEMI, ""}}},
		{"IGNORE  ignore;", []TokenValuePair{{IGNORE, "IGNORE"}, {IGNORE, "ignore"}, {SEMI, ""}}},
		{"IMMEDIATE  immediate;", []TokenValuePair{{IMMEDIATE, "IMMEDIATE"}, {IMMEDIATE, "immediate"}, {SEMI, ""}}},
		{"IN  in;", []TokenValuePair{{IN, "IN"}, {IN, "in"}, {SEMI, ""}}},
		{"INDEX  index;", []TokenValuePair{{INDEX, "INDEX"}, {INDEX, "index"}, {SEMI, ""}}},
		{"INDEXED  indexed;", []TokenValuePair{{INDEXED, "INDEXED"}, {INDEXED, "indexed"}, {SEMI, ""}}},
		{"INITIALLY  initially;", []TokenValuePair{{INITIALLY, "INITIALLY"}, {INITIALLY, "initially"}, {SEMI, ""}}},
		{"INNER  inner;", []TokenValuePair{{INNER, "INNER"}, {INNER, "inner"}, {SEMI, ""}}},
		{"INSERT  insert;", []TokenValuePair{{INSERT, "INSERT"}, {INSERT, "insert"}, {SEMI, ""}}},
		{"INSTEAD  instead;", []TokenValuePair{{INSTEAD, "INSTEAD"}, {INSTEAD, "instead"}, {SEMI, ""}}},
		{"INTERSECT  intersect;", []TokenValuePair{{INTERSECT, "INTERSECT"}, {INTERSECT, "intersect"}, {SEMI, ""}}},
		{"INTO  into;", []TokenValuePair{{INTO, "INTO"}, {INTO, "into"}, {SEMI, ""}}},
		{"IS  is;", []TokenValuePair{{IS, "IS"}, {IS, "is"}, {SEMI, ""}}},
		{"ISNULL  isnull;", []TokenValuePair{{ISNULL, "ISNULL"}, {ISNULL, "isnull"}, {SEMI, ""}}},
		{"JOIN  join;", []TokenValuePair{{JOIN, "JOIN"}, {JOIN, "join"}, {SEMI, ""}}},
		{"KEY  key;", []TokenValuePair{{KEY, "KEY"}, {KEY, "key"}, {SEMI, ""}}},
		{"LAST  last;", []TokenValuePair{{LAST, "LAST"}, {LAST, "last"}, {SEMI, ""}}},
		{"LEFT  left;", []TokenValuePair{{LEFT, "LEFT"}, {LEFT, "left"}, {SEMI, ""}}},
		{"LIKE  like;", []TokenValuePair{{LIKE, "LIKE"}, {LIKE, "like"}, {SEMI, ""}}},
		{"LIMIT  limit;", []TokenValuePair{{LIMIT, "LIMIT"}, {LIMIT, "limit"}, {SEMI, ""}}},
		{"MATCH  match;", []TokenValuePair{{MATCH, "MATCH"}, {MATCH, "match"}, {SEMI, ""}}},
		{"NATURAL  natural;", []TokenValuePair{{NATURAL, "NATURAL"}, {NATURAL, "natural"}, {SEMI, ""}}},
		{"NO  no;", []TokenValuePair{{NO, "NO"}, {NO, "no"}, {SEMI, ""}}},
		{"NOT  not;", []TokenValuePair{{NOT, "NOT"}, {NOT, "not"}, {SEMI, ""}}},
		{"NOTHING  nothing;", []TokenValuePair{{NOTHING, "NOTHING"}, {NOTHING, "nothing"}, {SEMI, ""}}},
		{"NOTNULL  notnull;", []TokenValuePair{{NOTNULL, "NOTNULL"}, {NOTNULL, "notnull"}, {SEMI, ""}}},
		{"NULL  null;", []TokenValuePair{{NULL, "NULL"}, {NULL, "null"}, {SEMI, ""}}},
		{"NULLS  nulls;", []TokenValuePair{{NULLS, "NULLS"}, {NULLS, "nulls"}, {SEMI, ""}}},
		{"OF  of;", []TokenValuePair{{OF, "OF"}, {OF, "of"}, {SEMI, ""}}},
		{"OFFSET  offset;", []TokenValuePair{{OFFSET, "OFFSET"}, {OFFSET, "offset"}, {SEMI, ""}}},
		{"ON  on;", []TokenValuePair{{ON, "ON"}, {ON, "on"}, {SEMI, ""}}},
		{"OR  or;", []TokenValuePair{{OR, "OR"}, {OR, "or"}, {SEMI, ""}}},
		{"ORDER  order;", []TokenValuePair{{ORDER, "ORDER"}, {ORDER, "order"}, {SEMI, ""}}},
		{"OTHERS  others;", []TokenValuePair{{OTHERS, "OTHERS"}, {OTHERS, "others"}, {SEMI, ""}}},
		{"OUTER  outer;", []TokenValuePair{{OUTER, "OUTER"}, {OUTER, "outer"}, {SEMI, ""}}},
		{"OVER  over;", []TokenValuePair{{OVER, "OVER"}, {OVER, "over"}, {SEMI, ""}}},
		{"PARTITION  partition;", []TokenValuePair{{PARTITION, "PARTITION"}, {PARTITION, "partition"}, {SEMI, ""}}},
		{"PLAN  plan;", []TokenValuePair{{PLAN, "PLAN"}, {PLAN, "plan"}, {SEMI, ""}}},
		{"PRAGMA  pragma;", []TokenValuePair{{PRAGMA, "PRAGMA"}, {PRAGMA, "pragma"}, {SEMI, ""}}},
		{"PRECEDING  preceding;", []TokenValuePair{{PRECEDING, "PRECEDING"}, {PRECEDING, "preceding"}, {SEMI, ""}}},
		{"PRIMARY  primary;", []TokenValuePair{{PRIMARY, "PRIMARY"}, {PRIMARY, "primary"}, {SEMI, ""}}},
		{"QUERY  query;", []TokenValuePair{{QUERY, "QUERY"}, {QUERY, "query"}, {SEMI, ""}}},
		{"RAISE  raise;", []TokenValuePair{{RAISE, "RAISE"}, {RAISE, "raise"}, {SEMI, ""}}},
		{"RANGE  range;", []TokenValuePair{{RANGE, "RANGE"}, {RANGE, "range"}, {SEMI, ""}}},
		{"RECURSIVE  recursive;", []TokenValuePair{{RECURSIVE, "RECURSIVE"}, {RECURSIVE, "recursive"}, {SEMI, ""}}},
		{"REFERENCES  references;", []TokenValuePair{{REFERENCES, "REFERENCES"}, {REFERENCES, "references"}, {SEMI, ""}}},
		{"REGEXP  regexp;", []TokenValuePair{{REGEXP, "REGEXP"}, {REGEXP, "regexp"}, {SEMI, ""}}},
		{"REINDEX  reindex;", []TokenValuePair{{REINDEX, "REINDEX"}, {REINDEX, "reindex"}, {SEMI, ""}}},
		{"RELEASE  release;", []TokenValuePair{{RELEASE, "RELEASE"}, {RELEASE, "release"}, {SEMI, ""}}},
		{"RENAME  rename;", []TokenValuePair{{RENAME, "RENAME"}, {RENAME, "rename"}, {SEMI, ""}}},
		{"REPLACE  replace;", []TokenValuePair{{REPLACE, "REPLACE"}, {REPLACE, "replace"}, {SEMI, ""}}},
		{"RESTRICT  restrict;", []TokenValuePair{{RESTRICT, "RESTRICT"}, {RESTRICT, "restrict"}, {SEMI, ""}}},
		{"RIGHT  right;", []TokenValuePair{{RIGHT, "RIGHT"}, {RIGHT, "right"}, {SEMI, ""}}},
		{"ROLLBACK  rollback;", []TokenValuePair{{ROLLBACK, "ROLLBACK"}, {ROLLBACK, "rollback"}, {SEMI, ""}}},
		{"ROW  row;", []TokenValuePair{{ROW, "ROW"}, {ROW, "row"}, {SEMI, ""}}},
		{"ROWS  rows;", []TokenValuePair{{ROWS, "ROWS"}, {ROWS, "rows"}, {SEMI, ""}}},
		{"SAVEPOINT  savepoint;", []TokenValuePair{{SAVEPOINT, "SAVEPOINT"}, {SAVEPOINT, "savepoint"}, {SEMI, ""}}},
		{"SELECT  select;", []TokenValuePair{{SELECT, "SELECT"}, {SELECT, "select"}, {SEMI, ""}}},
		{"SET  set;", []TokenValuePair{{SET, "SET"}, {SET, "set"}, {SEMI, ""}}},
		{"TABLE  table;", []TokenValuePair{{TABLE, "TABLE"}, {TABLE, "table"}, {SEMI, ""}}},
		{"TEMP  temp;", []TokenValuePair{{TEMP, "TEMP"}, {TEMP, "temp"}, {SEMI, ""}}},
		{"TEMPORARY  temporary;", []TokenValuePair{{TEMPORARY, "TEMPORARY"}, {TEMPORARY, "temporary"}, {SEMI, ""}}},
		{"THEN  then;", []TokenValuePair{{THEN, "THEN"}, {THEN, "then"}, {SEMI, ""}}},
		{"TIES  ties;", []TokenValuePair{{TIES, "TIES"}, {TIES, "ties"}, {SEMI, ""}}},
		{"TO  to;", []TokenValuePair{{TO, "TO"}, {TO, "to"}, {SEMI, ""}}},
		{"TRANSACTION  transaction;", []TokenValuePair{{TRANSACTION, "TRANSACTION"}, {TRANSACTION, "transaction"}, {SEMI, ""}}},
		{"TRIGGER  trigger;", []TokenValuePair{{TRIGGER, "TRIGGER"}, {TRIGGER, "trigger"}, {SEMI, ""}}},
		{"UNBOUNDED  unbounded;", []TokenValuePair{{UNBOUNDED, "UNBOUNDED"}, {UNBOUNDED, "unbounded"}, {SEMI, ""}}},
		{"UNION  union;", []TokenValuePair{{UNION, "UNION"}, {UNION, "union"}, {SEMI, ""}}},
		{"UNIQUE  unique;", []TokenValuePair{{UNIQUE, "UNIQUE"}, {UNIQUE, "unique"}, {SEMI, ""}}},
		{"UPDATE  update;", []TokenValuePair{{UPDATE, "UPDATE"}, {UPDATE, "update"}, {SEMI, ""}}},
		{"USING  using;", []TokenValuePair{{USING, "USING"}, {USING, "using"}, {SEMI, ""}}},
		{"VACUUM  vacuum;", []TokenValuePair{{VACUUM, "VACUUM"}, {VACUUM, "vacuum"}, {SEMI, ""}}},
		{"VALUES  values;", []TokenValuePair{{VALUES, "VALUES"}, {VALUES, "values"}, {SEMI, ""}}},
		{"VIEW  view;", []TokenValuePair{{VIEW, "VIEW"}, {VIEW, "view"}, {SEMI, ""}}},
		{"VIRTUAL  virtual;", []TokenValuePair{{VIRTUAL, "VIRTUAL"}, {VIRTUAL, "virtual"}, {SEMI, ""}}},
		{"WHEN  when;", []TokenValuePair{{WHEN, "WHEN"}, {WHEN, "when"}, {SEMI, ""}}},
		{"WHERE  where;", []TokenValuePair{{WHERE, "WHERE"}, {WHERE, "where"}, {SEMI, ""}}},
		{"WINDOW  window;", []TokenValuePair{{WINDOW, "WINDOW"}, {WINDOW, "window"}, {SEMI, ""}}},
		{"WITH  with;", []TokenValuePair{{WITH, "WITH"}, {WITH, "with"}, {SEMI, ""}}},
		{"WITHOUT  without;", []TokenValuePair{{WITHOUT, "WITHOUT"}, {WITHOUT, "without"}, {SEMI, ""}}},
	}

	for _, _case := range cases {
		scanner := NewScanner([]byte(_case.statement))
		tokens := scanAll(scanner)
		checkEquality(t, tokens, _case.expected)
	}
}

func TestSymbols(t *testing.T) {
	var cases = []TestCase{
		{"-;", []TokenValuePair{{MINUS, ""}, {SEMI, ""}}},
		{"(;", []TokenValuePair{{LP, ""}, {SEMI, ""}}},
		{");", []TokenValuePair{{RP, ""}, {SEMI, ""}}},
		{";;", []TokenValuePair{{SEMI, ""}, {SEMI, ""}}},
		{"+;", []TokenValuePair{{PLUS, ""}, {SEMI, ""}}},
		{"*;", []TokenValuePair{{STAR, ""}, {SEMI, ""}}},
		{"/;", []TokenValuePair{{SLASH, ""}, {SEMI, ""}}},
		{"%;", []TokenValuePair{{REM, ""}, {SEMI, ""}}},
		{"=;", []TokenValuePair{{EQ, ""}, {SEMI, ""}}},
		{"<=;", []TokenValuePair{{LE, ""}, {SEMI, ""}}},
		{"<<;", []TokenValuePair{{LSHIFT, ""}, {SEMI, ""}}},
		{"<;", []TokenValuePair{{LT, ""}, {SEMI, ""}}},
		{">=;", []TokenValuePair{{GE, ""}, {SEMI, ""}}},
		{">>;", []TokenValuePair{{RSHIFT, ""}, {SEMI, ""}}},
		{">;", []TokenValuePair{{GT, ""}, {SEMI, ""}}},
		{"!=;", []TokenValuePair{{NE, ""}, {SEMI, ""}}},
		{",;", []TokenValuePair{{COMMA, ""}, {SEMI, ""}}},
		{"&;", []TokenValuePair{{BITAND, ""}, {SEMI, ""}}},
		{"|;", []TokenValuePair{{BITOR, ""}, {SEMI, ""}}},
		{"!;", []TokenValuePair{{BITNOT, ""}, {SEMI, ""}}},
		{"||;", []TokenValuePair{{CONCAT, ""}, {SEMI, ""}}},
		{".;", []TokenValuePair{{DOT, ""}, {SEMI, ""}}},
	}

	for _, _case := range cases {
		scanner := NewScanner([]byte(_case.statement))
		tokens := scanAll(scanner)
		checkEquality(t, tokens, _case.expected)
	}
}
