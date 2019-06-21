// +build !appengine

package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/itsubaki/mackerel-api/pkg/infrastructure"
)

// CommandLine endpoint
func main() {
	handler := infrastructure.NewSQLHandler()

	c := make(chan os.Signal, 2)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-c
		if err := handler.Close(); err != nil {
			panic(err)
		}

		os.Exit(0)
	}()

	if err := infrastructure.Default(handler).Run(":8080"); err != nil {
		log.Fatalf("run mackerel-api: %v", err)
	}
}
