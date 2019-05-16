// +build appengine

package main

import (
	"net/http"

	"github.com/itsubaki/mackerel-api/pkg/infrastructure"
)

// GoogleAppEngine endpoint
func init() {
	h := infrastructure.NewSQLHandler()
	r := infrastructure.Router(h)

	http.Handle("/", r)
}
