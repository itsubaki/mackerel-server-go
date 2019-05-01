package usecase

import "github.com/itsubaki/mackerel-api/pkg/domain"

type ServiceRepository interface {
	List() (*domain.Services, error)
	Service(serviceName string) (*domain.Service, error)
	Save(service *domain.Service) error
	Delete(serviceName string) error
	Exists(serviceName string) bool

	RoleList(serviceName string) (*domain.Roles, error)
	Role(serviceName, roleName string) (*domain.Role, error)
	SaveRole(serviceName string, role *domain.Role) error
	DeleteRole(serviceName, roleName string) error

	MetricNames(serviceName string) (*domain.ServiceMetricValueNames, error)
	MetricValues(serviceName, metricName string, from, to int64) (*domain.ServiceMetricValues, error)
	SaveMetricValues(serviceName string, values []domain.ServiceMetricValue) error
}
