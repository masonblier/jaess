package jsparser

import (
	// "fmt"
	"io"
	"strings"
)

// Parser instance, consumes a TokenScanner
type Parser struct {
	scanner *TokenScanner
	state _ParserState
}

// state of parser machine
type _ParserState struct {
	knownKind AstKind
}

// parses a string into an *AstNode{type:Program,...}
func parse(source string) (*AstNode, error) {
	parser := NewParser(strings.NewReader(source))
	return parser.Next()
}

// create a new parser
func NewParser(input io.RuneScanner) *Parser {
	parser := new(Parser)
	parser.scanner = NewTokenScanner(input)
	return parser
}

// token, err := scanner.Next()
// fmt.Printf("\x1b[90m%+v\x1b[0m\n", token)

// gets the next statement
func (self *Parser) Next() (*AstNode, error) {
	var statement *AstNode

	self.state.knownKind = AST_UNKNOWN

	// for {
	// 	token, err := self.scanner

	// 	if err != nil {
	// 		return nil, err
	// 	}
	// 	if token == nil {
	// 		return nil, nil
	// 	}

	// 	switch self.state.knownKind {
	// 	case AST_UNKNOWN:

	// 		switch token.kind {

	// 		}

	// 	}
	// }

	// return NewAstNode(AST_PROGRAM)
	return statement, nil
}
