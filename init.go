package init

import (
	"net/http"

	"github.com/itsubaki/mackerel-api/pkg/mackerel"
)

// GoogleAppEngine endpoint
func init() {
	http.Handle("/", mackerel.Default())
}
