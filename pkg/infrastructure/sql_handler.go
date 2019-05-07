package infrastructure

import (
	"database/sql"
	"os"
	"os/signal"
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

func (handler *SQLHandler) ShutdownHook() {
	c := make(chan os.Signal, 2)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		if err := handler.DB.Close(); err != nil {
			panic(err)
		}
		os.Exit(0)
	}()
}
