package locks

import (
	"github.com/bmizerany/assert"
	"testing"
)

func TestKeyLock_TryLock(t *testing.T) {
	lockKey := "test"

	// lock
	l1 := NewKeyLock(lockKey)
	unlockKey, ok := l1.TryLock()
	assert.Equal(t, true, ok)

	// lock again
	l2 := NewKeyLock(lockKey)
	_, ok = l2.TryLock()
	assert.Equal(t, false, ok)

	// lock another
	lockKey3 := "test3"
	l3 := NewKeyLock(lockKey3)
	_, ok = l3.TryLock()
	assert.Equal(t, true, ok)

	// unlock
	err := l2.Unlock("wrong key")
	assert.Equal(t, UnlockKeyInvalidErr, err)
	err = l2.Unlock(unlockKey)
	assert.Equal(t, nil, err)

	// lock and again
	l4 := NewKeyLock(lockKey)
	_, ok = l4.TryLock()
	assert.Equal(t, true, ok)
}
