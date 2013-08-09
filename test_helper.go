package jsparser

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"reflect"
	"strings"
	"testing"
)

// testing error checking
type TestWrapper struct {
	t     *testing.T
	Trace bool
}

// constructs a TestWrapper instance with default options
func NewTestWrapper(raw_t *testing.T) *TestWrapper {
	wrapper := new(TestWrapper)
	wrapper.t = raw_t
	return wrapper
}

// open text file as buffered reader
func (self *TestWrapper) ReadFile(filename string) *bufio.Reader {
	inputFile, err := os.Open(filename)
	stream := bufio.NewReader(inputFile)
	if self.AssertNoError(err) {
		return stream
	} else {
		return nil
	}
}

// assert err != nil, otherwise report err
func (self *TestWrapper) AssertNoError(err error) bool {
	if err != nil {
		self.t.Error(err)
		return false
	}
	return true
}

// assert result of bool expression
func (self *TestWrapper) Assert(condition bool, message string, args ...interface{}) bool {
	if !condition {
		self.t.Errorf(message, args...)
		return false
	}
	return true
}

// assert deep equal
func (self *TestWrapper) AssertEqual(expected interface{}, actual interface{}) bool {
	if !reflect.DeepEqual(expected, actual) {
		self.t.Errorf("expected %v, actual %v", expected, actual)
		return false
	}
	return true
}

// assert buffers equal line-by-line
func (self *TestWrapper) AssertEqualLines(expected *bufio.Reader, actual *bufio.Reader) bool {
	i := 0

	// line by line check, returns out on fail
	for {
		aline, aerr := actual.ReadString('\n')
		if aerr != nil && aerr != io.EOF {
			self.t.Error(aerr)
			return false
		}
		eline, eerr := expected.ReadString('\n')
		if eerr != nil && eerr != io.EOF {
			self.t.Error(eerr)
			return false
		}
		actual_line := strings.TrimRight(aline, "\n")
		expected_line := strings.TrimRight(eline, "\n")

		if expected_line != actual_line {
			fmt.Printf("\x1b[92m=%3d %s\x1b[0m\n", i, expected_line)
			fmt.Printf("\x1b[91m-%3d %s\x1b[0m\n", i, actual_line)
			for j := 1; j < 3; j++ {
				line, err := actual.ReadString('\n')
				if err == nil {
					line = strings.TrimRight(line, "\n")
					fmt.Printf("\x1b[90m-%3d %s\x1b[0m\n", i+j, line)
				}
			}
			self.Assert(false, "see diff")
			return false
		}
		if self.Trace {
			fmt.Printf("\x1b[90m%%%3d %s\x1b[0m\n", i, actual_line)
		}

		i++

		if aerr == io.EOF {
			break
		}
	}

	// last check: expected has more lines
	is_fail := false
	for {
		line, err := expected.ReadString('\n')
		if err == nil {
			is_fail = true
			line = strings.TrimRight(line, "\n")
			fmt.Printf("\x1b[92m=%3d %s\x1b[0m\n", i, line)
		} else {
			break
		}
		i++
	}
	if is_fail {
		self.Assert(false, "see diff")
	}
	return is_fail
}
