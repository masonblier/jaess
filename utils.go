package jaess

import (
  "unicode"
  "regexp"
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
  case '<', '>', '+', '-', '*', '/', '%', '=', '&', '|', '^', '!', '~', '?', ':', '.':
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

// unary operator token
func IsUnaryOperator(token *Token) bool {
  switch token.Value {
  case "-":
    return true
  }
  return false
}

// formats a block of source code to a function body
func _TrimFunctionSource(s string) string {
  btR, err := regexp.Compile(`^\s*\{`)
  if err != nil {
    panic(err)
  }
  if btR.MatchString(s) {
    bR, err := regexp.Compile(`^\s*{\s*?\n?`)
    if err != nil {
      panic(err)
    }
    eR, err := regexp.Compile(`\s*}\s*;?\s*$`)
    if err != nil {
      panic(err)
    }
    return bR.ReplaceAllLiteralString(eR.ReplaceAllLiteralString(s, ""), "")
  }
  wspbR, err := regexp.Compile(`^\s*?\n?`)
  if err != nil {
    panic(err)
  }
  wspeR, err := regexp.Compile(`\s*$`)
  if err != nil {
    panic(err)
  }
  return wspbR.ReplaceAllLiteralString(wspeR.ReplaceAllLiteralString(s, ""), "")
}