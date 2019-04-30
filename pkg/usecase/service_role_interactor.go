package usecase

import (
	"regexp"

	"github.com/itsubaki/mackerel-api/pkg/domain"
)

type ServiceRoleInteractor struct {
	ServiceNameRule       *regexp.Regexp
	ServiceRoleNameRule   *regexp.Regexp
	ServiceRepository     ServiceRepository
	ServiceRoleRepository ServiceRoleRepository
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
