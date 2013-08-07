package jsparser

import (
	"fmt"
	"io"
)

// a scanner of tokens
type TokenScanner struct {
	input     io.RuneScanner
	offset    int64
	Location  Cursor
	lastToken *Token
	unToken   *Token
}

// location within the source input
type Cursor struct {
	line   int
	column int
}

// creates a new token scanner
func NewTokenScanner(input io.RuneScanner) *TokenScanner {
	return &TokenScanner{input, 0, Cursor{0, 0}, nil, nil}
}

func (self *Cursor) _IncrementByRune(r rune) {
	if r == '\n' {
		self.line += 1
		self.column = 0
	} else {
		self.column += 1
	}
}

// gets the next token available
// returns a token or nil, or nil with an error
func (self *TokenScanner) Next() (*Token, error) {
	var token *Token

	// check undo cache
	if self.unToken != nil {
		token = self.unToken
		self.unToken = nil
		return token, nil
	}

	// scan
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
			token.Location = self.Location
		}

		if ok, err := token.ConsumeRune(r); ok {
			self.offset += int64(rlen)
			self.Location._IncrementByRune(r)
		} else {
			self.input.UnreadRune()
			if err != nil {
				err.Location = self.Location
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

	// check for internal enums
	if token != nil && token.Kind >= _TOKEN_HIDDEN {
		if token.Kind == _TOKEN_COMMENT_MULTI_LINE ||
			token.Kind == _TOKEN_COMMENT_MULTI_LINE_MAY_END {
			return nil, &SyntaxError{"incomplete multiline comment", token.Location}
		} else {
			return nil, &SyntaxError{"unexpected eof", self.Location}
		}
	}

	// cache last token
	self.lastToken = token
	return token, nil
}

// moves the scanner back one. cannot go back more than one.
func (self *TokenScanner) UnNext() error {
	if self.unToken != nil {
		return ScannerError{"consecutive UnNext calls unsupported, must call Next between each"}
	}
	self.unToken = self.lastToken
	return nil
}

// peeks at the next value. cannot call UnNext after peeking.
func (self *TokenScanner) Peek() (*Token, error) {
	if self.unToken == nil {
		var err error
		self.unToken, err = self.Next()
		if err != nil {
			return nil, err
		}
	}
	return self.unToken, nil
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
			return false, &SyntaxError{fmt.Sprintf("Invalid Rune %c", r), Cursor{-1, -1}}
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
