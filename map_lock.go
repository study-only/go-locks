package golocks

import (
	"sync"
)

var (
	mapLockMu = new(sync.Mutex)
	mapLocked = make(map[string]struct{})
)

func NewMapLock(name string) *mapLock {
	return &mapLock{
		name: name,
	}
}

type mapLock struct {
	name    string
	isOwner bool
}

func (l *mapLock) TryLock() error {
	mapLockMu.Lock()
	defer mapLockMu.Unlock()

	if _, locked := mapLocked[l.name]; locked {
		return errorf("map lock: already locked")
	}

	mapLocked[l.name] = struct{}{}
	l.isOwner = true
	return nil
}

func (l *mapLock) Unlock() error {
	mapLockMu.Lock()
	defer mapLockMu.Unlock()

	if !l.isOwner {
		return errorf("map lock: not owner")
	}

	_, locked := mapLocked[l.name]
	if !locked {
		return errorf("map lock: unlock of unlocked")
	}

	delete(mapLocked, l.name)
	return nil
}
