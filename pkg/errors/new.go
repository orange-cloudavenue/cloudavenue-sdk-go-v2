package errors

import (
	"fmt"
)

// Newf creates a new error with a formatted message.
func Newf(format string, args ...interface{}) error {
	return fmt.Errorf(format, args...)
}
