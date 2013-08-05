package jsparser

import (
  "fmt"
  "testing"
)

func TestParse(raw_t *testing.T) {
  t := EemtoTest{raw_t}

  test_source := "anatøm + 1.20"
  // fmt.Printf("\x1b[96m-- lexing ----\n%s\n--------------\x1b[0m\n", test_source)
  ast, err := parse(test_source)
  t.AssertNoError(err)

  fmt.Printf("\x1b[93m-- ast ----\n%+v\n--------------\x1b[0m\n", ast)
}
