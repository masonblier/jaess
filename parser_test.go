package jaess

import (
	"os"
	"fmt"
	"bufio"
	"testing"
)

func TestBasicParse(raw_t *testing.T) {
	t := NewTestWrapper(raw_t)
	// t.Trace = true
	_RunParserTest("basic-parse", t)
}

func TestParseExportConstants(raw_t *testing.T) {
	t := NewTestWrapper(raw_t)
	// t.Trace = true
	_RunParserTest("exported-constants", t)
}

func TestParseNegatives(raw_t *testing.T) {
	t := NewTestWrapper(raw_t)
	// t.Trace = true
	_RunParserTest("negatives", t)
}

func TestParseShapeObjects(raw_t *testing.T) {
	t := NewTestWrapper(raw_t)
	t.Trace = true
	_RunParserTest("shape-objects", t)
}

func _RunParserTest(fixture_name string, t *TestWrapper) {
	test_input, err := os.Open(fmt.Sprintf("fixtures/%s.js", fixture_name))
	test_source := bufio.NewReader(test_input)
	t.AssertNoError(err)

	ast, err := NewParser(test_source).Parse()
	if !t.AssertNoError(err) {
		return
	}

	astBuffer, err := FormattedAstBuffer(ast)
	t.AssertNoError(err)

	expected_ast := t.ReadFile(fmt.Sprintf("fixtures/%s-ast.json", fixture_name))
	actual_ast := bufio.NewReader(astBuffer)

	t.AssertEqualLines(expected_ast, actual_ast)
}
