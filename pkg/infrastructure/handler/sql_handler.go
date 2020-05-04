package handler

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/itsubaki/mackerel-api/pkg/infrastructure/config"
	"github.com/itsubaki/mackerel-api/pkg/interface/database"
)

type SQLHandler struct {
	DB *sql.DB
}

func New(config *config.Config) (database.SQLHandler, error) {
	q := fmt.Sprintf("create database if not exists %s", config.DatabaseName)
	return NewWith(config, []string{q})
}

func NewWith(config *config.Config, query []string) (database.SQLHandler, error) {
	db, err := Wait(config)
	if err != nil {
		return nil, fmt.Errorf("wait: %v", err)
	}

	if err := db.Transact(func(tx database.Tx) error {
		for _, q := range query {
			if _, err := tx.Exec(q); err != nil {
				return fmt.Errorf("exec: %v", err)
			}
		}

		return nil
	}); err != nil {
		panic(fmt.Errorf("transaction: %v", err))
	}

	return Open(config)
}

func Open(config *config.Config) (database.SQLHandler, error) {
	source := fmt.Sprintf("%s%s", config.DataSourceName, config.DatabaseName)
	db, err := sql.Open(config.Driver, source)
	if err != nil {
		return nil, fmt.Errorf("sql open: %v", err)
	}

	return &SQLHandler{
		DB: db,
	}, nil
}

func Wait(config *config.Config) (database.SQLHandler, error) {
	db, err := sql.Open(config.Driver, config.DataSourceName)
	if err != nil {
		return nil, fmt.Errorf("sql open: %v", err)
	}

	start := time.Now()
	for {
		if time.Since(start) > 10*time.Minute {
			return nil, fmt.Errorf("db ping time over")
		}

		if err := db.Ping(); err != nil {
			log.Printf("db ping: %v", err)
			time.Sleep(10 * time.Second)
			continue
		}

		break
	}

	return &SQLHandler{
		DB: db,
	}, nil
}

func (h *SQLHandler) Query(query string, args ...interface{}) (database.Rows, error) {
	rows, err := h.DB.Query(query, args...)
	if err != nil {
		return nil, err
	}

	return &Rows{rows}, nil
}

func (h *SQLHandler) QueryRow(query string, args ...interface{}) database.Row {
	return h.DB.QueryRow(query, args...)
}

func (h *SQLHandler) Transact(txFunc func(tx database.Tx) error) (err error) {
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

func (tx *Tx) Commit() error {
	return tx.Tx.Commit()
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

func (tx *Tx) QueryRow(query string, args ...interface{}) database.Row {
	row := tx.Tx.QueryRow(query, args...)
	return &Row{row}
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

func (r *Rows) Close() error {
	return r.Rows.Close()
}

func (r *Rows) Columns() ([]string, error) {
	return r.Rows.Columns()
}

func (r *Rows) Err() error {
	return r.Rows.Err()
}

func (r *Rows) Next() bool {
	return r.Rows.Next()
}

func (r *Rows) NextResultSet() bool {
	return r.Rows.NextResultSet()
}

func (r *Rows) Scan(dest ...interface{}) error {
	return r.Rows.Scan(dest...)
}

type Row struct {
	Row *sql.Row
}

func (r *Row) Scan(dest ...interface{}) error {
	return r.Row.Scan(dest...)
}
