package domain

import (
	"database/sql"
	"reflect"
	"time"
)

type SQLDriver interface {
	LookupDefaultType(typ reflect.Type) string
}

// Conn is database connection like DB or Tx. This is also implemented by MockDB and MockTx.
type Conn interface {
	GetDriver() SQLDriver
	Ping() error
	Query(query string, args ...interface{}) (*sql.Rows, error)
	Exec(query string, args ...interface{}) (sql.Result, error)
}

// DB is the interface of database.
type DB interface {
	Conn
	SetConnMaxLifetime(n time.Duration) error
	SetMaxIdleConns(n int) error
	SetMaxOpenConns(n int) error
	Close() error
	Begin() (Tx, error)
}

// Tx is the interface of database transaction.
type Tx interface {
	Conn
	Commit() error
	Rollback() error
}

// Mock is mock database conneciton pool.
type Mock interface {
	Conn
	Complete() error
	CompareWith(Stmt) (interface{}, error)
	Expect(s Stmt)
	ExpectWithReturn(s Stmt, v interface{})
}

type MockDB interface {
	Mock
	SetConnMaxLifetime(n time.Duration) error
	SetMaxIdleConns(n int) error
	SetMaxOpenConns(n int) error
	Close() error
	Begin() (MockTx, error)
	ExpectBegin() MockTx
}

type MockTx interface {
	Mock
	Commit() error
	Rollback() error
	ExpectCommit()
	ExpectRollback()
}
