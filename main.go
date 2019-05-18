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
	var handler database.SQLHandler
	if os.Getenv("MACKEREL_API_PERSISTENCE") == "database" {
		handler = infrastructure.NewSQLHandler()
		c := make(chan os.Signal, 2)
		signal.Notify(c, os.Interrupt, syscall.SIGTERM)

		go func() {
			<-c
			if err := handler.Close(); err != nil {
				panic(err)
			}

			os.Exit(0)
		}()
	}

	if err := infrastructure.Router(handler).Run(":8080"); err != nil {
		log.Fatalf("run mackerel-api: %v", err)
	}
}
