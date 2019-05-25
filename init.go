// +build appengine

package main

import (
	"net/http"

	"github.com/itsubaki/mackerel-api/pkg/infrastructure"
)

// GoogleAppEngine endpoint
func init() {
	handler := infrastructure.NewSQLHandler()
	http.Handle("/", infrastructure.Router(handler))
}
