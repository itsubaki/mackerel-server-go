package mackerel

import (
	"github.com/itsubaki/mackerel-api/internal/services"
)

func New() (*Mackerel, error) {
	return &Mackerel{}, nil
}

type Mackerel struct{}

func (m *Mackerel) GetServices() *services.GetServicesOutput {
	return &services.GetServicesOutput{
		Services: []services.Service{},
	}
}
