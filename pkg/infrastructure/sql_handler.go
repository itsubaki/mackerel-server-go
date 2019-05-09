package infrastructure

import (
	"database/sql"
	"os"
	"os/signal"
	"syscall"

	"github.com/itsubaki/mackerel-api/pkg/interfaces/database"
)

type SQLHandler struct {
	DB *sql.DB
}

func NewSQLHandler() database.SQLHandler {
	db, err := sql.Open("mysql", "root:secret@tcp(127.0.0.1:3307)/mackerel")
	if err != nil {
		panic(err)
	}

	return &SQLHandler{
		DB: db,
	}
}

func (h *SQLHandler) Transact(txFunc func(database.Tx) error) (err error) {
	tx, err := h.DB.Begin()
	if err != nil {
		return
	}

	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p)
		}

		if err != nil {
			tx.Rollback()
			return
		}

		err = tx.Commit()
	}()

	return txFunc(&Tx{tx})
}

func (h *SQLHandler) ShutdownHook() {
	c := make(chan os.Signal, 2)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		if err := h.Close(); err != nil {
			panic(err)
		}
		os.Exit(0)
	}()
}

func (h *SQLHandler) Close() error {
	return h.DB.Close()
}

func (h *SQLHandler) Begin() (database.Tx, error) {
	tx, err := h.DB.Begin()
	if err != nil {
		return nil, err
	}

	return &Tx{tx}, nil
}

type Tx struct {
	Tx *sql.Tx
}

func (tx *Tx) Exec(statement string, args ...interface{}) (database.Result, error) {
	result, err := tx.Tx.Exec(statement, args...)
	if err != nil {
		return nil, err
	}

	return &Result{result}, nil
}

func (tx *Tx) Query(statement string, args ...interface{}) (database.Rows, error) {
	rows, err := tx.Tx.Query(statement, args...)
	if err != nil {
		return nil, err
	}

	return &Rows{rows}, nil
}

func (tx *Tx) Commit() error {
	return tx.Tx.Commit()
}

func (tx *Tx) Rollback() error {
	return tx.Tx.Rollback()
}

type Result struct {
	Result sql.Result
}

func (r *Result) LastInsertId() (int64, error) {
	return r.Result.LastInsertId()
}

func (r *Result) RowsAffected() (int64, error) {
	return r.Result.RowsAffected()
}

type Rows struct {
	Rows *sql.Rows
}

func (r *Rows) Scan(dest ...interface{}) error {
	return r.Rows.Scan(dest...)
}

func (r *Rows) Next() bool {
	return r.Rows.Next()
}

func (r *Rows) Close() error {
	return r.Rows.Close()
}
