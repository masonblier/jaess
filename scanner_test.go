package jsparser

import (
	// "fmt"
	"bufio"
	"os"
	"strings"
	"testing"
)

func TestBasicScanning(raw_t *testing.T) {
	t := EemtoTest{raw_t}

	test_source := "anatøm + 1.20"
	// fmt.Printf("\x1b[96m-- lexing ----\n%s\n--------------\x1b[0m\n", test_source)
	inputReader := strings.NewReader(test_source)
	scanner := NewTokenScanner(inputReader)

	var token *Token
	var err error

	token, err = scanner.Next()
	t.AssertNoError(err)
	t.AssertEqual(Token{TOKEN_ATOM, Cursor{0, 0}, "anatøm"}, *token)
	// fmt.Printf("\x1b[90m%+v\x1b[0m\n", token)

	token, err = scanner.Next()
	t.AssertNoError(err)
	t.AssertEqual(Token{TOKEN_OPERATOR, Cursor{0, 7}, "+"}, *token)
	// fmt.Printf("\x1b[90m%+v\x1b[0m\n", token)

	token, err = scanner.Next()
	t.AssertNoError(err)
	t.AssertEqual(Token{TOKEN_NUMBER, Cursor{0, 9}, "1.20"}, *token)
	// fmt.Printf("\x1b[90m%+v\x1b[0m\n", token)
}

func TestMultilineWithComments(raw_t *testing.T) {
	t := EemtoTest{raw_t}

	test_source := "can + /* it \nhandle */ // maybe\n{ \"this\" } \nasdf"
	tokens := make([]Token, 10)
	tokens[0] = Token{TOKEN_ATOM, Cursor{0, 0}, "can"}
	tokens[1] = Token{TOKEN_OPERATOR, Cursor{0, 4}, "+"}
	tokens[2] = Token{TOKEN_COMMENT, Cursor{0, 6}, "/* it \nhandle */"}
	tokens[3] = Token{TOKEN_COMMENT, Cursor{1, 10}, "// maybe"}
	tokens[4] = Token{TOKEN_NEWLINE, Cursor{1, 18}, "\n"}
	tokens[5] = Token{TOKEN_DELIMITER, Cursor{2, 0}, "{"}
	tokens[6] = Token{TOKEN_STRING, Cursor{2, 2}, "\"this\""}
	tokens[7] = Token{TOKEN_DELIMITER, Cursor{2, 9}, "}"}
	tokens[8] = Token{TOKEN_NEWLINE, Cursor{2, 11}, "\n"}
	tokens[9] = Token{TOKEN_ATOM, Cursor{3, 0}, "asdf"}

	// fmt.Printf("\x1b[96m-- expected tokens ----\n%v\n--------------\x1b[0m\n", tokens)
	// fmt.Printf("\x1b[96m-- lexing ----\n%s\n--------------\x1b[0m\n", test_source)
	inputReader := strings.NewReader(test_source)
	scanner := NewTokenScanner(inputReader)

	for _, etkn := range tokens {
		token, err := scanner.Next()
		if !(t.AssertNoError(err) &&
			t.Assert(token != nil, "unexpected end of scanner") &&
			t.AssertEqual(etkn, *token)) {
			return
		}

		// fmt.Printf("\x1b[90m%+v\n%+v\x1b[0m\n", etkn, token)
	}

	token, _ := scanner.Next()
	t.Assert(token == nil, "scanner emitting excessive symbols")
}


func TestAgainstShapeObjects(raw_t *testing.T) {
	t := EemtoTest{raw_t}

	source, err := os.Open("fixtures/shape-objects.js")
	if !t.AssertNoError(err) {
		return
	}

	tokens := make([]Token, 3)
	tokens[0] = Token{TOKEN_COMMENT, Cursor{0, 0}, "// @src https://developer.mozilla.org/en-US/docs/Web/JavaScript/Reference/Global_Objects/Object/create"}
	tokens[1] = Token{TOKEN_NEWLINE, Cursor{0, 102}, "\n"}
	// todo fix: don't emit newlines for blank lines
	tokens[2] = Token{TOKEN_NEWLINE, Cursor{1, 0}, "\n"}

	inputReader := bufio.NewReader(source)
	scanner := NewTokenScanner(inputReader)

	for _, etkn := range tokens {
		token, err := scanner.Next()
		if !(t.AssertNoError(err) &&
			t.Assert(token != nil, "unexpected end of scanner") &&
			t.AssertEqual(etkn, *token)) {
			return
		}

		// fmt.Printf("\x1b[90m%+v\x1b[0m\n", token)
		// fmt.Printf("\x1b[90m%+v\n%+v\x1b[0m\n", etkn, token)
	}

	// todo finish test
	// token, _ := scanner.Next()
	// t.Assert(token == nil, "scanner emitting excessive symbols")
}
