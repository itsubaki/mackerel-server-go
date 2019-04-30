package database

import (
	"fmt"

	"github.com/itsubaki/mackerel-api/pkg/domain"
)

type ServiceRoleRepository struct {
	Internal domain.ServiceRoles
}

func NewServiceRoleRepository() *ServiceRoleRepository {
	return &ServiceRoleRepository{
		Internal: domain.ServiceRoles{},
	}
}

func (repo *ServiceRoleRepository) ExistsByName(serviceName, roleName string) bool {
	for i := range repo.Internal {
		if repo.Internal[i].ServiceName == serviceName && repo.Internal[i].Name == roleName {
			return true
		}
	}

	return false
}

func (repo *ServiceRoleRepository) FindByName(serviceName, roleName string) (domain.ServiceRole, error) {
	for i := range repo.Internal {
		if repo.Internal[i].ServiceName == serviceName && repo.Internal[i].Name == roleName {
			return repo.Internal[i], nil
		}
	}

	return domain.ServiceRole{}, fmt.Errorf("role not found")
}

func (repo *ServiceRoleRepository) FindAll(serviceName string) (domain.ServiceRoles, error) {
	list := domain.ServiceRoles{}
	for i := range repo.Internal {
		if repo.Internal[i].ServiceName == serviceName {
			list = append(list, repo.Internal[i])
		}
	}

	return list, nil
}

func (repo *ServiceRoleRepository) Save(r domain.ServiceRole) error {
	repo.Internal = append(repo.Internal, r)
	return nil
}

func (repo *ServiceRoleRepository) Delete(serviceName, roleName string) error {
	list := domain.ServiceRoles{}
	for i := range repo.Internal {
		if repo.Internal[i].ServiceName != serviceName || repo.Internal[i].Name != roleName {
			list = append(list, repo.Internal[i])
		}
	}
	repo.Internal = list
	return nil
}
