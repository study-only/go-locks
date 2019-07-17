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


// map lock
mapLock := golocks.NewMapLock("lock1")
err := mapLock.TryLock()
err = mapLock.Unlock()

// redis lock
InitRedisLock(redisClient)
redisLock := NewRedisLock("lock2", time.Second)
err = redisLock.TryLock()
err = redisLock.Unlock()

// mysql lock
InitMysqlLock(db, "go_lock", 5*time.Second)
mysqlLock := NewMysqlLock("lock3", time.Second)
err = mysqlLock.TryLock()
err = mysqlLock.Unlock()

// spin lock
mapSpinLock := NewSpinLock(NewMapLock(lockKey), spinTries, spinInterval)
err = mapSpinLock.Lock()
err = mapSpinLock.Unlock()

```