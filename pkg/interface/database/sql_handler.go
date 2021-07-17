package database

import "time"

type SQLHandler interface {
	Transact(txFunc func(tx Tx) error) error
	SetMaxIdleConns(n int)
	SetMaxOpenConns(n int)
	SetConnMaxLifetime(d time.Duration)
	Raw() interface{}
	IsDebugMode() bool
	Close() error
}

type Tx interface {
	Exec(query string, args ...interface{}) error
}
