package golocks

import (
	"fmt"
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

func NewRedisLock(name string, expiry time.Duration, spinTries int, spinInterval time.Duration) Locker {
	return &redisLock{
		name:         name,
		expiry:       expiry,
		spinTries:    spinTries,
		spinInterval: spinInterval,
		isOwner:      false,
	}
}

type redisLock struct {
	name         string
	expiry       time.Duration
	spinTries    int
	spinInterval time.Duration
	startAt      time.Time
	isOwner      bool
}

func (l *redisLock) Lock() error {
	for i := 0; i < l.spinTries; i++ {
		if ok, _ := redisClient.SetNX(l.key(), 1, l.expiry).Result(); ok {
			l.startAt = time.Now()
			l.isOwner = true
			return nil
		}

		time.Sleep(l.spinInterval)
	}

	return errorf(fmt.Sprintf("redis lock: lock %s failed after %f seconds", l.key(), float64(l.spinTries)*l.spinInterval.Seconds()))
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
