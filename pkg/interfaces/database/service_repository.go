package database

import (
	"fmt"

	"github.com/itsubaki/mackerel-api/pkg/domain"
)

type ServiceRepository struct {
	SQLHandler SQLHandler
}

func NewServiceRepository(handler SQLHandler) *ServiceRepository {
	return &ServiceRepository{
		SQLHandler: handler,
	}
}

// select * from services
func (repo *ServiceRepository) List() (*domain.Services, error) {
	return nil, nil
}

// insert into services values()
func (repo *ServiceRepository) Save(s *domain.Service) error {
	return nil
}

// select * from services where service_name=${serviceName}
func (repo *ServiceRepository) Service(serviceName string) (*domain.Service, error) {
	return nil, fmt.Errorf("service not found")
}

// select * from services where service_name=${serviceName} limit=1
func (repo *ServiceRepository) Exists(serviceName string) bool {
	return false
}

// delete from services where service_name=${serviceName}
func (repo *ServiceRepository) Delete(serviceName string) error {
	return nil
}

// select * from service_metadata where service_name=${serviceName} and namespace=${namespace} limit=1
func (repo *ServiceRepository) ExistsMetadata(serviceName, namespace string) bool {
	return true
}

// select namespacee from service_metadata where service_name=${serviceName}
func (repo *ServiceRepository) MetadataList(serviceName string) (*domain.ServiceMetadataList, error) {
	return nil, nil
}

// select * from service_metadata where service_name=${serviceName} and namespace=${namespace}
func (repo *ServiceRepository) Metadata(serviceName, namespace string) (interface{}, error) {
	return nil, nil
}

// insert into service_metadata values(${serviceName}, ${namespace}, ${metadata})
func (repo *ServiceRepository) SaveMetadata(serviceName, namespace string, metadata interface{}) (*domain.Success, error) {
	return nil, nil
}

// delete from service_metadata where service_name=${serviceName} and namespace=${namespace}
func (repo *ServiceRepository) DeleteMetadata(serviceName, namespace string) (*domain.Success, error) {
	return nil, nil
}

// select * from service_roles where service_name=${serviceName} and role_name=${roleName}
func (repo *ServiceRepository) ExistsRole(serviceName, roleName string) bool {
	return false
}

// select * from service_roles where service_name=${serviceName}
func (repo *ServiceRepository) RoleList(serviceName string) (*domain.Roles, error) {
	return nil, nil
}

// select * from service_roles where service_name=${serviceName} and role_name=${roleName}
func (repo *ServiceRepository) Role(serviceName, roleName string) (*domain.Role, error) {
	return nil, fmt.Errorf("role not found")
}

// insert into service_roles values(${serviceName}, ${roleName}, ${Memo})
func (repo *ServiceRepository) SaveRole(serviceName string, r *domain.Role) error {
	return nil
}

// delete from service_roles where service_name=${serviceName} and role_name=${roleName}
func (repo *ServiceRepository) DeleteRole(serviceName, roleName string) error {
	return nil
}

// select * from role_metadata where service_name=${serviceName} and role_name=${roleName} and namespace=${namespace} limit=1
func (repo *ServiceRepository) ExistsRoleMetadata(serviceName, roleName, namespace string) bool {
	return true
}

// select namespace from role_metadata where service_name=${serviceName} and role_name=${roleName}
func (repo *ServiceRepository) RoleMetadataList(serviceName, roleName string) (*domain.RoleMetadataList, error) {
	return nil, nil
}

// select * from role_metadata where service_name=${serviceName} and role_name=${roleName} and namespace=${namespace}
func (repo *ServiceRepository) RoleMetadata(serviceName, roleName, namespace string) (interface{}, error) {
	return nil, nil
}

// insert into role_metadata values(${serviveName}, ${roleName}, ${namespace}, ${metadata})
func (repo *ServiceRepository) SaveRoleMetadata(serviceName, roleName, namespace string, metadata interface{}) (*domain.Success, error) {
	return nil, nil
}

// delete from role_metadata where service_name=${serviceName} and role_name=${roleName} and namespace=${namespace}
func (repo *ServiceRepository) DeleteRoleMetadata(serviceName, roleName, namespace string) (*domain.Success, error) {
	return nil, nil
}

// select * from service_metrics where service_name=${serviceName} and name=${metricName}
func (repo *ServiceRepository) ExistsMetric(serviceName, metricName string) bool {
	return true
}

// select distinct name from service_metrics where service_name=${serviceName}
func (repo *ServiceRepository) MetricNames(serviceName string) (*domain.ServiceMetricValueNames, error) {
	return &domain.ServiceMetricValueNames{}, nil
}

// select * from service_metrics where service_name=${serviceName} and name=${metricName}  and ${from} < from and to < ${to}
func (repo *ServiceRepository) MetricValues(serviceName, metricName string, from, to int) (*domain.ServiceMetricValues, error) {
	return &domain.ServiceMetricValues{}, nil
}

// insert into service_metrics values(${serviceName}, ${name}, ${time}, ${value})
func (repo *ServiceRepository) SaveMetricValues(serviceName string, values []domain.ServiceMetricValue) (*domain.Success, error) {
	return &domain.Success{Success: true}, nil
}
