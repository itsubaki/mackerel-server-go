// +build appengine

package main

import (
	"net/http"
	"os"

	"github.com/itsubaki/mackerel-api/pkg/infrastructure"
	"github.com/itsubaki/mackerel-api/pkg/interfaces/database"
)

// GoogleAppEngine endpoint
func init() {
	var handler database.SQLHandler
	if os.Getenv("MACKEREL_API_PERSISTENCE") == "database" {
		handler = infrastructure.NewSQLHandler()
	}

	http.Handle("/", infrastructure.Router(handler))
}
