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
l1 := golocks.NewMapLock("lock1", 10 * time.MilliSecond)
err := l1.Lock()
err = l1.Unlock()

// redis lock
InitRedisLock(redisClient)
l2 := NewRedisLock("lock2", time.Second, 50, 100*time.Millisecond)
err = l2.Lock()
err = l2.Unlock()

// mysql lock
InitMysqlLock(db, "go_lock")
l3 := NewMysqlLock("lock3", time.Second, 5, time.Second)
err = l3.Lock()
err = l3.Unlock()

```