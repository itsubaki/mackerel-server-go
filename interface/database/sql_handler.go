package database

import "time"

type SQLHandler interface {
	Begin() (Tx, error)
	Transact(txFunc func(tx Tx) error) error
	SetMaxIdleConns(n int)
	SetMaxOpenConns(n int)
	SetConnMaxLifetime(d time.Duration)
	Raw() any
	IsDebugMode() bool
	Close() error
}

type Tx interface {
	Exec(query string, args ...any) error
	Commit() error
	Rollback() error
}
