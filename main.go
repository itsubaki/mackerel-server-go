package main

import (
	"log"

	"github.com/itsubaki/mackerel-api/pkg/infrastructure"
)

// GoogleAppEngine endpoint
// func init() {
// 	http.Handle("/", infrastructure.Default())
// }

// CommandLine endpoint
func main() {
	if err := infrastructure.Default().Run(":8080"); err != nil {
		log.Fatalf("run mackerel-api: %v", err)
	}
}
