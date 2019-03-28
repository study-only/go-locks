package locks

import (
	"fmt"
	"github.com/go-redis/redis"
	"github.com/stretchr/testify/assert"
	"sync"
	"testing"
	"time"
)

func getRedisClient(host string, port int64, pwd string) *redis.Client {
	addr := fmt.Sprintf("%s:%d", host, port)
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: pwd,
		DB:       0,
	})

	if err := client.Ping().Err(); err == nil {
		client.FlushAll()
		return client
	} else {
		return nil
	}
}

func TestRedisLock_Lock(t *testing.T) {
	client := getRedisClient("127.0.0.1", 6379, "")
	mutex := NewRedisLock(client, "test", time.Second, 3, 100*time.Millisecond)

	assert.Equal(t, nil, mutex.Lock())
	assert.NotEqual(t, nil, mutex.Lock())
	time.Sleep(2000 * time.Millisecond)
	assert.Equal(t, nil, mutex.Lock())
}

func TestRedisLock_Unlock(t *testing.T) {
	client := getRedisClient("127.0.0.1", 6379, "")
	mutex := NewRedisLock(client, "test", time.Second, 3, 100*time.Millisecond)

	assert.Equal(t, nil, mutex.Lock())
	assert.Equal(t, nil, mutex.Unlock())
	assert.NotEqual(t, nil, mutex.Unlock())
}

func TestRedisLock_UnLockExpired(t *testing.T) {
	client := getRedisClient("127.0.0.1", 6379, "")
	mutex := NewRedisLock(client, "test", 500*time.Millisecond, 3, 100*time.Millisecond)

	assert.Equal(t, nil, mutex.Lock())
	time.Sleep(600 * time.Millisecond)
	assert.NotEqual(t, nil, mutex.Unlock())
}

func TestRedisLock_ConcurrentLock(t *testing.T) {
	client := getRedisClient("127.0.0.1", 6379, "")
	mutex1 := NewRedisLock(client, "concurrent", time.Second, 3, 100*time.Millisecond)
	mutex2 := NewRedisLock(client, "concurrent", time.Second, 3, 100*time.Millisecond)
	wg := sync.WaitGroup{}
	wg.Add(2)

	var err1, err2 error
	go func() {
		err1 = mutex1.Lock()
		wg.Add(-1)
	}()
	go func() {
		err2 = mutex2.Lock()
		wg.Add(-1)
	}()
	wg.Wait()

	assert.Equal(t, false, err1 == nil && err2 == nil)
	assert.Equal(t, true, err1 == nil || err2 == nil)
}

func TestRedisLock_ConcurrentUnlock(t *testing.T) {
	client := getRedisClient("127.0.0.1", 6379, "")
	mutex1 := NewRedisLock(client, "test", time.Second, 3, 100*time.Millisecond)
	mutex2 := NewRedisLock(client, "test", time.Second, 3, 100*time.Millisecond)

	assert.Equal(t, nil, mutex1.Lock())
	assert.NotEqual(t, nil, mutex2.Unlock())
	assert.Equal(t, nil, mutex1.Unlock())
}
