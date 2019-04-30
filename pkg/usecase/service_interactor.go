package usecase

import (
	"github.com/itsubaki/mackerel-api/pkg/domain"
	"github.com/itsubaki/mackerel-api/pkg/interfaces/database"
)

type ServiceInteractor struct {
	ServiceRepository     *database.ServiceRepository
	ServiceRoleRepository *database.ServiceRoleRepository
}

func (s *ServiceInteractor) List() (domain.Services, error) {
	return s.ServiceRepository.FindAll()
}

func (s *ServiceInteractor) Delete(serviceName string) (*domain.Service, error) {
	service, err := s.ServiceRepository.FindByName(serviceName)
	if err != nil {
		return nil, &ServiceNotFound{}
	}

	roles, err := s.ServiceRoleRepository.FindAll(serviceName)
	if err != nil {
		return nil, err
	}

	for i := range roles {
		if err := s.ServiceRoleRepository.Delete(roles[i].ServiceName, roles[i].Name); err != nil {
			return nil, err
		}
	}

	if err := s.ServiceRepository.Delete(serviceName); err != nil {
		return nil, err
	}

	return &domain.Service{
		Name:  service.Name,
		Memo:  service.Memo,
		Roles: roles.Array(),
	}, nil
}
