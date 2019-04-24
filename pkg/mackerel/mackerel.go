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

func (m *Mackerel) PostServices(in *services.PostServicesInput) *services.PostServicesOutput {
	return &services.PostServicesOutput{}
}

func (m *Mackerel) DeleteServices(in *services.DeleteServicesInput) *services.DeleteServicesOutput {
	return &services.DeleteServicesOutput{}
}

func (m *Mackerel) GetRoles(in *services.GetRolesInput) *services.GetRolesOutput {
	return &services.GetRolesOutput{}
}

func (m *Mackerel) PostRoles(in *services.PostRolesInput) *services.PostRolesOutput {
	return &services.PostRolesOutput{}
}

func (m *Mackerel) DeleteRoles(in *services.DeleteRolesInput) *services.DeleteRolesOutput {
	return &services.DeleteRolesOutput{}
}

func (m *Mackerel) GetMetricNames(in *services.GetMetricNamesInput) *services.GetMetricNamesOutput {
	return &services.GetMetricNamesOutput{}
}
