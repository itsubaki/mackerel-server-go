package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/go-sql-driver/mysql"

	"github.com/itsubaki/mackerel-server-go/pkg/infrastructure"
	"github.com/itsubaki/mackerel-server-go/pkg/infrastructure/config"
	"github.com/itsubaki/mackerel-server-go/pkg/infrastructure/handler"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	c := config.New()
	log.Printf("%#v\n", c)

	h, err := handler.New(c.Driver, c.Host, c.Database, handler.Opt{
		SQLMode:         c.SQLMode,
		Timeout:         &c.Timeout,
		Sleep:           &c.Sleep,
		MaxIdleConns:    &c.MaxIdleConns,
		MaxOpenConns:    &c.MaxOpenConns,
		ConnMaxLifetime: &c.ConnMaxLifetime,
	})
	if err != nil {
		log.Printf("handler new: %v", err)
	}
	log.Printf("db connected")

	if err := infrastructure.RunFixture(h); err != nil {
		log.Fatalf("run fixture: %v", err)
	}

	r := infrastructure.Default()
	infrastructure.UseSession(r)
	infrastructure.APIv0(r, h)

	s := &http.Server{
		Addr:    fmt.Sprintf(":%s", c.Port),
		Handler: r,
	}

	go func() {
		log.Println("http server listen and serve")
		if err := s.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	ch := make(chan os.Signal, 2)
	signal.Notify(ch, os.Interrupt, syscall.SIGTERM)
	<-ch

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := s.Shutdown(ctx); err != nil {
		log.Fatalf("http server shutdown: %v\n", err)
	}

	if err := h.Close(); err != nil {
		log.Fatalf("handler closed: %v\n", err)
	}

	log.Println("shutdown finished")
}
