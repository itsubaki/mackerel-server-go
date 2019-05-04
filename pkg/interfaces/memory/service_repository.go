package memory

import (
	"fmt"

	"github.com/itsubaki/mackerel-api/pkg/domain"
)

type ServiceRepository struct {
	Services            *domain.Services
	ServiceMetadata     *domain.ServiceMetadataList
	ServiceMetricValues *domain.ServiceMetricValues
	Roles               *domain.Roles
	RoleMetadataL       *domain.RoleMetadataList
}

func (repo *ServiceRepository) List() (*domain.Services, error) {
	return repo.Services, nil
}

func (repo *ServiceRepository) Exists(serviceName string) bool {
	for i := range repo.Services.Services {
		if repo.Services.Services[i].Name == serviceName {
			return true
		}
	}

	return false
}

func (repo *ServiceRepository) Service(serviceName string) (*domain.Service, error) {
	for i := range repo.Services.Services {
		if repo.Services.Services[i].Name != serviceName {
			continue
		}

		s := repo.Services.Services[i]
		if s.Roles == nil {
			s.Roles = []string{}
		}

		return &s, nil
	}

	return nil, fmt.Errorf("service not found")
}

func (repo *ServiceRepository) Save(s *domain.Service) error {
	repo.Services.Services = append(repo.Services.Services, *s)
	return nil
}

func (repo *ServiceRepository) Delete(serviceName string) error {
	services := []domain.Service{}
	for i := range repo.Services.Services {
		if repo.Services.Services[i].Name != serviceName {
			services = append(services, repo.Services.Services[i])
		}
	}
	repo.Services.Services = services

	roles := []domain.Role{}
	for i := range repo.Roles.Roles {
		if repo.Roles.Roles[i].ServiceName != serviceName {
			roles = append(roles, repo.Roles.Roles[i])
		}
	}
	repo.Roles.Roles = roles

	return nil
}

func (repo *ServiceRepository) ExistsMetadata(serviceName, namespace string) bool {
	for _, m := range repo.ServiceMetadata.Metadata {
		if m.ServiceName == serviceName && m.Namespace == namespace {
			return true
		}
	}

	return false
}

func (repo *ServiceRepository) MetadataList(serviceName string) (*domain.ServiceMetadataList, error) {
	list := []domain.ServiceMetadata{}
	for i := range repo.ServiceMetadata.Metadata {
		if repo.ServiceMetadata.Metadata[i].ServiceName == serviceName {
			list = append(list, repo.ServiceMetadata.Metadata[i])
		}
	}

	return &domain.ServiceMetadataList{
		Metadata: list,
	}, nil
}

func (repo *ServiceRepository) Metadata(serviceName, namespace string) (interface{}, error) {
	for i := range repo.ServiceMetadata.Metadata {
		if repo.ServiceMetadata.Metadata[i].ServiceName == serviceName && repo.ServiceMetadata.Metadata[i].Namespace == namespace {
			return repo.ServiceMetadata.Metadata[i].Metadata, nil
		}
	}

	return nil, fmt.Errorf("serviceName/namespace not found")
}

func (repo *ServiceRepository) SaveMetadata(serviceName, namespace string, metadata interface{}) (*domain.Success, error) {
	for i := range repo.ServiceMetadata.Metadata {
		if repo.ServiceMetadata.Metadata[i].ServiceName == serviceName && repo.ServiceMetadata.Metadata[i].Namespace == namespace {
			repo.ServiceMetadata.Metadata[i].Metadata = metadata
			return &domain.Success{Success: true}, nil
		}
	}

	repo.ServiceMetadata.Metadata = append(repo.ServiceMetadata.Metadata, domain.ServiceMetadata{
		ServiceName: serviceName,
		Namespace:   namespace,
		Metadata:    metadata,
	})

	return &domain.Success{Success: true}, nil
}

func (repo *ServiceRepository) DeleteMetadata(serviceName, namespace string) (*domain.Success, error) {
	list := []domain.ServiceMetadata{}
	for i := range repo.ServiceMetadata.Metadata {
		if repo.ServiceMetadata.Metadata[i].ServiceName == serviceName && repo.ServiceMetadata.Metadata[i].Namespace == namespace {
			continue
		}
		list = append(list, repo.ServiceMetadata.Metadata[i])
	}
	repo.ServiceMetadata.Metadata = list

	return &domain.Success{Success: true}, nil
}

func (repo *ServiceRepository) ExistsRole(serviceName, roleName string) bool {
	for i := range repo.Roles.Roles {
		if repo.Roles.Roles[i].ServiceName == serviceName && repo.Roles.Roles[i].Name == roleName {
			return true
		}
	}

	return false
}

func (repo *ServiceRepository) RoleList(serviceName string) (*domain.Roles, error) {
	list := &domain.Roles{
		Roles: []domain.Role{},
	}
	for i := range repo.Roles.Roles {
		if repo.Roles.Roles[i].ServiceName == serviceName {
			list.Roles = append(list.Roles, repo.Roles.Roles[i])
		}
	}

	return list, nil
}

func (repo *ServiceRepository) Role(serviceName, roleName string) (*domain.Role, error) {
	for i := range repo.Roles.Roles {
		if repo.Roles.Roles[i].ServiceName == serviceName && repo.Roles.Roles[i].Name == roleName {
			return &repo.Roles.Roles[i], nil
		}
	}

	return nil, fmt.Errorf("role not found")
}

func (repo *ServiceRepository) SaveRole(serviceName string, r *domain.Role) error {
	r.ServiceName = serviceName
	repo.Roles.Roles = append(repo.Roles.Roles, *r)
	return nil
}

func (repo *ServiceRepository) DeleteRole(serviceName, roleName string) error {
	list := &domain.Roles{
		Roles: []domain.Role{},
	}
	for i := range repo.Roles.Roles {
		if repo.Roles.Roles[i].ServiceName != serviceName || repo.Roles.Roles[i].Name != roleName {
			list.Roles = append(list.Roles, repo.Roles.Roles[i])
		}
	}

	repo.Roles = list

	return nil
}

func (repo *ServiceRepository) ExistsRoleMetadata(serviceName, roleName, namespace string) bool {
	for _, m := range repo.RoleMetadataL.Metadata {
		if m.ServiceName == serviceName && m.RoleName == roleName && m.Namespace == namespace {
			return true
		}
	}

	return false
}

func (repo *ServiceRepository) RoleMetadataList(serviceName, roleName string) (*domain.RoleMetadataList, error) {
	list := []domain.RoleMetadata{}
	for i := range repo.RoleMetadataL.Metadata {
		if repo.RoleMetadataL.Metadata[i].ServiceName == serviceName && repo.RoleMetadataL.Metadata[i].RoleName == roleName {
			list = append(list, repo.RoleMetadataL.Metadata[i])
		}
	}

	return &domain.RoleMetadataList{
		Metadata: list,
	}, nil
}

func (repo *ServiceRepository) RoleMetadata(serviceName, roleName, namespace string) (interface{}, error) {
	for i := range repo.RoleMetadataL.Metadata {
		if repo.RoleMetadataL.Metadata[i].ServiceName == serviceName &&
			repo.RoleMetadataL.Metadata[i].RoleName == roleName &&
			repo.RoleMetadataL.Metadata[i].Namespace == namespace {
			return repo.RoleMetadataL.Metadata[i].Metadata, nil
		}
	}

	return nil, fmt.Errorf("serviceName/roleName/namespace not found")
}

func (repo *ServiceRepository) SaveRoleMetadata(serviceName, roleName, namespace string, metadata interface{}) (*domain.Success, error) {
	for i := range repo.RoleMetadataL.Metadata {
		if repo.RoleMetadataL.Metadata[i].ServiceName == serviceName &&
			repo.RoleMetadataL.Metadata[i].RoleName == roleName &&
			repo.RoleMetadataL.Metadata[i].Namespace == namespace {
			repo.RoleMetadataL.Metadata[i].Metadata = metadata
			return &domain.Success{Success: true}, nil
		}
	}

	repo.RoleMetadataL.Metadata = append(repo.RoleMetadataL.Metadata, domain.RoleMetadata{
		ServiceName: serviceName,
		RoleName:    roleName,
		Namespace:   namespace,
		Metadata:    metadata,
	})

	return &domain.Success{Success: true}, nil
}

func (repo *ServiceRepository) DeleteRoleMetadata(serviceName, roleName, namespace string) (*domain.Success, error) {
	list := []domain.RoleMetadata{}
	for i := range repo.RoleMetadataL.Metadata {
		if repo.RoleMetadataL.Metadata[i].ServiceName == serviceName &&
			repo.RoleMetadataL.Metadata[i].RoleName == roleName &&
			repo.RoleMetadataL.Metadata[i].Namespace == namespace {
			continue
		}
		list = append(list, repo.RoleMetadataL.Metadata[i])
	}
	repo.RoleMetadataL.Metadata = list

	return &domain.Success{Success: true}, nil
}

func (repo *ServiceRepository) ExistsMetric(serviceName, metricName string) bool {
	for _, m := range repo.ServiceMetricValues.Metrics {
		if m.ServiceName == serviceName && m.Name == metricName {
			return true
		}
	}

	return false
}

func (repo *ServiceRepository) MetricNames(serviceName string) (*domain.ServiceMetricValueNames, error) {
	return &domain.ServiceMetricValueNames{
		Names: repo.ServiceMetricValues.MetricNames().Names,
	}, nil
}

func (repo *ServiceRepository) MetricValues(serviceName, metricName string, from, to int) (*domain.ServiceMetricValues, error) {
	metrics := []domain.ServiceMetricValue{}

	for i := range repo.ServiceMetricValues.Metrics {
		if repo.ServiceMetricValues.Metrics[i].ServiceName != serviceName {
			continue
		}
		if repo.ServiceMetricValues.Metrics[i].Name != metricName {
			continue
		}
		if from > repo.ServiceMetricValues.Metrics[i].Time {
			continue
		}
		if repo.ServiceMetricValues.Metrics[i].Time > to {
			continue
		}

		metrics = append(metrics, repo.ServiceMetricValues.Metrics[i])
	}

	return &domain.ServiceMetricValues{Metrics: metrics}, nil
}

func (repo *ServiceRepository) SaveMetricValues(serviceName string, values []domain.ServiceMetricValue) (*domain.Success, error) {
	for i := range values {
		values[i].ServiceName = serviceName
		repo.ServiceMetricValues.Metrics = append(repo.ServiceMetricValues.Metrics, values[i])
	}

	return &domain.Success{Success: true}, nil
}
