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
