package main

import (
	"net/http"

	"github.com/itsubaki/mackerel-api/pkg/api"
	"github.com/itsubaki/mackerel-api/pkg/mackerel"
)

func init() {
	h := api.Handler(api.Must(mackerel.New()))
	http.Handle("/", h)
}
