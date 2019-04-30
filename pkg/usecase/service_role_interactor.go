package usecase

import (
	"regexp"

	"github.com/itsubaki/mackerel-api/pkg/domain"
	"github.com/itsubaki/mackerel-api/pkg/interfaces/database"
)

func NewServiceRoleInteractor() *ServiceRoleInteractor {
	return &ServiceRoleInteractor{
		ServiceNameRule:       regexp.MustCompile(`^[a-zA-Z0-9]{1,1}[a-zA-Z0-9_-]{1,62}`),
		ServiceRoleNameRule:   regexp.MustCompile(`^[a-zA-Z0-9]{1,1}[a-zA-Z0-9_-]{1,62}`),
		ServiceRepository:     database.NewServiceRepository(),
		ServiceRoleRepository: database.NewServiceRoleRepository(),
	}
}

type ServiceRoleInteractor struct {
	ServiceNameRule       *regexp.Regexp
	ServiceRoleNameRule   *regexp.Regexp
	ServiceRepository     *database.ServiceRepository
	ServiceRoleRepository *database.ServiceRoleRepository
}

func (s *ServiceRoleInteractor) FindAll(serviceName string) (domain.ServiceRoles, error) {
	list, err := s.ServiceRoleRepository.FindAll(serviceName)
	if err != nil {
		return nil, &ServiceNotFound{}
	}

	return list, nil
}

func (s *ServiceRoleInteractor) Save(role *domain.ServiceRole) (*domain.ServiceRole, error) {
	if !s.ServiceNameRule.Match([]byte(role.ServiceName)) {
		return nil, &InvalidServiceName{}
	}

	if !s.ServiceRoleNameRule.Match([]byte(role.Name)) {
		return nil, &InvalidRoleName{}
	}

	if s.ServiceRoleRepository.ExistsByName(role.ServiceName, role.Name) {
		return nil, &InvalidServiceName{}
	}

	if err := s.ServiceRoleRepository.Save(*role); err != nil {
		return nil, err
	}

	return role, nil
}

func (s *ServiceRoleInteractor) Delete(serviceName, roleName string) (*domain.ServiceRole, error) {
	r, err := s.ServiceRoleRepository.FindByName(serviceName, roleName)
	if err != nil {
		return nil, &RoleNotFound{}
	}

	if err := s.ServiceRoleRepository.Delete(serviceName, roleName); err != nil {
		return nil, err
	}

	return &r, nil
}
