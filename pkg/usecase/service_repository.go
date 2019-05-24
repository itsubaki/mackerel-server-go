package usecase

import "github.com/itsubaki/mackerel-api/pkg/domain"

type ServiceRepository interface {
	List(org string) (*domain.Services, error)
	Exists(org, serviceName string) bool
	Service(org, serviceName string) (*domain.Service, error)
	Save(org string, service *domain.Service) error
	Delete(org, serviceName string) error

	ExistsMetadata(org, serviceName, namespace string) bool
	MetadataList(org, serviceName string) (*domain.ServiceMetadataList, error)
	Metadata(org, serviceName, namespace string) (interface{}, error)
	SaveMetadata(org, serviceName, namespace string, metadata interface{}) (*domain.Success, error)
	DeleteMetadata(org, serviceName, namespace string) (*domain.Success, error)

	ExistsRole(org, serviceName, roleName string) bool
	RoleList(org, serviceName string) (*domain.Roles, error)
	Role(org, serviceName, roleName string) (*domain.Role, error)
	SaveRole(org, serviceName string, role *domain.Role) error
	DeleteRole(org, serviceName, roleName string) error

	ExistsRoleMetadata(org, serviceName, roleName, namespace string) bool
	RoleMetadataList(org, serviceName, roleName string) (*domain.RoleMetadataList, error)
	RoleMetadata(org, serviceName, roleName, namespace string) (interface{}, error)
	SaveRoleMetadata(org, serviceName, roleName, namespace string, metadata interface{}) (*domain.Success, error)
	DeleteRoleMetadata(org, serviceName, roleName, namespace string) (*domain.Success, error)

	ExistsMetric(org, serviceName, metricName string) bool
	MetricNames(org, serviceName string) (*domain.ServiceMetricValueNames, error)
	MetricValues(org, serviceName, metricName string, from, to int64) (*domain.ServiceMetricValues, error)
	SaveMetricValues(org, serviceName string, values []domain.ServiceMetricValue) (*domain.Success, error)
}
