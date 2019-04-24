package api

import (
	"github.com/itsubaki/mackerel-api/internal/services"
)

type Mackerel interface {
	GetServices() *services.GetServicesOutput
}
