package golocks

import (
	"database/sql"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var (
	testMysqlHost   = env("TEST_MYSQL_HOST", "127.0.0.1")
	testMysqlPort   = env("TEST_MYSQL_PORT", "3306")
	testMysqlUser   = env("TEST_MYSQL_USER", "root")
	testMysqlPwd    = env("TEST_MYSQL_PWD", "")
	testMysqlDbName = env("TEST_MYSQL_DB_NAME", "test")
)

func TestMysqlLock_Lock(t *testing.T) {
	db := getDB(t, testMysqlHost, testMysqlPort, testMysqlUser, testMysqlPwd, testMysqlDbName)
	InitMysqlLock(db, "go_lock", time.Second)
	lock := NewMysqlLock("lock", time.Second)

	err := lock.TryLock()
	assert.Nil(t, err)
	err = lock.Unlock()
	assert.Nil(t, err)
}

func TestMysqlLock_Expired(t *testing.T) {
	db := getDB(t, testMysqlHost, testMysqlPort, testMysqlUser, testMysqlPwd, testMysqlDbName)
	InitMysqlLock(db, "go_lock", 200*time.Millisecond)
	lock := NewMysqlLock("expiry", 500*time.Millisecond)

	err := lock.TryLock()
	assert.Nil(t, err)

	time.Sleep(2 * time.Second)
	assert.Nil(t, err)

	err = lock.Unlock()
	assert.NotNil(t, err)
	err = lock.TryLock()
	assert.Nil(t, err)
}

func getDB(t *testing.T, host, port, user, pwd, dbName string) *sql.DB {
	connStr := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?loc=Local&parseTime=true", user, pwd, host, port, dbName)
	var err error
	db, err := sql.Open("mysql", connStr)

	if err == nil {
		err = db.Ping()
	}

	if err != nil {
		t.Fatalf("sql db connect error %s : %#v", connStr, err)
	}

	return db
}

func env(name, defaultValue string) string {
	val := os.Getenv(name)
	if val != "" {
		return val
	}

	return defaultValue
}
