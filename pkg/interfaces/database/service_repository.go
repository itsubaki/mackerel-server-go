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

func (repo *ServiceRepository) MetricValues(serviceName, metricName string, from, to int64) (*domain.ServiceMetricValues, error) {
	list := &domain.ServiceMetricValues{
		Metrics: []domain.ServiceMetricValue{},
	}

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

		list.Metrics = append(list.Metrics, repo.ServiceMetricValues.Metrics[i])
	}

	return list, nil
}

func (repo *ServiceRepository) SaveMetricValues(v domain.ServiceMetricValues) error {
	repo.ServiceMetricValues.Metrics = append(repo.ServiceMetricValues.Metrics, v.Metrics...)
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

func (repo *ServiceRepository) SaveRole(r domain.Role) error {
	repo.Roles.Roles = append(repo.Roles.Roles, r)
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
	for i := range repo.Services {
		if repo.Services[i].Name == serviceName {
			return true
		}
	}

	return false
}

func (repo *ServiceRepository) Service(serviceName string) (*domain.Service, error) {
	for i := range repo.Services {
		if repo.Services[i].Name == serviceName {
			return &repo.Services[i], nil
		}
	}

	return nil, fmt.Errorf("service not found")
}

func (repo *ServiceRepository) List() (domain.Services, error) {
	return repo.Services, nil
}

func (repo *ServiceRepository) Save(s domain.Service) error {
	repo.Services = append(repo.Services, s)
	return nil
}

func (repo *ServiceRepository) Delete(serviceName string) error {
	services := domain.Services{}
	for i := range repo.Services {
		if repo.Services[i].Name != serviceName {
			services = append(services, repo.Services[i])
		}
	}
	repo.Services = services

	roles := domain.Roles{}
	for i := range repo.Roles {
		if repo.Roles[i].ServiceName != serviceName {
			roles = append(roles, repo.Roles[i])
		}
	}
	repo.Roles = roles

	return nil
}
