package init

import (
	"net/http"

	"github.com/itsubaki/mackerel-api/pkg/mackerel"
)

// GoogleAppEngine Endpoint
func init() {
	http.Handle("/", mackerel.Default())
}
