package locks

import (
	"fmt"
	"github.com/go-redis/redis"
	"time"
)

type RedisLock struct {
	client *redis.Client
	name   string
	expiry time.Duration
	tries  int64
	delay  time.Duration

	startAt time.Time
	isOwner bool
}

func NewRedisLock(client *redis.Client, name string, expiry time.Duration, tries int64, delay time.Duration) RedisLock {
	return RedisLock{
		client:  client,
		name:    name,
		expiry:  expiry,
		tries:   tries,
		delay:   delay,
		isOwner: false,
	}
}

func (m *RedisLock) Lock() error {
	for i := int64(0); i < m.tries; i++ {
		if ok, _ := m.client.SetNX(m.key(), 1, m.expiry).Result(); ok {
			m.startAt = time.Now()
			m.isOwner = true
			return nil
		}

		time.Sleep(m.delay)
	}

	return fmt.Errorf("lock %s failed after %f seconds", m.key(), float64(m.tries)*m.delay.Seconds())
}

func (m *RedisLock) Unlock() error {
	if !m.isOwner {
		return fmt.Errorf("no permmision")
	}
	if time.Now().UnixNano()-m.startAt.UnixNano() >= m.expiry.Nanoseconds() {
		return fmt.Errorf("lock %s already expired", m.key())
	}

	if err := m.client.Del(m.key()).Err(); err != nil {
		return err
	}

	m.isOwner = false
	return nil
}

func (m RedisLock) key() string {
	return "redis_lock:" + m.name
}
