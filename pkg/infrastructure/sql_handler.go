package infrastructure

import (
	"os"
	"os/signal"
	"syscall"
)

type SQLHandler struct {
}

func NewSQLHandler() *SQLHandler {
	return &SQLHandler{}
}

func (handler *SQLHandler) ShutdownHook() {
	c := make(chan os.Signal, 2)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		os.Exit(0)
	}()
}
