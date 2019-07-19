package golocks

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMapLock_TryLock(t *testing.T) {
	lockAKey := "a"
	lockBKey := "b"

	la1 := NewMapLock(lockAKey)
	la2 := NewMapLock(lockAKey)
	la3 := NewMapLock(lockAKey)
	lb := NewMapLock(lockBKey)

	// lock
	err := la1.TryLock()
	assert.Nil(t, err)

	// lock another
	err = lb.TryLock()
	assert.Nil(t, err)

	// unlock
	err = la2.Unlock()
	assert.NotNil(t, err)
	err = la1.Unlock()
	assert.Nil(t, err)

	// lock and again
	err = la3.TryLock()
	assert.Nil(t, err)
	err = la3.Unlock()
	assert.Nil(t, err)
}
