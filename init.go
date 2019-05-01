package init

import (
	"net/http"

	"github.com/itsubaki/mackerel-api/pkg/infrastructure"
)

// GoogleAppEngine endpoint
func init() {
	http.Handle("/", infrastructure.Default())
}
