package handler

import (
	"database/sql"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/itsubaki/mackerel-server-go/interface/database"
)

var _ database.SQLHandler = (*SQLHandler)(nil)

type SQLHandler struct {
	DB      *sql.DB
	Driver  string
	Dsn     string
	SQLMode string
	Timeout time.Duration
	Sleep   time.Duration
}

type Opt struct {
	SQLMode         string
	Timeout         *time.Duration
	Sleep           *time.Duration
	MaxIdleConns    *int
	MaxOpenConns    *int
	ConnMaxLifetime *time.Duration
}

func New(driver, host, database string, opt ...Opt) (database.SQLHandler, error) {
	sql := fmt.Sprintf("create database if not exists %s", database)
	if err := Exec(driver, host, []string{sql}, opt...); err != nil {
		return nil, fmt.Errorf("query: %v", err)
	}

	return Open(driver, DSN(host, database), opt...)
}

func DSN(host, database string) string {
	if !strings.HasSuffix(host, "/") && !strings.HasPrefix(database, "/") {
		return fmt.Sprintf("%s/%s", host, database)
	}

	return fmt.Sprintf("%s%s", host, database)
}

func Exec(driver, dsn string, query []string, opt ...Opt) error {
	h, err := Open(driver, dsn, opt...)
	if err != nil {
		return fmt.Errorf("open: %v", err)
	}
	defer h.Close()

	if err := h.Transact(func(tx database.Tx) error {
		for _, q := range query {
			if err := tx.Exec(q); err != nil {
				return fmt.Errorf("exec: %v", err)
			}
		}

		return nil
	}); err != nil {
		return fmt.Errorf("transaction: %v", err)
	}

	return nil
}

func Open(driver, dsn string, opt ...Opt) (database.SQLHandler, error) {
	db, err := sql.Open(driver, dsn)
	if err != nil {
		return nil, fmt.Errorf("sql open: %v", err)
	}

	h := &SQLHandler{
		DB:      db,
		Driver:  driver,
		Dsn:     dsn,
		SQLMode: "release",
		Timeout: 10 * time.Minute,
		Sleep:   10 * time.Second,
	}

	if len(opt) > 0 {
		h.SetOpt(opt[0])
	}

	if err := h.Ping(); err != nil {
		return nil, fmt.Errorf("ping: %v", err)
	}

	return h, nil
}

func (h *SQLHandler) SetOpt(opt Opt) {
	if opt.SQLMode != "" {
		h.SQLMode = opt.SQLMode
	}

	if opt.Timeout != nil {
		h.Timeout = *opt.Timeout
	}

	if opt.Sleep != nil {
		h.Sleep = *opt.Sleep
	}

	if opt.MaxIdleConns != nil {
		h.SetMaxIdleConns(*opt.MaxIdleConns)
	}

	if opt.MaxOpenConns != nil {
		h.SetMaxOpenConns(*opt.MaxOpenConns)
	}

	if opt.ConnMaxLifetime != nil {
		h.SetConnMaxLifetime(*opt.ConnMaxLifetime)
	}
}

func (h *SQLHandler) SetMaxIdleConns(n int) {
	h.DB.SetMaxIdleConns(n)
}

func (h *SQLHandler) SetMaxOpenConns(n int) {
	h.DB.SetMaxOpenConns(n)
}

func (h *SQLHandler) SetConnMaxLifetime(d time.Duration) {
	h.DB.SetConnMaxLifetime(d)
}

func (h *SQLHandler) Ping() error {
	start := time.Now()
	for {
		if time.Since(start) > h.Timeout {
			return fmt.Errorf("db ping time over")
		}

		if err := h.DB.Ping(); err != nil {
			log.Printf("db ping: %v", err)
			time.Sleep(h.Sleep)
			continue
		}

		break
	}

	return nil
}

func (h *SQLHandler) Begin() (database.Tx, error) {
	tx, err := h.DB.Begin()
	if err != nil {
		return nil, err
	}

	return &Tx{tx}, nil
}

func (h *SQLHandler) Transact(txFunc func(tx database.Tx) error) (err error) {
	tx, err := h.Begin()
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

	return txFunc(tx)
}

func (h *SQLHandler) Raw() any {
	return h.DB
}

func (h *SQLHandler) IsDebugMode() bool {
	if strings.ToLower(h.SQLMode) == "debug" {
		return true
	}

	return false
}

func (h *SQLHandler) Close() error {
	return h.DB.Close()
}

type Tx struct {
	Tx *sql.Tx
}

func (tx *Tx) Exec(statement string, args ...any) error {
	if _, err := tx.Tx.Exec(statement, args...); err != nil {
		return fmt.Errorf("exec: %v", err)
	}

	return nil
}

func (tx *Tx) Commit() error {
	return tx.Tx.Commit()
}

func (tx *Tx) Rollback() error {
	return tx.Tx.Rollback()
}
