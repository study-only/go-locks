package golocks

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestKeyLock_TryLock(t *testing.T) {
	lockAKey := "a"
	lockBKey := "b"
	spinInterval := 10 * time.Millisecond

	la1 := NewMapLock(lockAKey, spinInterval)
	la2 := NewMapLock(lockAKey, spinInterval)
	la3 := NewMapLock(lockAKey, spinInterval)
	lb := NewMapLock(lockBKey, spinInterval)

	// lock
	err := la1.Lock()
	assert.Nil(t, err)

	// lock another
	err = lb.Lock()
	assert.Nil(t, err)

	// unlock
	err = la2.Unlock()
	assert.NotNil(t, err)
	err = la1.Unlock()
	assert.Nil(t, err)

	// lock and again
	err = la3.Lock()
	assert.Nil(t, err)
	la3.Unlock()
}

func TestKeyLock_Concurrent(t *testing.T) {
	testValue := 1
	lockKey := "test"
	spinInterval := 10 * time.Millisecond

	l1 := NewMapLock(lockKey, spinInterval)
	l1.Lock()

	go func() {
		l2 := NewMapLock(lockKey, spinInterval)
		l2.Lock()
		defer l2.Unlock()
		assert.Equal(t, 2, testValue)
	}()

	time.Sleep(3 * spinInterval)
	assert.Equal(t, 1, testValue)
	testValue++
	l1.Unlock()

	time.Sleep(3 * spinInterval)
}
