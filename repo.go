package golocks

import (
	"database/sql"
	"github.com/go-redis/redis"
	"time"
)

type TryLocker interface {
	TryLock() error
	Unlock() error
}

type LockFactory interface {
	NewLock(key string) TryLocker
}

type ExpiryLockFactory interface {
	NewLock(key string, expiry time.Duration) TryLocker
}

type mapLockFactory struct{}

func NewMapLockFactory() *mapLockFactory {
	return &mapLockFactory{}
}

func (*mapLockFactory) NewLock(key string) TryLocker {
	return NewMapLock(key)
}

type mysqlLockFactory struct{}

func NewMysqlLockFactory(db *sql.DB, tableName string, clearExpiryInterval time.Duration) *mysqlLockFactory {
	InitMysqlLock(db, tableName, clearExpiryInterval)
	return &mysqlLockFactory{}
}

func (*mysqlLockFactory) NewLock(key string, expiry time.Duration) TryLocker {
	return NewMysqlLock(key, expiry)
}

type redisLockFactory struct{}

func NewRedisLockFactory(client *redis.Client) *redisLockFactory {
	InitRedisLock(client)
	return &redisLockFactory{}
}

func (*redisLockFactory) NewLock(key string, expiry time.Duration) TryLocker {
	return NewRedisLock(key, expiry)
}
