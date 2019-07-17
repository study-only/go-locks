package golocks

import (
	"errors"
	"fmt"
	"time"
)

type TryLocker interface {
	TryLock() error
	Unlock() error
}

func NewSpinLock(lock TryLocker, spinTries int, spinInterval time.Duration) *spinLock {
	return &spinLock{
		lock:         lock,
		spinTries:    spinTries,
		spinInterval: spinInterval,
	}
}

type spinLock struct {
	lock         TryLocker
	spinTries    int
	spinInterval time.Duration
}

func (l *spinLock) Lock() error {
	for i := 0; i < l.spinTries; i++ {
		if err := l.lock.TryLock(); err == nil {
			return nil
		}

		time.Sleep(l.spinInterval)
	}

	return errorf(fmt.Sprintf("spin lock: failed after %f seconds", float64(l.spinTries)*l.spinInterval.Seconds()))
}

func (l *spinLock) Unlock() error {
	return l.lock.Unlock()
}

func errorf(format string, args ...interface{}) error {
	return errors.New(fmt.Sprintf(format, args...))
}
