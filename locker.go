package golocks

import (
	"errors"
	"fmt"
)

type Locker interface {
	Lock() error
	Unlock() error
}

func errorf(format string, args ...interface{}) error {
	return errors.New(fmt.Sprintf(format, args...))
}
