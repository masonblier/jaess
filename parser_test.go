package jsparser

import (
	"bufio"
	"bytes"
	"encoding/json"
	"os"
	"testing"
)

func TestBasicParse(raw_t *testing.T) {
	t := NewTestWrapper(raw_t)

	test_source := "anatøm + 1.20"
	ast, err := parse(test_source)
	t.AssertNoError(err)

	actual_ast := _MarshalAndIndentJsonBufioReader(t, ast)
	expected_ast := t.ReadFile("fixtures/basic-parse-ast.json")

	t.AssertEqualLines(expected_ast, actual_ast)
}

func TestParseExportConstants(raw_t *testing.T) {
	t := NewTestWrapper(raw_t)
	// t.Trace = true

	test_input, err := os.Open("fixtures/exported-constants.js")
	test_source := bufio.NewReader(test_input)
	t.AssertNoError(err)

	ast, err := NewParser(test_source).Parse()
	if !t.AssertNoError(err) {
		return
	}

	actual_ast := _MarshalAndIndentJsonBufioReader(t, ast)
	expected_ast := t.ReadFile("fixtures/exported-constants-ast.json")

	t.AssertEqualLines(expected_ast, actual_ast)
}

func TestParseShapeObjects(raw_t *testing.T) {
	t := NewTestWrapper(raw_t)
	// t.Trace = true

	test_input, err := os.Open("fixtures/shape-objects.js")
	test_source := bufio.NewReader(test_input)
	t.AssertNoError(err)

	ast, err := NewParser(test_source).Parse()
	if !t.AssertNoError(err) {
		return
	}

	actual_ast := _MarshalAndIndentJsonBufioReader(t, ast)
	expected_ast := t.ReadFile("fixtures/shape-objects-ast.json")

	t.AssertEqualLines(expected_ast, actual_ast)
}

func _MarshalAndIndentJsonBufioReader(t *TestWrapper, ast AstNode) *bufio.Reader {
	jsonStr, err := json.Marshal(ast)
	t.AssertNoError(err)

	var dst bytes.Buffer
	err = json.Indent(&dst, jsonStr, "", "    ")
	t.AssertNoError(err)

	return bufio.NewReader(&dst)
}
