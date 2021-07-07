package database_test

import (
	"time"

	"github.com/itsubaki/mackerel-server-go/pkg/interface/database"
)

type SQLHandlerMock struct {
}

func (h *SQLHandlerMock) Transact(txFunc func(tx database.Tx) error) error {
	return nil
}

func (h *SQLHandlerMock) Query(query string, args ...interface{}) (database.Rows, error) {
	return &RowsMock{}, nil
}

func (h *SQLHandlerMock) QueryRow(query string, args ...interface{}) database.Row {
	return &RowMock{}
}

func (h *SQLHandlerMock) Close() error {
	return nil
}

func (h *SQLHandlerMock) Begin() (database.Tx, error) {
	return &TxMock{}, nil
}

func (h *SQLHandlerMock) SetMaxIdleConns(n int) {

}

func (h *SQLHandlerMock) SetMaxOpenConns(n int) {

}

func (h *SQLHandlerMock) SetConnMaxLifetime(d time.Duration) {

}
func (h *SQLHandlerMock) Raw() interface{} {
	return nil
}

func (h *SQLHandlerMock) IsDebugMode() bool {
	return false
}

type TxMock struct {
}

func (tx *TxMock) Commit() error {
	return nil
}

func (tx *TxMock) Exec(query string, args ...interface{}) (database.Result, error) {
	return &ResultMock{}, nil
}

func (tx *TxMock) Query(query string, args ...interface{}) (database.Rows, error) {
	return &RowsMock{}, nil
}

func (tx *TxMock) QueryRow(query string, args ...interface{}) database.Row {
	return &RowMock{}
}

func (tx *TxMock) Rollback() error {
	return nil
}

type ResultMock struct {
}

func (r *ResultMock) LastInsertId() (int64, error) {
	return -1, nil
}

func (r *ResultMock) RowsAffected() (int64, error) {
	return -1, nil
}

type RowsMock struct {
}

func (r *RowsMock) Close() error {
	return nil
}

func (r *RowsMock) Columns() ([]string, error) {
	return []string{}, nil
}

func (r *RowsMock) Err() error {
	return nil
}

func (r *RowsMock) Next() bool {
	return false
}

func (r *RowsMock) NextResultSet() bool {
	return false
}

func (r *RowsMock) Scan(...interface{}) error {
	return nil
}

type RowMock struct {
}

func (r *RowMock) Scan(...interface{}) error {
	return nil
}
