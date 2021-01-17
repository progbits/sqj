package main

import (
	"fmt"
	"io"
	"math"
	"strconv"
)

// Allowed JSON values.
//
// RFC 7159 - Section 2.
type JSONValue int

const (
	JSON_VALUE_OBJECT JSONValue = iota
	JSON_VALUE_ARRAY
	JSON_VALUE_NUMBER
	JSON_VALUE_STRING
	JSON_VALUE_NULL
	JSON_VALUE_TRUE
	JSON_VALUE_FALSE
)

type ASTNode struct {
	// The value type of this node.
	value JSONValue

	// Name of object member.
	name string

	// Object members.
	members []*ASTNode

	// Array values.
	values []*ASTNode

	// Value for tokens of type NUMBER.
	number float64

	// Value for tokens of type STRING.
	string string
}

// Compare two ASTs for equality.
func equal(a, b *ASTNode) bool {
	if a == nil && b == nil {
		return true
	}

	if a == nil || b == nil {
		return false
	}

	if a.name != b.name || a.value != b.value {
		return false
	}

	switch a.value {
	case JSON_VALUE_OBJECT:
		if len(a.members) != len(b.members) {
			return false
		}

		for i := 0; i < len(a.members); i++ {
			if !equal(a.members[i], b.members[i]) {
				return false
			}
		}
	case JSON_VALUE_ARRAY:
		if len(a.values) != len(b.values) {
			return false
		}

		for i := 0; i < len(a.values); i++ {
			if !equal(a.values[i], b.values[i]) {
				return false
			}
		}
	case JSON_VALUE_NUMBER:
		if a.number != b.number {
			return false
		}
	case JSON_VALUE_STRING:
		if a.string != b.string {
			return false
		}
	case JSON_VALUE_NULL, JSON_VALUE_TRUE, JSON_VALUE_FALSE:
		if a.value != b.value {
			return false
		}
	default:
		panic("unexpected value\n")
	}
	return true
}

// Find an AST node by name.
func findNode(ast *ASTNode, name string) *ASTNode {
	var result *ASTNode
	findNodeImpl(ast, &result, name, "")
	return result
}

func findNodeImpl(ast *ASTNode, result **ASTNode, target, prefix string) bool {
	if ast == nil {
		return false
	}

	columnName := concatPrefixName(prefix, ast.name)
	if columnName == target {
		*result = ast
		return true
	}

	if ast.value == JSON_VALUE_OBJECT {
		for i := 0; i < len(ast.members); i++ {
			found := findNodeImpl(ast.members[i], result, target, columnName)
			if found {
				return true
			}
		}
	}
	return false
}

func prettyPrintImpl(writer io.Writer, ast *ASTNode, compact bool, depth int) {
	lineTerm := "\n"
	valueSep := "  "
	if compact {
		lineTerm = ""
		valueSep = ""
	}

	// Indent to the current depth.
	if !compact {
		for i := 0; i < depth; i++ {
			_, _ = fmt.Fprintf(writer, "%s", valueSep)
		}
	}

	literal := ""
	switch ast.value {
	case JSON_VALUE_OBJECT:
		if ast.name != "" && depth > 0 {
			_, _ = fmt.Fprintf(writer, "\"%s\": {%s", ast.name, lineTerm)
		} else {
			_, _ = fmt.Fprintf(writer, "{%s", lineTerm)
		}

		for i := 0; i < len(ast.members); i++ {
			prettyPrintImpl(writer, ast.members[i], compact, depth+1)
			if i < len(ast.members)-1 {
				_, _ = fmt.Fprintf(writer, ",%s", lineTerm)
			}
		}
		_, _ = fmt.Fprintf(writer, "%s", lineTerm)

		for i := 0; i < depth; i++ {
			_, _ = fmt.Fprintf(writer, "%s", valueSep)
		}

		if depth == 0 {
			_, _ = fmt.Fprintf(writer, "}%s", lineTerm)
		} else {
			_, _ = fmt.Fprintf(writer, "}%s", "")
		}
		return
	case JSON_VALUE_ARRAY:
		if ast.name != "" && depth > 0 {
			_, _ = fmt.Fprintf(writer, "\"%s\": [%s", ast.name, lineTerm)
		} else {
			_, _ = fmt.Fprintf(writer, "[%s", lineTerm)
		}

		for i := 0; i < len(ast.values); i++ {
			prettyPrintImpl(writer, ast.values[i], compact, depth+1)
			if i < len(ast.values)-1 {
				_, _ = fmt.Fprintf(writer, ",%s", lineTerm)
			}
		}
		_, _ = fmt.Fprintf(writer, "%s", lineTerm)

		for i := 0; i < depth; i++ {
			_, _ = fmt.Fprintf(writer, "%s", valueSep)
		}

		if depth == 0 {
			_, _ = fmt.Fprintf(writer, "]%s", lineTerm)
		} else {
			_, _ = fmt.Fprintf(writer, "]")
		}

		return
	case JSON_VALUE_NUMBER:
		if ast.name != "" && depth > 0 {
			_, _ = fmt.Fprintf(writer, "\"%s\": ", ast.name)
		}

		// This is not very pretty...
		//
		// Try and work ioOut what precision will recover the original number
		// exactly and if we can't recover it, resort to %1.17g.
		candidates := []string{"%1.15g", "%1.16g", "%1.17g"}
		for i := 0; i < 3; i++ {
			buffer := fmt.Sprintf(candidates[i], ast.number)
			rcvd, _ := strconv.ParseFloat(buffer, 64)
			diff := math.Abs(rcvd - ast.number)
			if diff == 0.0 {
				_, _ = fmt.Fprintf(writer, "%s", buffer)
				return
			}
		}
		_, _ = fmt.Fprintf(writer, "%1.17g", ast.number)
		return
	case JSON_VALUE_STRING:
		if ast.name != "" && depth > 0 {
			_, _ = fmt.Fprintf(writer, "\"%s\": \"%s\"", ast.name, ast.string)
		} else {
			_, _ = fmt.Fprintf(writer, "\"%s\"", ast.string)
		}
		return
	case JSON_VALUE_NULL:
		literal = "null"
	case JSON_VALUE_TRUE:
		literal = "true"
	case JSON_VALUE_FALSE:
		literal = "false"
	default:
		panic("unexpected value\n")
	}

	// Handle literal values.
	if ast.name != "" && depth > 0 {
		_, _ = fmt.Fprintf(writer, "\"%s\": %s", ast.name, literal)
	} else {
		_, _ = fmt.Fprintf(writer, "%s", literal)
	}
}

// prettyPrint pretty prints an AST.
func prettyPrint(writer io.Writer, ast *ASTNode, compact bool) {
	prettyPrintImpl(writer, ast, compact, 0)
}
