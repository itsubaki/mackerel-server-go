package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/itsubaki/mackerel-api/pkg/infrastructure"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	config := infrastructure.NewConfig()
	log.Printf("%#v\n", config)

	handler := infrastructure.NewSQLHandler(config)

	s := &http.Server{
		Addr:    config.Port,
		Handler: infrastructure.Router(handler),
	}

	go func() {
		log.Println("http server listen and serve")
		if err := s.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	c := make(chan os.Signal, 2)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := s.Shutdown(ctx); err != nil {
		log.Fatalf("http server shutdown: %v\n", err)
	}
	log.Println("http server closed")

	if err := handler.Close(); err != nil {
		log.Fatalf("handler closed: %v\n", err)
	}
	log.Println("db disconnected")

	log.Println("shutdown finished")
}
