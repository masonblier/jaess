package jsparser

import (
	"fmt"
	"reflect"
	"testing"
)

// scanner errors
type SyntaxError struct {
	Message  string
	Location Cursor
}

func NewSyntaxError(message string) *SyntaxError {
	err := new(SyntaxError)
	err.Message = message
	return err
}

func (self *SyntaxError) Error() string {
	return fmt.Sprintf("%s at %+v", self.Message, self.Location)
}

// testing error checking
type EemtoTest struct {
	t *testing.T
}

func (self *EemtoTest) AssertNoError(err error) bool {
	if err != nil {
		self.t.Error(err)
		return false
	}
	return true
}

func (self *EemtoTest) Assert(condition bool, message string, args ...interface{}) bool {
	if !condition {
		self.t.Errorf(message, args...)
		return false
	}
	return true
}

func (self *EemtoTest) AssertEqual(expected interface{}, actual interface{}) bool {
	if !reflect.DeepEqual(expected, actual) {
		self.t.Errorf("expected %v, actual %v", expected, actual)
		return false
	}
	return true
}
