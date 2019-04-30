package usecase

import "github.com/itsubaki/mackerel-api/pkg/domain"

type ServiceRepository interface {
	MetricNames(serviceName string) ([]string, error)
	MetricValues(serviceName, metricName string, from, to int64) (domain.ServiceMetricValues, error)
	SaveMetricValues(values domain.ServiceMetricValues) error
	RoleList(serviceName string) (domain.Roles, error)
	Role(serviceName, roleName string) (*domain.Role, error)
	SaveRole(role domain.Role) error
	DeleteRole(serviceName, roleName string) error
	List() (domain.Services, error)
	Service(serviceName string) (domain.Service, error)
	Save(service domain.Service) error
	Delete(serviceName string) error
	Exists(serviceName string) bool
}
