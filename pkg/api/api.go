package api

import (
	"github.com/itsubaki/mackerel-api/internal/services"
)

type Mackerel interface {
	GetServices() *services.GetServicesOutput
	PostServices(in *services.PostServicesInput) *services.PostServicesOutput
	DeleteServices(in *services.DeleteServicesInput) *services.DeleteServicesOutput

	GetRoles(in *services.GetRolesInput) *services.GetRolesOutput
	PostRoles(in *services.PostRolesInput) *services.PostRolesOutput
	DeleteRoles(in *services.DeleteRolesInput) *services.DeleteRolesOutput

	GetMetricNames(in *services.GetMetricNamesInput) *services.GetMetricNamesOutput
}
