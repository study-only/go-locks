package golocks

import (
	"errors"
	"fmt"
)

func errorf(format string, args ...interface{}) error {
	return errors.New(fmt.Sprintf(format, args...))
}
