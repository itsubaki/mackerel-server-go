// +build !appengine
package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/itsubaki/mackerel-api/pkg/infrastructure"
	"github.com/itsubaki/mackerel-api/pkg/interfaces/database"
)

// CommandLine endpoint
func main() {
	h := infrastructure.NewSQLHandler()
	r := infrastructure.Router(h)

	ShutdownHook(h)
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("run mackerel-api: %v", err)
	}
}

func ShutdownHook(h database.SQLHandler) {
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
