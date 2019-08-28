package golocks

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

const (
	createTableSql = `
CREATE TABLE IF NOT EXISTS %s (
  name       varchar(255) NOT NULL,
  expire_at  timestamp    NOT NULL,
  created_at timestamp    NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (name) USING HASH,
  KEY idx_expire_at (expire_at) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8
`
	insertRowSql = `
INSERT INTO %s (name, expire_at, created_at) VALUES (?, ?, ?)
`
	deleteRowSql = `
DELETE FROM %s WHERE name=? LIMIT 1
`
	deleteExprieRowsSql = `
DELETE FROM %s WHERE expire_at<?
`
)

var lockDb *sql.DB
var lockTableName string

func InitMysqlLock(db *sql.DB, tableName string, clearExpiryInterval time.Duration) {
	if db == nil {
		panic("db is nil")
	}

	lockDb = db
	lockTableName = tableName

	if err := createTable(); err != nil {
		panic(err)
	}
	if err := deleteExpireLocks(); err != nil {
		panic(err)
	}

	go func() {
		for {
			deleteExpireLocks()
			time.Sleep(clearExpiryInterval)
		}
	}()
}

func NewMysqlLock(name string, expiry time.Duration) *mysqlLock {
	return &mysqlLock{
		name:   name,
		expiry: expiry,
	}
}

type mysqlLock struct {
	name   string
	expiry time.Duration

	startAt time.Time
	isOwner bool
}

func (l *mysqlLock) TryLock() error {
	createdAt := time.Now()
	expireAt := createdAt.Add(l.expiry)
	if err := insertRow(l.name, expireAt, createdAt); err != nil {
		return errorf(fmt.Sprintf("mysql lock: %s already locked", l.name))
	}

	l.startAt = time.Now()
	l.isOwner = true
	return nil
}

func (l *mysqlLock) Unlock() error {
	if !l.isOwner {
		return errorf("mysql lock: not owner")
	}
	if time.Now().UnixNano()-l.startAt.UnixNano() >= l.expiry.Nanoseconds() {
		return errorf("mysql lock: lock expired")
	}

	if err := deleteRow(l.name); err != nil {
		return err
	}

	l.isOwner = false
	return nil
}

func createTable() error {
	query := fmt.Sprintf(createTableSql, lockTableName)
	if _, err := lockDb.Exec(query); err != nil {
		return errorf("mysql lock: %s", err)
	}

	return nil
}

func insertRow(name string, expireAt, createdAt time.Time) error {
	query := fmt.Sprintf(insertRowSql, lockTableName)
	if _, err := lockDb.Exec(query, name, expireAt, createdAt); err != nil {
		return errorf("mysql lock: %s", err)
	}

	return nil
}

func deleteRow(name string) error {
	query := fmt.Sprintf(deleteRowSql, lockTableName)
	if _, err := lockDb.Exec(query, name); err != nil {
		return errorf("mysql lock: %s", err)
	}

	return nil
}

func deleteExpireLocks() error {
	query := fmt.Sprintf(deleteExprieRowsSql, lockTableName)
	if _, err := lockDb.Exec(query, time.Now()); err != nil {
		return errorf("mysql lock: %s", err)
	}

	return nil
}
