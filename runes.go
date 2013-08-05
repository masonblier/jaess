package jsparser

import (
	"unicode"
)

// spaces excluding newlines
func IsInlineWhitespaceRune(r rune) bool {
	return unicode.IsSpace(r) && r != '\n'
}

// delimiter: punct breaks up the token
// each rune emits a seperate token
func IsDelimeterRune(r rune) bool {
	switch r {
	case '[', ']', '(', ')', '{', '}', ';', ',':
		return true
	}
	return false
}

// pretty much punct thats not a delimeter
// contiguous operator runes emit one token
func IsOperatorRune(r rune) bool {
	switch r {
	case '<', '>', '+', '-', '*', '/', '%', '=', '&', '|', '^', '!', '~', '?', ':','.':
		return true
	}
	return false
}

// unicode digits only
func IsDigitRune(r rune) bool {
	if unicode.IsDigit(r) {
		return true
	}
	return false
}

// runes for keywords or identifiers
func IsAtomRune(r rune) bool {
	if r == '_' || r == '$' || unicode.IsLetter(r) {
		return true
	}
	return false
}
