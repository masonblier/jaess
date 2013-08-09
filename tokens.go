package jsparser

import (
// "io"
)

// scanned token
type Token struct {
	Type     TokenType
	Location Cursor
	Value    string
}

type TokenType int

const (
	TOKEN_UNKNOWN TokenType = iota
	DELIMITER
	OPERATOR
	ATOM
	NUMBER
	STRING
	COMMENT
	NEWLINE

	_HIDDEN
	_SPACE
	_ONE_SLASH
	_COMMENT_SINGLE_LINE
	_COMMENT_MULTI_LINE
	_COMMENT_MULTI_LINE_MAY_END
	_STRING_SINGLE_QUOTE
	_STRING_DOUBLE_QUOTE
)

func (self TokenType) String() string {
	switch self {

	case TOKEN_UNKNOWN:
		return "TOKEN_UNKNOWN"
	case DELIMITER:
		return "DELIMITER"
	case OPERATOR:
		return "OPERATOR"
	case ATOM:
		return "ATOM"
	case NUMBER:
		return "NUMBER"
	case STRING:
		return "STRING"
	case COMMENT:
		return "COMMENT"
	case NEWLINE:
		return "NEWLINE"

	case _SPACE:
		return "_SPACE"
	case _ONE_SLASH:
		return "_ONE_SLASH"
	case _COMMENT_SINGLE_LINE:
		return "_COMMENT_SINGLE_LINE"
	case _COMMENT_MULTI_LINE:
		return "_COMMENT_MULTI_LINE"
	case _COMMENT_MULTI_LINE_MAY_END:
		return "_COMMENT_MULTI_LINE_MAY_END"
	case _STRING_SINGLE_QUOTE:
		return "_STRING_SINGLE_QUOTE"
	case _STRING_DOUBLE_QUOTE:
		return "_STRING_DOUBLE_QUOTE"

	}
	return "<#error: bad value>"
}
