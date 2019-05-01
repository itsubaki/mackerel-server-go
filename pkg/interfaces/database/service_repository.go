package database

import (
	"fmt"

	"github.com/itsubaki/mackerel-api/pkg/domain"
)

type ServiceRepository struct {
	SQLHandler SQLHandler
}

func (repo *ServiceRepository) List() (*domain.Services, error) {
	return nil, nil
}

func (repo *ServiceRepository) Exists(serviceName string) bool {
	return false
}

func (repo *ServiceRepository) Service(serviceName string) (*domain.Service, error) {
	return nil, fmt.Errorf("service not found")
}

func (repo *ServiceRepository) Save(s *domain.Service) error {
	return nil
}

func (repo *ServiceRepository) Delete(serviceName string) error {
	return nil
}

func (repo *ServiceRepository) ExistsMetadata(serviceName, namespace string) bool {
	return true
}

func (repo *ServiceRepository) MetadataList(serviceName string) (*domain.ServiceMetadataList, error) {
	return nil, nil
}

func (repo *ServiceRepository) Metadata(serviceName, namespace string) (interface{}, error) {
	return nil, nil
}

func (repo *ServiceRepository) SaveMetadata(serviceName, namespace string, metadata interface{}) (*domain.Success, error) {
	return nil, nil
}

func (repo *ServiceRepository) DeleteMetadata(serviceName, namespace string) (*domain.Success, error) {
	return nil, nil
}

func (repo *ServiceRepository) ExistsRole(serviceName, roleName string) bool {
	return false
}

func (repo *ServiceRepository) RoleList(serviceName string) (*domain.Roles, error) {
	return nil, nil
}

func (repo *ServiceRepository) Role(serviceName, roleName string) (*domain.Role, error) {
	return nil, fmt.Errorf("role not found")
}

func (repo *ServiceRepository) SaveRole(serviceName string, r *domain.Role) error {
	return nil
}

func (repo *ServiceRepository) DeleteRole(serviceName, roleName string) error {
	return nil
}

func (repo *ServiceRepository) ExistsRoleMetadata(serviceName, roleName, namespace string) bool {
	return true
}

func (repo *ServiceRepository) RoleMetadataList(serviceName, roleName string) (*domain.RoleMetadataList, error) {
	return nil, nil
}

func (repo *ServiceRepository) RoleMetadata(serviceName, roleName, namespace string) (interface{}, error) {
	return nil, nil
}

func (repo *ServiceRepository) SaveRoleMetadata(serviceName, roleName, namespace string, metadata interface{}) (*domain.Success, error) {
	return nil, nil
}

func (repo *ServiceRepository) DeleteRoleMetadata(serviceName, roleName, namespace string) (*domain.Success, error) {
	return nil, nil
}

func (repo *ServiceRepository) ExistsMetric(serviceName, metricName string) bool {
	return true
}

func (repo *ServiceRepository) MetricNames(serviceName string) (*domain.ServiceMetricValueNames, error) {
	return &domain.ServiceMetricValueNames{}, nil
}

func (repo *ServiceRepository) MetricValues(serviceName, metricName string, from, to int) (*domain.ServiceMetricValues, error) {
	return &domain.ServiceMetricValues{}, nil
}

func (repo *ServiceRepository) SaveMetricValues(serviceName string, values []domain.ServiceMetricValue) (*domain.Success, error) {
	return &domain.Success{Success: true}, nil
}
