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

func (h *SQLHandlerMock) Close() error {
	return nil
}

type TxMock struct {
}

func (tx *TxMock) Exec(query string, args ...interface{}) error {
	return nil
}
