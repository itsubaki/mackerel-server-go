package usecase

import (
	"regexp"

	"github.com/itsubaki/mackerel-api/pkg/domain"
)

type ServiceInteractor struct {
	ServiceNameRule       *regexp.Regexp
	ServiceRoleNameRule   *regexp.Regexp
	ServiceRepository     ServiceRepository
	ServiceRoleRepository ServiceRoleRepository
}

func (s *ServiceInteractor) FindAll() (domain.Services, error) {
	return s.ServiceRepository.FindAll()
}

func (s *ServiceInteractor) Save(service *domain.Service) (*domain.Service, error) {
	if !s.ServiceNameRule.Match([]byte(service.Name)) {
		return nil, &InvalidServiceName{}
	}

	if s.ServiceRepository.ExistsByName(service.Name) {
		return nil, &InvalidServiceName{}
	}

	if err := s.ServiceRepository.Save(*service); err != nil {
		return nil, err
	}

	return service, nil
}

func (s *ServiceInteractor) Delete(serviceName string) (*domain.Service, error) {
	service, err := s.ServiceRepository.FindByName(serviceName)
	if err != nil {
		return nil, &ServiceNotFound{}
	}

	if err := s.ServiceRoleRepository.DeleteAll(serviceName); err != nil {
		return nil, err
	}

	if err := s.ServiceRepository.Delete(serviceName); err != nil {
		return nil, err
	}

	return &service, nil
}
