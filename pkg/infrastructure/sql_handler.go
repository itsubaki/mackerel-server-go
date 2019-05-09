package infrastructure

import (
	"database/sql"
	"os"
	"os/signal"
	"reflect"
	"syscall"

	_ "github.com/go-sql-driver/mysql"
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

func (h *SQLHandler) Exec(query string, args ...interface{}) (database.Result, error) {
	result, err := h.DB.Exec(query, args...)
	if err != nil {
		return nil, err
	}

	return &Result{result}, nil
}

func (h *SQLHandler) Prepare(query string) (database.Stmt, error) {
	stmt, err := h.DB.Prepare(query)
	if err != nil {
		return nil, err
	}

	return &Stmt{stmt}, nil
}

func (h *SQLHandler) Query(query string, args ...interface{}) (database.Rows, error) {
	rows, err := h.DB.Query(query, args...)
	if err != nil {
		return nil, err
	}

	return &Rows{rows}, nil
}

func (h *SQLHandler) QueryRow(query string, args ...interface{}) database.Row {
	row := h.DB.QueryRow(query, args...)
	return &Row{row}
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

type Stmt struct {
	stmt *sql.Stmt
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

func (r *Rows) Close() error {
	return r.Rows.Close()
}

func (r *Rows) ColumnTypes() ([]database.ColumnType, error) {
	out := []database.ColumnType{}

	list, err := r.Rows.ColumnTypes()
	if err != nil {
		return out, err
	}

	for i := range list {
		out = append(out, list[i])
	}

	return out, nil
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

type ColumnType struct {
	ColumnType *sql.ColumnType
}

func (c *ColumnType) DatabaseTypeName() string {
	return c.ColumnType.DatabaseTypeName()
}

func (c *ColumnType) DecimalSize() (precision, scale int64, ok bool) {
	return c.ColumnType.DecimalSize()
}

func (c *ColumnType) Length() (length int64, ok bool) {
	return c.ColumnType.Length()
}

func (c *ColumnType) Name() string {
	return c.ColumnType.Name()
}

func (c *ColumnType) Nullable() (nullable, ok bool) {
	return c.ColumnType.Nullable()
}

func (c *ColumnType) ScanType() reflect.Type {
	return c.ColumnType.ScanType()
}
