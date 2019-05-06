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

func (s *ServiceInteractor) List() (*domain.Services, error) {
	return s.ServiceRepository.List()
}

func (s *ServiceInteractor) Save(service *domain.Service) (*domain.Service, error) {
	if !s.NameRule.Match([]byte(service.Name)) {
		return nil, &InvalidServiceName{}
	}

	if err := s.ServiceRepository.Save(service); err != nil {
		return nil, err
	}

	return service, nil
}

func (s *ServiceInteractor) Delete(serviceName string) (*domain.Service, error) {
	if !s.ServiceRepository.Exists(serviceName) {
		return nil, &ServiceNotFound{Err{errors.New("the Service corresponding to <serviceName> can't be found")}}
	}

	service, err := s.ServiceRepository.Service(serviceName)
	if err != nil {
		return nil, err
	}

	if err := s.ServiceRepository.Delete(serviceName); err != nil {
		return nil, err
	}

	return service, nil
}

func (s *ServiceInteractor) RoleList(serviceName string) (*domain.Roles, error) {
	if !s.ServiceRepository.Exists(serviceName) {
		return nil, &ServiceNotFound{Err{errors.New("the Service corresponding to <serviceName> can't be found")}}
	}

	list, err := s.ServiceRepository.RoleList(serviceName)
	if err != nil {
		return nil, err
	}

	return list, nil
}

func (s *ServiceInteractor) SaveRole(serviceName string, role *domain.Role) (*domain.Role, error) {
	if !s.ServiceRepository.Exists(serviceName) {
		return nil, &ServiceNotFound{Err{errors.New("the Service corresponding to <serviceName> can't be found")}}
	}

	if !s.RoleNameRule.Match([]byte(role.Name)) {
		return nil, &InvalidRoleName{}
	}

	if err := s.ServiceRepository.SaveRole(serviceName, role); err != nil {
		return nil, err
	}

	return role, nil
}

func (s *ServiceInteractor) DeleteRole(serviceName, roleName string) (*domain.Role, error) {
	if !s.ServiceRepository.Exists(serviceName) {
		return nil, &ServiceNotFound{Err{errors.New("the Service corresponding to <serviceName> can't be found")}}
	}

	if !s.ServiceRepository.ExistsRole(serviceName, roleName) {
		return nil, &RoleNotFound{Err{errors.New("the role corresponding to <roleName> can't be found")}}
	}

	r, err := s.ServiceRepository.Role(serviceName, roleName)
	if err != nil {
		return nil, err
	}

	if err := s.ServiceRepository.DeleteRole(serviceName, roleName); err != nil {
		return nil, err
	}

	return r, nil
}

func (s *ServiceInteractor) MetadataList(serviceName string) (*domain.ServiceMetadataList, error) {
	if !s.ServiceRepository.Exists(serviceName) {
		return nil, &ServiceNotFound{Err{errors.New("the Service corresponding to <serviceName> can't be found")}}
	}

	return s.ServiceRepository.MetadataList(serviceName)
}

func (s *ServiceInteractor) Metadata(serviceName, namespace string) (interface{}, error) {
	if !s.ServiceRepository.Exists(serviceName) {
		return nil, &ServiceNotFound{Err{errors.New("the Service corresponding to <serviceName> can't be found")}}
	}

	if !s.ServiceRepository.ExistsMetadata(serviceName, namespace) {
		return nil, &ServiceMetadataNotFound{Err{errors.New("metadata specified for the service does not exist")}}
	}

	return s.ServiceRepository.Metadata(serviceName, namespace)
}

func (s *ServiceInteractor) SaveMetadata(serviceName, namespace string, metadata interface{}) (*domain.Success, error) {
	if !s.ServiceRepository.Exists(serviceName) {
		return nil, &ServiceNotFound{Err{errors.New("the Service corresponding to <serviceName> can't be found")}}
	}

	meta, err := s.ServiceRepository.MetadataList(serviceName)
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

	return s.ServiceRepository.SaveMetadata(serviceName, namespace, metadata)
}

func (s *ServiceInteractor) DeleteMetadata(serviceName, namespace string) (*domain.Success, error) {
	if !s.ServiceRepository.Exists(serviceName) {
		return nil, &ServiceNotFound{Err{errors.New("the Service corresponding to <serviceName> can't be found")}}
	}

	if !s.ServiceRepository.ExistsMetadata(serviceName, namespace) {
		return nil, &ServiceMetadataNotFound{Err{errors.New("the specified metadata does not exist for the host")}}
	}

	return s.ServiceRepository.DeleteMetadata(serviceName, namespace)
}

func (s *ServiceInteractor) RoleMetadataList(serviceName, roleName string) (*domain.RoleMetadataList, error) {
	if !s.ServiceRepository.Exists(serviceName) {
		return nil, &ServiceNotFound{Err{errors.New("the service does not exist")}}
	}

	if !s.ServiceRepository.ExistsRole(serviceName, roleName) {
		return nil, &RoleNotFound{Err{errors.New("the role does not exist")}}
	}

	return s.ServiceRepository.RoleMetadataList(serviceName, roleName)
}

func (s *ServiceInteractor) RoleMetadata(serviceName, roleName, namespace string) (interface{}, error) {
	if !s.ServiceRepository.Exists(serviceName) {
		return nil, &ServiceNotFound{Err{errors.New("the service does not exist")}}
	}

	if !s.ServiceRepository.ExistsRole(serviceName, roleName) {
		return nil, &RoleNotFound{Err{errors.New("the role does not exist")}}
	}

	if !s.ServiceRepository.ExistsRoleMetadata(serviceName, roleName, namespace) {
		return nil, &RoleMetadataNotFound{Err{errors.New("the metadata specified for the role does not exist")}}
	}

	return s.ServiceRepository.RoleMetadata(serviceName, roleName, namespace)
}

func (s *ServiceInteractor) SaveRoleMetadata(serviceName, roleName, namespace string, metadata interface{}) (*domain.Success, error) {
	if !s.ServiceRepository.Exists(serviceName) {
		return nil, &ServiceNotFound{Err{errors.New("the service does not exist")}}
	}

	if !s.ServiceRepository.ExistsRole(serviceName, roleName) {
		return nil, &RoleNotFound{Err{errors.New("the role does not exist")}}
	}

	meta, err := s.ServiceRepository.RoleMetadataList(serviceName, roleName)
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

	return s.ServiceRepository.SaveRoleMetadata(serviceName, roleName, namespace, metadata)
}

func (s *ServiceInteractor) DeleteRoleMetadata(serviceName, roleName, namespace string) (*domain.Success, error) {
	if !s.ServiceRepository.Exists(serviceName) {
		return nil, &ServiceNotFound{Err{errors.New("the service does not exist")}}
	}

	if !s.ServiceRepository.ExistsRole(serviceName, roleName) {
		return nil, &RoleNotFound{Err{errors.New("the role does not exist")}}
	}

	if !s.ServiceRepository.ExistsRoleMetadata(serviceName, roleName, namespace) {
		return nil, &RoleMetadataNotFound{Err{errors.New("the metadata specified for the role does not exist")}}
	}

	return s.ServiceRepository.DeleteRoleMetadata(serviceName, roleName, namespace)
}

func (s *ServiceInteractor) MetricNames(serviceName string) (*domain.ServiceMetricValueNames, error) {
	if !s.ServiceRepository.Exists(serviceName) {
		return nil, &ServiceNotFound{Err{errors.New("the Service corresponding to <serviceName> can't be found")}}
	}

	return s.ServiceRepository.MetricNames(serviceName)
}

func (s *ServiceInteractor) MetricValues(serviceName, metricName string, from, to int64) (*domain.ServiceMetricValues, error) {
	if !s.ServiceRepository.Exists(serviceName) {
		return nil, &ServiceNotFound{Err{errors.New("the service does not exist")}}
	}

	if !s.ServiceRepository.ExistsMetric(serviceName, metricName) {
		return nil, &ServiceMetricNotFound{Err{errors.New("the metric does not exist")}}
	}

	return s.ServiceRepository.MetricValues(serviceName, metricName, from, to)
}

func (s *ServiceInteractor) SaveMetricValues(serviceName string, values []domain.ServiceMetricValue) (*domain.Success, error) {
	// TODO
	// When the number of requests per minute is exceeded. Correct this by setting the posting frequency to a 1 minute interval, or posting multiple metrics at once, etc.
	return s.ServiceRepository.SaveMetricValues(serviceName, values)
}
