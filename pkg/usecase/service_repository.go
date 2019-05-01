package usecase

import "github.com/itsubaki/mackerel-api/pkg/domain"

type ServiceRepository interface {
	List() (*domain.Services, error)
	Exists(serviceName string) bool
	Service(serviceName string) (*domain.Service, error)
	Save(service *domain.Service) error
	Delete(serviceName string) error

	ExistsMetadata(serviceName, namespace string) bool
	MetadataList(serviceName string) (*domain.ServiceMetadataList, error)
	Metadata(serviceName, namespace string) (interface{}, error)
	SaveMetadata(serviceName, namespace string, metadata interface{}) (*domain.Success, error)
	DeleteMetadata(serviceName, namespace string) (*domain.Success, error)

	ExistsRole(serviceName, roleName string) bool
	RoleList(serviceName string) (*domain.Roles, error)
	Role(serviceName, roleName string) (*domain.Role, error)
	SaveRole(serviceName string, role *domain.Role) error
	DeleteRole(serviceName, roleName string) error

	ExistsRoleMetadata(serviceName, roleName, namespace string) bool
	RoleMetadataList(serviceName, roleName string) (*domain.RoleMetadataList, error)
	RoleMetadata(serviceName, roleName, namespace string) (interface{}, error)
	SaveRoleMetadata(serviceName, roleName, namespace string, metadata interface{}) (*domain.Success, error)
	DeleteRoleMetadata(serviceName, roleName, namespace string) (*domain.Success, error)

	ExistsMetric(serviceName, metricName string) bool
	MetricNames(serviceName string) (*domain.ServiceMetricValueNames, error)
	MetricValues(serviceName, metricName string, from, to int) (*domain.ServiceMetricValues, error)
	SaveMetricValues(serviceName string, values []domain.ServiceMetricValue) (*domain.Success, error)
}
