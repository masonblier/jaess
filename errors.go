package jsparser

import (
	"fmt"
)

// input stream errors
type SyntaxError struct {
	Message  string
	Location Cursor
}

func (self SyntaxError) Error() string {
	return fmt.Sprintf("%s at %+v", self.Message, self.Location)
}

// scanner errors
type ScannerError struct {
	Message string
}

func (self ScannerError) Error() string {
	return self.Message
}

// parser errors
type ParseError struct {
	Message  string
	Location Cursor
}

func NewParseError(message string, args ...interface{}) *ParseError {
  err := new(ParseError)
  fmsg := fmt.Sprintf(message, args...)
  err.Message = fmt.Sprintf("\x1b[91m= %s\x1b[0m", fmsg)
  return err
}

func (self ParseError) SetLocation(location Cursor) ParseError {
  self.Location = location
  return self
}

func (self ParseError) Error() string {
	return fmt.Sprintf("%s at %+v", self.Message, self.Location)
}
