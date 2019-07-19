# Go Locks

- map lock
- redis lock
- mysql lock

## Usage

```go
import (
	"time"
	"github.com/study-only/go-locks"
)


// Map Lock
mapLock := golocks.NewMapLock("lock1")
err := mapLock.TryLock()
err = mapLock.Unlock()

// Redis Lock
InitRedisLock(redisClient)
redisLock := NewRedisLock("lock2", time.Second)
err = redisLock.TryLock()
err = redisLock.Unlock()

// Mysql Lock
InitMysqlLock(db, "go_lock", 5*time.Second)
mysqlLock := NewMysqlLock("lock3", time.Second)
err = mysqlLock.TryLock()
err = mysqlLock.Unlock()

// Upgrade to Spin Lock
mapSpinLock := NewSpinLock(NewMapLock(lockKey), spinTries, spinInterval)
err = mapSpinLock.Lock()
err = mapSpinLock.Unlock()

```