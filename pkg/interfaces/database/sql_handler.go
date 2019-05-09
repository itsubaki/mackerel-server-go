package database

type SQLHandler interface {
	ShutdownHook()
	Close() error
	Begin() (Tx, error)
	Transact(txFunc func(tx Tx) error) (err error)
}

type Tx interface {
	Exec(string, ...interface{}) (Result, error)
	Query(string, ...interface{}) (Rows, error)
	Commit() error
	Rollback() error
}

type Result interface {
	LastInsertId() (int64, error)
	RowsAffected() (int64, error)
}

type Rows interface {
	Scan(...interface{}) error
	Next() bool
	Close() error
}
