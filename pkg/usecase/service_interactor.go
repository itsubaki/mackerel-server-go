package usecase

import (
	"encoding/json"
	"errors"
	"regexp"

	"github.com/itsubaki/mackerel-api/pkg/domain"
)

type ServiceInteractor struct {
	NameRule          *regexp.Regexp
	RoleNameRule      *regexp.Regexp
	ServiceRepository ServiceRepository
}

func (s *ServiceInteractor) List(org string) (*domain.Services, error) {
	return s.ServiceRepository.List(org)
}

func (s *ServiceInteractor) Save(org string, service *domain.Service) (*domain.Service, error) {
	if !s.NameRule.Match([]byte(service.Name)) {
		return nil, &InvalidServiceName{}
	}

	if err := s.ServiceRepository.Save(org, service); err != nil {
		return nil, err
	}

	return &domain.Service{
		Name:  service.Name,
		Memo:  service.Memo,
		Roles: []string{},
	}, nil
}

func (s *ServiceInteractor) Delete(org, serviceName string) (*domain.Service, error) {
	if !s.ServiceRepository.Exists(org, serviceName) {
		return nil, &ServiceNotFound{Err{errors.New("the Service corresponding to <serviceName> can't be found")}}
	}

	service, err := s.ServiceRepository.Service(org, serviceName)
	if err != nil {
		return nil, err
	}

	if err := s.ServiceRepository.Delete(org, serviceName); err != nil {
		return nil, err
	}

	return service, nil
}

func (s *ServiceInteractor) RoleList(org, serviceName string) (*domain.Roles, error) {
	if !s.ServiceRepository.Exists(org, serviceName) {
		return nil, &ServiceNotFound{Err{errors.New("the Service corresponding to <serviceName> can't be found")}}
	}

	list, err := s.ServiceRepository.RoleList(org, serviceName)
	if err != nil {
		return nil, err
	}

	return list, nil
}

func (s *ServiceInteractor) SaveRole(org, serviceName string, role *domain.Role) (*domain.Role, error) {
	if !s.ServiceRepository.Exists(org, serviceName) {
		return nil, &ServiceNotFound{Err{errors.New("the Service corresponding to <serviceName> can't be found")}}
	}

	if !s.RoleNameRule.Match([]byte(role.Name)) {
		return nil, &InvalidRoleName{}
	}

	if err := s.ServiceRepository.SaveRole(org, serviceName, role); err != nil {
		return nil, err
	}

	return role, nil
}

func (s *ServiceInteractor) DeleteRole(org, serviceName, roleName string) (*domain.Role, error) {
	if !s.ServiceRepository.Exists(org, serviceName) {
		return nil, &ServiceNotFound{Err{errors.New("the Service corresponding to <serviceName> can't be found")}}
	}

	if !s.ServiceRepository.ExistsRole(org, serviceName, roleName) {
		return nil, &RoleNotFound{Err{errors.New("the role corresponding to <roleName> can't be found")}}
	}

	r, err := s.ServiceRepository.Role(org, serviceName, roleName)
	if err != nil {
		return nil, err
	}

	if err := s.ServiceRepository.DeleteRole(org, serviceName, roleName); err != nil {
		return nil, err
	}

	return r, nil
}

func (s *ServiceInteractor) MetadataList(org, serviceName string) (*domain.ServiceMetadataList, error) {
	if !s.ServiceRepository.Exists(org, serviceName) {
		return nil, &ServiceNotFound{Err{errors.New("the Service corresponding to <serviceName> can't be found")}}
	}

	return s.ServiceRepository.MetadataList(org, serviceName)
}

func (s *ServiceInteractor) Metadata(org, serviceName, namespace string) (interface{}, error) {
	if !s.ServiceRepository.Exists(org, serviceName) {
		return nil, &ServiceNotFound{Err{errors.New("the Service corresponding to <serviceName> can't be found")}}
	}

	if !s.ServiceRepository.ExistsMetadata(org, serviceName, namespace) {
		return nil, &ServiceMetadataNotFound{Err{errors.New("metadata specified for the service does not exist")}}
	}

	return s.ServiceRepository.Metadata(org, serviceName, namespace)
}

func (s *ServiceInteractor) SaveMetadata(org, serviceName, namespace string, metadata interface{}) (*domain.Success, error) {
	if !s.ServiceRepository.Exists(org, serviceName) {
		return nil, &ServiceNotFound{Err{errors.New("the Service corresponding to <serviceName> can't be found")}}
	}

	meta, err := s.ServiceRepository.MetadataList(org, serviceName)
	if err != nil {
		return nil, err
	}

	if len(meta.Metadata) > 50 {
		return nil, &MetadataLimitExceeded{Err{errors.New("trying to register while exceeding the limit of metadata per service (50 per service)")}}
	}

	b, err := json.Marshal(metadata)
	if err != nil {
		return nil, err
	}

	if len(b) > 100000 {
		return nil, &MetadataTooLarge{Err{errors.New("the metadata exceeds 100KB")}}
	}

	return s.ServiceRepository.SaveMetadata(org, serviceName, namespace, metadata)
}

func (s *ServiceInteractor) DeleteMetadata(org, serviceName, namespace string) (*domain.Success, error) {
	if !s.ServiceRepository.Exists(org, serviceName) {
		return nil, &ServiceNotFound{Err{errors.New("the Service corresponding to <serviceName> can't be found")}}
	}

	if !s.ServiceRepository.ExistsMetadata(org, serviceName, namespace) {
		return nil, &ServiceMetadataNotFound{Err{errors.New("the specified metadata does not exist for the host")}}
	}

	return s.ServiceRepository.DeleteMetadata(org, serviceName, namespace)
}

func (s *ServiceInteractor) RoleMetadataList(org, serviceName, roleName string) (*domain.RoleMetadataList, error) {
	if !s.ServiceRepository.Exists(org, serviceName) {
		return nil, &ServiceNotFound{Err{errors.New("the service does not exist")}}
	}

	if !s.ServiceRepository.ExistsRole(org, serviceName, roleName) {
		return nil, &RoleNotFound{Err{errors.New("the role does not exist")}}
	}

	return s.ServiceRepository.RoleMetadataList(org, serviceName, roleName)
}

func (s *ServiceInteractor) RoleMetadata(org, serviceName, roleName, namespace string) (interface{}, error) {
	if !s.ServiceRepository.Exists(org, serviceName) {
		return nil, &ServiceNotFound{Err{errors.New("the service does not exist")}}
	}

	if !s.ServiceRepository.ExistsRole(org, serviceName, roleName) {
		return nil, &RoleNotFound{Err{errors.New("the role does not exist")}}
	}

	if !s.ServiceRepository.ExistsRoleMetadata(org, serviceName, roleName, namespace) {
		return nil, &RoleMetadataNotFound{Err{errors.New("the metadata specified for the role does not exist")}}
	}

	return s.ServiceRepository.RoleMetadata(org, serviceName, roleName, namespace)
}

func (s *ServiceInteractor) SaveRoleMetadata(org, serviceName, roleName, namespace string, metadata interface{}) (*domain.Success, error) {
	if !s.ServiceRepository.Exists(org, serviceName) {
		return nil, &ServiceNotFound{Err{errors.New("the service does not exist")}}
	}

	if !s.ServiceRepository.ExistsRole(org, serviceName, roleName) {
		return nil, &RoleNotFound{Err{errors.New("the role does not exist")}}
	}

	meta, err := s.ServiceRepository.RoleMetadataList(org, serviceName, roleName)
	if err != nil {
		return nil, err
	}

	if len(meta.Metadata) > 50 {
		return nil, &MetadataLimitExceeded{Err{errors.New("trying to register while exceeding the limit of metadata per service (50 per service)")}}
	}

	b, err := json.Marshal(metadata)
	if err != nil {
		return nil, err
	}

	if len(b) > 100000 {
		return nil, &MetadataTooLarge{Err{errors.New("the metadata exceeds 100KB")}}
	}

	return s.ServiceRepository.SaveRoleMetadata(org, serviceName, roleName, namespace, metadata)
}

func (s *ServiceInteractor) DeleteRoleMetadata(org, serviceName, roleName, namespace string) (*domain.Success, error) {
	if !s.ServiceRepository.Exists(org, serviceName) {
		return nil, &ServiceNotFound{Err{errors.New("the service does not exist")}}
	}

	if !s.ServiceRepository.ExistsRole(org, serviceName, roleName) {
		return nil, &RoleNotFound{Err{errors.New("the role does not exist")}}
	}

	if !s.ServiceRepository.ExistsRoleMetadata(org, serviceName, roleName, namespace) {
		return nil, &RoleMetadataNotFound{Err{errors.New("the metadata specified for the role does not exist")}}
	}

	return s.ServiceRepository.DeleteRoleMetadata(org, serviceName, roleName, namespace)
}

func (s *ServiceInteractor) MetricNames(org, serviceName string) (*domain.ServiceMetricValueNames, error) {
	if !s.ServiceRepository.Exists(org, serviceName) {
		return nil, &ServiceNotFound{Err{errors.New("the Service corresponding to <serviceName> can't be found")}}
	}

	return s.ServiceRepository.MetricNames(org, serviceName)
}

func (s *ServiceInteractor) MetricValues(org, serviceName, metricName string, from, to int64) (*domain.ServiceMetricValues, error) {
	if !s.ServiceRepository.Exists(org, serviceName) {
		return nil, &ServiceNotFound{Err{errors.New("the service does not exist")}}
	}

	if !s.ServiceRepository.ExistsMetric(org, serviceName, metricName) {
		return nil, &ServiceMetricNotFound{Err{errors.New("the metric does not exist")}}
	}

	return s.ServiceRepository.MetricValues(org, serviceName, metricName, from, to)
}

func (s *ServiceInteractor) SaveMetricValues(org, serviceName string, values []domain.ServiceMetricValue) (*domain.Success, error) {
	// TODO
	// When the number of requests per minute is exceeded. Correct this by setting the posting frequency to a 1 minute interval, or posting multiple metrics at once, etc.
	return s.ServiceRepository.SaveMetricValues(org, serviceName, values)
}
