// +build !appengine
package infrastructure

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/itsubaki/mackerel-api/pkg/interfaces/database"
)

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
