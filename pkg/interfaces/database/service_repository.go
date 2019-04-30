package database

import (
	"fmt"

	"github.com/itsubaki/mackerel-api/pkg/domain"
)

type ServiceRepository struct {
	SQLHandler          SQLHandler
	Services            domain.Services
	ServiceMetadata     domain.ServiceMetadataList
	ServiceMetricValues domain.ServiceMetricValues
	Roles               domain.Roles
	RoleMetadata        domain.RoleMetadataList
}

func (repo *ServiceRepository) MetricNames(serviceName string) ([]string, error) {
	return []string{}, nil
}

func (repo *ServiceRepository) MetricValues(serviceName, metricName string, from, to int64) (domain.ServiceMetricValues, error) {
	list := domain.ServiceMetricValues{}
	for i := range repo.ServiceMetricValues {
		if repo.ServiceMetricValues[i].ServiceName != serviceName {
			continue
		}
		if repo.ServiceMetricValues[i].Name != metricName {
			continue
		}
		if from > repo.ServiceMetricValues[i].Time {
			continue
		}
		if repo.ServiceMetricValues[i].Time > to {
			continue
		}

		list = append(list, repo.ServiceMetricValues[i])
	}

	return list, nil
}

func (repo *ServiceRepository) SaveMetricValues(v domain.ServiceMetricValues) error {
	repo.ServiceMetricValues = append(repo.ServiceMetricValues, v...)
	return nil
}

func (repo *ServiceRepository) Role(serviceName, roleName string) (*domain.Role, error) {
	for i := range repo.Roles {
		if repo.Roles[i].ServiceName == serviceName && repo.Roles[i].Name == roleName {
			return &repo.Roles[i], nil
		}
	}

	return nil, fmt.Errorf("role not found")
}

func (repo *ServiceRepository) RoleList(serviceName string) (domain.Roles, error) {
	list := domain.Roles{}
	for i := range repo.Roles {
		if repo.Roles[i].ServiceName == serviceName {
			list = append(list, repo.Roles[i])
		}
	}

	return list, nil
}

func (repo *ServiceRepository) SaveRole(r domain.Role) error {
	repo.Roles = append(repo.Roles, r)
	return nil
}

func (repo *ServiceRepository) DeleteRole(serviceName, roleName string) error {
	list := domain.Roles{}
	for i := range repo.Roles {
		if repo.Roles[i].ServiceName != serviceName || repo.Roles[i].Name != roleName {
			list = append(list, repo.Roles[i])
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

func (repo *ServiceRepository) Service(serviceName string) (domain.Service, error) {
	for i := range repo.Services {
		if repo.Services[i].Name == serviceName {
			return repo.Services[i], nil
		}
	}

	return domain.Service{}, fmt.Errorf("service not found")
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
