package jsparser

import (
	"fmt"
	"io"
)

// a scanner of tokens
type TokenScanner struct {
	input io.RuneScanner
	state _ScannerState
}

// internal state of the scanner
type _ScannerState struct {
	offset   int64
	Location Cursor
}

// location within the source input
type Cursor struct {
	line   int
	column int
}

// creates a new token scanner with default state
func NewTokenScanner(input io.RuneScanner) *TokenScanner {
	return &TokenScanner{input, _ScannerState{0, Cursor{0, 0}}}
}

func (self *_ScannerState) _IncrementCursor(r rune) {
	if r == '\n' {
		self.Location.line += 1
		self.Location.column = 0
	} else {
		self.Location.column += 1
	}
}

// gets the next token available
// returns a token or nil, or nil with an error
func (self *TokenScanner) Next() (*Token, error) {
	var token *Token

	for {
		r, rlen, err := self.input.ReadRune()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}

		if token == nil {
			token = new(Token)
			token.Location = self.state.Location
		}

		if ok, err := token.ConsumeRune(r); ok {
			self.state.offset += int64(rlen)
			self.state._IncrementCursor(r)
		} else {
			self.input.UnreadRune()
			if err != nil {
				err.Location = self.state.Location
				return nil, err
			}
			if _TOKEN_SPACE == token.Kind {
				token = nil
			} else {
				break
			}
		}
	}

	// hack to correct type of single line comment at eof
	if token != nil && _TOKEN_COMMENT_SINGLE_LINE == token.Kind {
		token.Kind = TOKEN_COMMENT
	}

	if token != nil && token.Kind >= _TOKEN_HIDDEN_START {
		if token.Kind == _TOKEN_COMMENT_MULTI_LINE ||
			token.Kind == _TOKEN_COMMENT_MULTI_LINE_MAY_END {
			return nil, &SyntaxError{"incomplete multiline comment", token.Location}
		} else {
			return nil, &SyntaxError{"unexpected eof", self.state.Location}
		}
	}

	return token, nil
}

// reads the run into the token
// returning true if rune is accepted, false if rune is rejected
// error is nil unless a syntax error is detected
func (self *Token) ConsumeRune(r rune) (bool, *SyntaxError) {
	switch self.Kind {

	// stage 1: lexical identification
	case TOKEN_UNKNOWN:
		switch {
		case IsInlineWhitespaceRune(r):
			self.Kind = _TOKEN_SPACE
		case '\n' == r:
			self.Kind = TOKEN_NEWLINE
		// case '\'' == r:
		// 	self.Kind = _TOKEN_STRING_SINGLE_QUOTE
		case '"' == r:
			self.Kind = _TOKEN_STRING_DOUBLE_QUOTE
		case '/' == r:
			self.Kind = _TOKEN_ONE_SLASH
		case IsDelimeterRune(r):
			self.Kind = TOKEN_DELIMITER
		case IsOperatorRune(r):
			self.Kind = TOKEN_OPERATOR
		case IsAtomRune(r):
			self.Kind = TOKEN_ATOM
		case IsDigitRune(r):
			self.Kind = TOKEN_NUMBER
		default:
			return false, NewSyntaxError(fmt.Sprintf("Invalid Rune %c", r))
		}
		self.Value += string(r)
		return true, nil

	// stage 2: token scan
	case _TOKEN_SPACE:
		if IsInlineWhitespaceRune(r) {
			self.Value += string(r)
			return true, nil
		}
	case _TOKEN_ONE_SLASH:
		if r == '/' {
			self.Value += string(r)
			self.Kind = _TOKEN_COMMENT_SINGLE_LINE
			return true, nil
		}
		if r == '*' {
			self.Value += string(r)
			self.Kind = _TOKEN_COMMENT_MULTI_LINE
			return true, nil
		}
	case _TOKEN_COMMENT_SINGLE_LINE:
		if r == '\n' {
			self.Kind = TOKEN_COMMENT
			return false, nil
		} else {
			self.Value += string(r)
			return true, nil
		}
	case _TOKEN_COMMENT_MULTI_LINE:
		self.Value += string(r)
		if r == '*' {
			self.Kind = _TOKEN_COMMENT_MULTI_LINE_MAY_END
		}
		return true, nil
	case _TOKEN_COMMENT_MULTI_LINE_MAY_END:
		switch r {
		case '/':
			self.Kind = TOKEN_COMMENT
		case '*':
			self.Kind = _TOKEN_COMMENT_MULTI_LINE_MAY_END
		default:
			self.Kind = _TOKEN_COMMENT_MULTI_LINE
		}
		self.Value += string(r)
		return true, nil

	case _TOKEN_STRING_DOUBLE_QUOTE:
		self.Value += string(r)
		if r == '"' {
			self.Kind = TOKEN_STRING
		}
		return true, nil

	case TOKEN_COMMENT:
		return false, nil
	case TOKEN_DELIMITER:
		return false, nil
	case TOKEN_STRING:
		return false, nil
	case TOKEN_OPERATOR:
		if IsOperatorRune(r) {
			self.Value += string(r)
			return true, nil
		}
	case TOKEN_ATOM:
		if IsAtomRune(r) || IsDigitRune(r) {
			self.Value += string(r)
			return true, nil
		}
	case TOKEN_NUMBER:
		if r == '.' || IsDigitRune(r) {
			self.Value += string(r)
			return true, nil
		}

	}
	return false, nil
}
