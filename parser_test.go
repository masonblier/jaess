package jaess

import (
	"os"
	"bufio"
	"testing"
)

func TestBasicParse(raw_t *testing.T) {
	t := NewTestWrapper(raw_t)

	test_source := ";anatøm + 1.20"
	ast, err := Parse(test_source)
	t.AssertNoError(err)

	astBuffer, err := FormattedAstBuffer(ast)
	t.AssertNoError(err)

	expected_ast := t.ReadFile("fixtures/basic-parse-ast.json")
	actual_ast := bufio.NewReader(astBuffer)


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

	astBuffer, err := FormattedAstBuffer(ast)
	t.AssertNoError(err)

	expected_ast := t.ReadFile("fixtures/exported-constants-ast.json")
	actual_ast := bufio.NewReader(astBuffer)

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

	astBuffer, err := FormattedAstBuffer(ast)
	t.AssertNoError(err)

	expected_ast := t.ReadFile("fixtures/shape-objects-ast.json")
	actual_ast := bufio.NewReader(astBuffer)


	t.AssertEqualLines(expected_ast, actual_ast)
}
