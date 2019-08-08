package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/itsubaki/mackerel-api/pkg/infrastructure"
)

func main() {
	config := infrastructure.NewConfig()
	fmt.Printf("%#v\n", config)

	handler := infrastructure.NewSQLHandler(config)

	c := make(chan os.Signal, 2)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-c
		if err := handler.Close(); err != nil {
			panic(err)
		}

		os.Exit(0)
	}()

	if err := infrastructure.Router(handler).Run(config.Port); err != nil {
		log.Fatalf("run mackerel-api: %v", err)
	}
}
