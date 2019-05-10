package database

import (
	"context"
	"reflect"
)

type SQLHandler interface {
	ShutdownHook()
	Close() error
	Begin() (Tx, error)
	Transact(txFunc func(tx Tx) error) (err error)
	Exec(query string, args ...interface{}) (Result, error)
	Prepare(query string) (Stmt, error)
	Query(query string, args ...interface{}) (Rows, error)
	QueryRow(query string, args ...interface{}) Row
}

type Stmt interface {
	Close() error
	Exec(args ...interface{}) (Result, error)
	ExecContext(ctx context.Context, args ...interface{}) (Result, error)
	Query(args ...interface{}) (Rows, error)
	QueryContext(ctx context.Context, args ...interface{}) (Rows, error)
	QueryRow(args ...interface{}) Row
	QueryRowContext(ctx context.Context, args ...interface{}) Row
}

type Tx interface {
	Commit() error
	Exec(query string, args ...interface{}) (Result, error)
	ExecContext(ctx context.Context, query string, args ...interface{}) (Result, error)
	Prepare(query string) (Stmt, error)
	PrepareContext(ctx context.Context, query string) (Stmt, error)
	Query(query string, args ...interface{}) (Rows, error)
	QueryContext(ctx context.Context, query string, args ...interface{}) (Rows, error)
	QueryRow(query string, args ...interface{}) Row
	QueryRowContext(ctx context.Context, query string, args ...interface{}) Row
	Rollback() error
}

type Result interface {
	LastInsertId() (int64, error)
	RowsAffected() (int64, error)
}

type Rows interface {
	Close() error
	ColumnTypes() ([]ColumnType, error)
	Columns() ([]string, error)
	Err() error
	Next() bool
	NextResultSet() bool
	Scan(...interface{}) error
}

type Row interface {
	Scan(...interface{}) error
}

type ColumnType interface {
	DatabaseTypeName() string
	DecimalSize() (precision, scale int64, ok bool)
	Length() (length int64, ok bool)
	Name() string
	Nullable() (nullable, ok bool)
	ScanType() reflect.Type
}
