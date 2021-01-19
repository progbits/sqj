package json

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
	Value JSONValue

	// Name of object member.
	Name string

	// Object members.
	Members []*ASTNode

	// Array values.
	Values []*ASTNode

	// Value for tokens of type NUMBER.
	Number float64

	// Value for tokens of type STRING.
	String string
}

// Compare two ASTs for equality.
func equal(a, b *ASTNode) bool {
	if a == nil && b == nil {
		return true
	}

	if a == nil || b == nil {
		return false
	}

	if a.Name != b.Name || a.Value != b.Value {
		return false
	}

	switch a.Value {
	case JSON_VALUE_OBJECT:
		if len(a.Members) != len(b.Members) {
			return false
		}

		for i := 0; i < len(a.Members); i++ {
			if !equal(a.Members[i], b.Members[i]) {
				return false
			}
		}
	case JSON_VALUE_ARRAY:
		if len(a.Values) != len(b.Values) {
			return false
		}

		for i := 0; i < len(a.Values); i++ {
			if !equal(a.Values[i], b.Values[i]) {
				return false
			}
		}
	case JSON_VALUE_NUMBER:
		if a.Number != b.Number {
			return false
		}
	case JSON_VALUE_STRING:
		if a.String != b.String {
			return false
		}
	case JSON_VALUE_NULL, JSON_VALUE_TRUE, JSON_VALUE_FALSE:
		if a.Value != b.Value {
			return false
		}
	default:
		panic("unexpected value\n")
	}
	return true
}

// Find an AST node by name.
func FindNode(ast *ASTNode, name string) *ASTNode {
	var result *ASTNode
	findNodeImpl(ast, &result, name, "")
	return result
}

func findNodeImpl(ast *ASTNode, result **ASTNode, target, prefix string) bool {
	if ast == nil {
		return false
	}

	columnName := concatPrefixName(prefix, ast.Name)
	if columnName == target {
		*result = ast
		return true
	}

	if ast.Value == JSON_VALUE_OBJECT {
		for i := 0; i < len(ast.Members); i++ {
			found := findNodeImpl(ast.Members[i], result, target, columnName)
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
	switch ast.Value {
	case JSON_VALUE_OBJECT:
		if ast.Name != "" && depth > 0 {
			_, _ = fmt.Fprintf(writer, "\"%s\": {%s", ast.Name, lineTerm)
		} else {
			_, _ = fmt.Fprintf(writer, "{%s", lineTerm)
		}

		for i := 0; i < len(ast.Members); i++ {
			prettyPrintImpl(writer, ast.Members[i], compact, depth+1)
			if i < len(ast.Members)-1 {
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
		if ast.Name != "" && depth > 0 {
			_, _ = fmt.Fprintf(writer, "\"%s\": [%s", ast.Name, lineTerm)
		} else {
			_, _ = fmt.Fprintf(writer, "[%s", lineTerm)
		}

		for i := 0; i < len(ast.Values); i++ {
			prettyPrintImpl(writer, ast.Values[i], compact, depth+1)
			if i < len(ast.Values)-1 {
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
		if ast.Name != "" && depth > 0 {
			_, _ = fmt.Fprintf(writer, "\"%s\": ", ast.Name)
		}

		// This is not very pretty...
		//
		// Try and work ioOut what precision will recover the original number
		// exactly and if we can't recover it, resort to %1.17g.
		candidates := []string{"%1.15g", "%1.16g", "%1.17g"}
		for i := 0; i < 3; i++ {
			buffer := fmt.Sprintf(candidates[i], ast.Number)
			rcvd, _ := strconv.ParseFloat(buffer, 64)
			diff := math.Abs(rcvd - ast.Number)
			if diff == 0.0 {
				_, _ = fmt.Fprintf(writer, "%s", buffer)
				return
			}
		}
		_, _ = fmt.Fprintf(writer, "%1.17g", ast.Number)
		return
	case JSON_VALUE_STRING:
		if ast.Name != "" && depth > 0 {
			_, _ = fmt.Fprintf(writer, "\"%s\": \"%s\"", ast.Name, ast.String)
		} else {
			_, _ = fmt.Fprintf(writer, "\"%s\"", ast.String)
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
	if ast.Name != "" && depth > 0 {
		_, _ = fmt.Fprintf(writer, "\"%s\": %s", ast.Name, literal)
	} else {
		_, _ = fmt.Fprintf(writer, "%s", literal)
	}
}

// prettyPrint pretty prints an AST.
func PrettyPrint(writer io.Writer, ast *ASTNode, compact bool) {
	prettyPrintImpl(writer, ast, compact, 0)
}
