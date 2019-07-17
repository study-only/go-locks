package golocks

import (
	"github.com/go-redis/redis"
	"time"
)

var redisClient *redis.Client

func InitRedisLock(client *redis.Client) {
	if client == nil {
		panic("client is nil")
	}
	redisClient = client
}

func NewRedisLock(name string, expiry time.Duration) TryLocker {
	return &redisLock{
		name:    name,
		expiry:  expiry,
		isOwner: false,
	}
}

type redisLock struct {
	name    string
	expiry  time.Duration
	startAt time.Time
	isOwner bool
}

func (l *redisLock) TryLock() error {
	if ok, _ := redisClient.SetNX(l.key(), 1, l.expiry).Result(); !ok {
		return errorf("redis lock: %s already locked", l.key())
	}

	l.startAt = time.Now()
	l.isOwner = true
	return nil

}

func (l *redisLock) Unlock() error {
	if !l.isOwner {
		return errorf("redis lock: not owner")
	}
	if time.Now().UnixNano()-l.startAt.UnixNano() >= l.expiry.Nanoseconds() {
		return errorf("redis lock: lock expired")
	}

	if err := redisClient.Del(l.key()).Err(); err != nil {
		return errorf("redis lock: %s", err)
	}

	l.isOwner = false
	return nil
}

func (l redisLock) key() string {
	return "redis_lock:" + l.name
}
