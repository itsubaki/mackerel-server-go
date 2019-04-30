package database

import (
	"fmt"

	"github.com/itsubaki/mackerel-api/pkg/domain"
)

type ServiceRepository struct {
	Internal domain.Services
}

func NewServiceRepository() *ServiceRepository {
	return &ServiceRepository{
		Internal: domain.Services{},
	}
}

func (repo *ServiceRepository) ExistsByName(serviceName string) bool {
	for i := range repo.Internal {
		if repo.Internal[i].Name == serviceName {
			return true
		}
	}

	return false
}

func (repo *ServiceRepository) FindByName(serviceName string) (domain.Service, error) {
	for i := range repo.Internal {
		if repo.Internal[i].Name == serviceName {
			return repo.Internal[i], nil
		}
	}

	return domain.Service{}, fmt.Errorf("service not found")
}

func (repo *ServiceRepository) FindAll() (domain.Services, error) {
	return repo.Internal, nil
}

func (repo *ServiceRepository) Save(s domain.Service) error {
	repo.Internal = append(repo.Internal, s)
	return nil
}

func (repo *ServiceRepository) Delete(serviceName string) error {
	list := domain.Services{}
	for i := range repo.Internal {
		if repo.Internal[i].Name != serviceName {
			list = append(list, repo.Internal[i])
		}
	}
	repo.Internal = list
	return nil
}
