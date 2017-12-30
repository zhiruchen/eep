package eep

import (
	"errors"
	"fmt"
)

var (
	UnterminatedStrError = errors.New("Unterminated string")
)

func newTokenError(line int, msg string) error {
	return fmt.Errorf("line: %d, %s", line, msg)
}
