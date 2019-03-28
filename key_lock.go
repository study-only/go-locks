package locks

import (
	"errors"
	"fmt"
	rd "math/rand"
	"sync"
	"time"
)

func init() {
	rd.Seed(time.Now().Unix())
}

var keyLocked = make(map[string]string)
var keyLockMu = new(sync.Mutex)
var UnlockKeyInvalidErr = errors.New("unlock key invalid")

func NewKeyLock(key string) *KeyLock {
	return &KeyLock{
		key: key,
	}
}

type KeyLock struct {
	key string
}

func (k *KeyLock) TryLock() (unlockKey string, success bool) {
	keyLockMu.Lock()
	defer keyLockMu.Unlock()
	unlockKey = fmt.Sprintf("unlockKey:%d:%d", time.Now().Unix(), rd.Int())

	if _, locked := keyLocked[k.key]; locked {
		return "", false
	}

	keyLocked[k.key] = unlockKey
	return unlockKey, true
}

func (k *KeyLock) Unlock(unlockKey string) error {
	keyLockMu.Lock()
	defer keyLockMu.Unlock()

	key, locked := keyLocked[k.key]
	if !locked {
		return nil
	}
	if unlockKey != key {
		return UnlockKeyInvalidErr
	}

	delete(keyLocked, k.key)
	return nil
}
