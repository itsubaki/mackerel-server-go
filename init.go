package init

import (
	"net/http"

	"github.com/itsubaki/mackerel-api/pkg/api"
	"github.com/itsubaki/mackerel-api/pkg/mackerel"
)

// GoogleAppEngine Endpoint
func init() {
	http.Handle("/", api.Router(api.Must(mackerel.New())))
}
