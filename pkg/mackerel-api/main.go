package main

import (
	"log"
	"os"

	"github.com/itsubaki/mackerel-api/pkg/mackerel"
)

func main() {
	port := ":8080"
	if p := os.Getenv("MACKEREL_API_PORT"); len(p) > 0 {
		port = p
	}

	if err := mackerel.Default().Run(port); err != nil {
		log.Fatalf("run mackerel-api: %v", err)
	}
}
