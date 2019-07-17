package golocks

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestSpinLock_Lock(t *testing.T) {
	testValue := 1
	lockKey := "test"
	spinTries := 5
	spinInterval := 10 * time.Millisecond

	l1 := NewSpinLock(NewMapLock(lockKey), spinTries, spinInterval)
	l2 := NewSpinLock(NewMapLock(lockKey), spinTries, spinInterval)

	err := l1.Lock()
	assert.Nil(t, err)

	go func() {
		err := l2.Lock()
		assert.Nil(t, err)
		err = l2.Unlock()
		assert.Nil(t, err)
		assert.Equal(t, 2, testValue)
	}()

	time.Sleep(3 * spinInterval)
	assert.Equal(t, 1, testValue)
	testValue++
	err = l1.Unlock()
	assert.Nil(t, err)

	time.Sleep(3 * spinInterval)
}