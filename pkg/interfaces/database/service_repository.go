package database

import (
	"fmt"

	"github.com/itsubaki/mackerel-api/pkg/domain"
)

type ServiceRepository struct {
	SQLHandler          SQLHandler
	Services            *domain.Services
	ServiceMetadata     *domain.ServiceMetadataList
	ServiceMetricValues *domain.ServiceMetricValues
	Roles               *domain.Roles
	RoleMetadata        *domain.RoleMetadataList
}

func (repo *ServiceRepository) MetricNames(serviceName string) (*domain.ServiceMetricValueNames, error) {
	return &domain.ServiceMetricValueNames{}, nil
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

func (repo *ServiceRepository) SaveMetricValues(serviceName string, values []domain.ServiceMetricValue) error {
	for i := range values {
		values[i].ServiceName = serviceName
		repo.ServiceMetricValues.Metrics = append(repo.ServiceMetricValues.Metrics, values[i])
	}

	return nil
}

func (repo *ServiceRepository) Role(serviceName, roleName string) (*domain.Role, error) {
	for i := range repo.Roles.Roles {
		if repo.Roles.Roles[i].ServiceName == serviceName && repo.Roles.Roles[i].Name == roleName {
			return &repo.Roles.Roles[i], nil
		}
	}

	return nil, fmt.Errorf("role not found")
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

func (repo *ServiceRepository) List() (*domain.Services, error) {
	return repo.Services, nil
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
