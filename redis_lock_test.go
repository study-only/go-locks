package golocks

import (
	"fmt"
	"github.com/go-redis/redis"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

var (
	testRedisHost = env("TEST_REDIS_HOST", "127.0.0.1")
	testRedisPORT = env("TEST_REDIS_PORT", "6379")
	testRedisPWD  = env("TEST_REDIS_PWD", "")
)

func TestRedisLock_Lock(t *testing.T) {
	client := getRedisClient(t, testRedisHost, testRedisPORT, testRedisPWD)
	InitRedisLock(client)
	lock := NewRedisLock("test", time.Second, 3, 100*time.Millisecond)

	err := lock.Lock()
	assert.Nil(t, err)
	err = lock.Unlock()
	assert.Nil(t, err)
}

func TestRedisLock_Expired(t *testing.T) {
	client := getRedisClient(t, testRedisHost, testRedisPORT, testRedisPWD)
	InitRedisLock(client)
	lock := NewRedisLock("expiry", 500*time.Millisecond, 3, 100*time.Millisecond)

	err := lock.Lock()
	assert.Nil(t, err)
	time.Sleep(600 * time.Millisecond)

	err = lock.Unlock()
	assert.NotNil(t, err)
	err = lock.Lock()
	assert.Nil(t, err)
}

func TestRedisLock_ConcurrentLock(t *testing.T) {
	spinInterval := 10 * time.Millisecond
	client := getRedisClient(t, testRedisHost, testRedisPORT, testRedisPWD)
	InitRedisLock(client)
	l1 := NewRedisLock("concurrent", time.Second, 10, spinInterval)
	l2 := NewRedisLock("concurrent", time.Second, 10, spinInterval)

	testValue := 1
	l1.Lock()
	go func() {
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

func getRedisClient(t *testing.T, host, port, pwd string) *redis.Client {
	addr := fmt.Sprintf("%s:%s", host, port)
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: pwd,
		DB:       0,
	})

	if err := client.Ping().Err(); err != nil {
		t.Fatal(err)
	}

	client.FlushAll()
	return client
}
