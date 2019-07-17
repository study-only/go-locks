package golocks

import (
	"sync"
	"time"
)

var (
	mapLockMu = new(sync.Mutex)
	mapLocked = make(map[string]struct{})
)

func NewMapLock(name string, spinInterval time.Duration) Locker {
	return &keyLock{
		name:         name,
		spinInterval: spinInterval,
	}
}

type keyLock struct {
	name         string
	spinInterval time.Duration
	isOwner      bool
}

func (l *keyLock) Lock() error {
	for {
		if ok := l.tryLock(); ok {
			return nil
		}
		time.Sleep(l.spinInterval)
	}
}

func (l *keyLock) Unlock() error {
	mapLockMu.Lock()
	defer mapLockMu.Unlock()

	if !l.isOwner {
		return errorf("key lock: not owner")
	}

	_, locked := mapLocked[l.name]
	if !locked {
		return errorf("key lock: unlock of unlocked keyLock")
	}

	delete(mapLocked, l.name)
	return nil
}

func (l *keyLock) tryLock() (success bool) {
	mapLockMu.Lock()
	defer mapLockMu.Unlock()

	if _, locked := mapLocked[l.name]; locked {
		return false
	}

	mapLocked[l.name] = struct{}{}
	l.isOwner = true
	return true
}
