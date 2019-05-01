package usecase

import (
	"regexp"

	"github.com/itsubaki/mackerel-api/pkg/domain"
)

type ServiceInteractor struct {
	ServiceNameRule     *regexp.Regexp
	ServiceRoleNameRule *regexp.Regexp
	ServiceRepository   ServiceRepository
}

func (s *ServiceInteractor) List() (domain.Services, error) {
	return s.ServiceRepository.List()
}

func (s *ServiceInteractor) Save(service *domain.Service) (*domain.Service, error) {
	if !s.ServiceNameRule.Match([]byte(service.Name)) {
		return nil, &InvalidServiceName{}
	}

	if s.ServiceRepository.Exists(service.Name) {
		return nil, &InvalidServiceName{}
	}

	if err := s.ServiceRepository.Save(*service); err != nil {
		return nil, err
	}

	return service, nil
}

func (s *ServiceInteractor) Delete(serviceName string) (*domain.Service, error) {
	service, err := s.ServiceRepository.Service(serviceName)
	if err != nil {
		return nil, &ServiceNotFound{}
	}

	if err := s.ServiceRepository.Delete(serviceName); err != nil {
		return nil, err
	}

	return service, nil
}

func (s *ServiceInteractor) RoleList(serviceName string) (domain.Roles, error) {
	list, err := s.ServiceRepository.RoleList(serviceName)
	if err != nil {
		return nil, &ServiceNotFound{}
	}

	return list, nil
}

func (s *ServiceInteractor) SaveRole(role *domain.Role) (*domain.Role, error) {
	if !s.ServiceNameRule.Match([]byte(role.ServiceName)) {
		return nil, &InvalidServiceName{}
	}

	if !s.ServiceRoleNameRule.Match([]byte(role.Name)) {
		return nil, &InvalidRoleName{}
	}

	if s.ServiceRepository.Exists(role.ServiceName) {
		return nil, &InvalidServiceName{}
	}

	if err := s.ServiceRepository.SaveRole(*role); err != nil {
		return nil, err
	}

	return role, nil
}

func (s *ServiceInteractor) DeleteRole(serviceName, roleName string) (*domain.Role, error) {
	r, err := s.ServiceRepository.Role(serviceName, roleName)
	if err != nil {
		return nil, &RoleNotFound{}
	}

	if err := s.ServiceRepository.DeleteRole(serviceName, roleName); err != nil {
		return nil, err
	}

	return r, nil
}

func (s *ServiceInteractor) MetadataList() (domain.ServiceMetadataList, error) {
	return domain.ServiceMetadataList{}, nil
}

func (s *ServiceInteractor) Metadata(serviceName, namespace string) (interface{}, error) {
	return nil, nil
}

func (s *ServiceInteractor) SaveMetadata(serviceName, namespace string, metadata interface{}) error {
	return nil
}

func (s *ServiceInteractor) DeleteMetadata(serviceName, namespace string) error {
	return nil
}

func (s *ServiceInteractor) RoleMetadata(serviceName, roleName, namespace string) (interface{}, error) {
	return nil, nil
}

func (s *ServiceInteractor) SaveRoleMetadata(serviceName, roleName, namespace string, metadata interface{}) error {
	return nil
}

func (s *ServiceInteractor) DeleteRoleMetadata(serviceName, roleName, namespace string) error {
	return nil
}

func (s *ServiceInteractor) RoleMetadataList(serviceName, roleName string) (domain.RoleMetadataList, error) {
	return domain.RoleMetadataList{}, nil
}

func (s *ServiceInteractor) SaveMetricValues(values domain.ServiceMetricValues) error {
	return s.ServiceRepository.SaveMetricValues(values)
}

func (s *ServiceInteractor) MetricValues(serviceName, metricName string, from, to int64) (*domain.ServiceMetricValues, error) {
	return s.ServiceRepository.MetricValues(serviceName, metricName, from, to)
}

func (s *ServiceInteractor) MetricNames(serviceName string) (*domain.ServiceMetricValueNames, error) {
	return s.ServiceRepository.MetricNames(serviceName)
}
