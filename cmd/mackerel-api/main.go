package main

import (
	"log"

	"github.com/itsubaki/mackerel-api/pkg/mackerel"
)

// CommandLine endpoint
func main() {
	if err := mackerel.Default().Run(":8080"); err != nil {
		log.Fatalf("run mackerel-api: %v", err)
	}
}
