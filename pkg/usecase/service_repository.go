package usecase

import "github.com/itsubaki/mackerel-api/pkg/domain"

type ServiceRepository interface {
	List(orgID string) (*domain.Services, error)
	Exists(orgID, serviceName string) bool
	Service(orgID, serviceName string) (*domain.Service, error)
	Save(orgID string, service *domain.Service) error
	Delete(orgID, serviceName string) error

	ExistsMetadata(orgID, serviceName, namespace string) bool
	MetadataList(orgID, serviceName string) (*domain.ServiceMetadataList, error)
	Metadata(orgID, serviceName, namespace string) (interface{}, error)
	SaveMetadata(orgID, serviceName, namespace string, metadata interface{}) (*domain.Success, error)
	DeleteMetadata(orgID, serviceName, namespace string) (*domain.Success, error)

	ExistsRole(orgID, serviceName, roleName string) bool
	RoleList(orgID, serviceName string) (*domain.Roles, error)
	Role(orgID, serviceName, roleName string) (*domain.Role, error)
	SaveRole(orgID, serviceName string, role *domain.Role) error
	DeleteRole(orgID, serviceName, roleName string) error

	ExistsRoleMetadata(orgID, serviceName, roleName, namespace string) bool
	RoleMetadataList(orgID, serviceName, roleName string) (*domain.RoleMetadataList, error)
	RoleMetadata(orgID, serviceName, roleName, namespace string) (interface{}, error)
	SaveRoleMetadata(orgID, serviceName, roleName, namespace string, metadata interface{}) (*domain.Success, error)
	DeleteRoleMetadata(orgID, serviceName, roleName, namespace string) (*domain.Success, error)

	ExistsMetric(orgID, serviceName, metricName string) bool
	MetricNames(orgID, serviceName string) (*domain.ServiceMetricValueNames, error)
	MetricValues(orgID, serviceName, metricName string, from, to int64) (*domain.ServiceMetricValues, error)
	SaveMetricValues(orgID, serviceName string, values []domain.ServiceMetricValue) (*domain.Success, error)
}
