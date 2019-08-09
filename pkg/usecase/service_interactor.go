package usecase

import (
	"encoding/json"
	"errors"
	"fmt"
	"regexp"

	"github.com/itsubaki/mackerel-api/pkg/domain"
)

type ServiceInteractor struct {
	NameRule          *regexp.Regexp
	RoleNameRule      *regexp.Regexp
	ServiceRepository ServiceRepository
}

func (s *ServiceInteractor) List(orgID string) (*domain.Services, error) {
	return s.ServiceRepository.List(orgID)
}

func (s *ServiceInteractor) Save(orgID string, service *domain.Service) (*domain.Service, error) {
	if !s.NameRule.Match([]byte(service.Name)) {
		return nil, &InvalidServiceName{}
	}

	if err := s.ServiceRepository.Save(orgID, service); err != nil {
		return nil, err
	}

	return &domain.Service{
		Name:  service.Name,
		Memo:  service.Memo,
		Roles: []string{},
	}, nil
}

func (s *ServiceInteractor) Delete(orgID, serviceName string) (*domain.Service, error) {
	if !s.ServiceRepository.Exists(orgID, serviceName) {
		return nil, &ServiceNotFound{Err{errors.New("the Service corresponding to <serviceName> can't be found")}}
	}

	service, err := s.ServiceRepository.Service(orgID, serviceName)
	if err != nil {
		return nil, err
	}

	if err := s.ServiceRepository.Delete(orgID, serviceName); err != nil {
		return nil, err
	}

	return service, nil
}

func (s *ServiceInteractor) RoleList(orgID, serviceName string) (*domain.Roles, error) {
	if !s.ServiceRepository.Exists(orgID, serviceName) {
		return nil, &ServiceNotFound{Err{errors.New("the Service corresponding to <serviceName> can't be found")}}
	}

	list, err := s.ServiceRepository.RoleList(orgID, serviceName)
	if err != nil {
		return nil, err
	}

	return list, nil
}

func (s *ServiceInteractor) SaveRole(orgID, serviceName string, role *domain.Role) (*domain.Role, error) {
	if !s.ServiceRepository.Exists(orgID, serviceName) {
		return nil, &ServiceNotFound{Err{errors.New("the Service corresponding to <serviceName> can't be found")}}
	}

	if !s.RoleNameRule.Match([]byte(role.Name)) {
		return nil, &InvalidRoleName{}
	}

	if err := s.ServiceRepository.SaveRole(orgID, serviceName, role); err != nil {
		return nil, err
	}

	return role, nil
}

func (s *ServiceInteractor) DeleteRole(orgID, serviceName, roleName string) (*domain.Role, error) {
	if !s.ServiceRepository.Exists(orgID, serviceName) {
		return nil, &ServiceNotFound{Err{errors.New("the Service corresponding to <serviceName> can't be found")}}
	}

	if !s.ServiceRepository.ExistsRole(orgID, serviceName, roleName) {
		return nil, &RoleNotFound{Err{errors.New("the role corresponding to <roleName> can't be found")}}
	}

	r, err := s.ServiceRepository.Role(orgID, serviceName, roleName)
	if err != nil {
		return nil, err
	}

	if err := s.ServiceRepository.DeleteRole(orgID, serviceName, roleName); err != nil {
		return nil, err
	}

	return r, nil
}

func (s *ServiceInteractor) MetadataList(orgID, serviceName string) (*domain.ServiceMetadataList, error) {
	if !s.ServiceRepository.Exists(orgID, serviceName) {
		return nil, &ServiceNotFound{Err{errors.New("the Service corresponding to <serviceName> can't be found")}}
	}

	return s.ServiceRepository.MetadataList(orgID, serviceName)
}

func (s *ServiceInteractor) Metadata(orgID, serviceName, namespace string) (interface{}, error) {
	if !s.ServiceRepository.Exists(orgID, serviceName) {
		return nil, &ServiceNotFound{Err{errors.New("the Service corresponding to <serviceName> can't be found")}}
	}

	if !s.ServiceRepository.ExistsMetadata(orgID, serviceName, namespace) {
		return nil, &ServiceMetadataNotFound{Err{errors.New("metadata specified for the service does not exist")}}
	}

	return s.ServiceRepository.Metadata(orgID, serviceName, namespace)
}

func (s *ServiceInteractor) SaveMetadata(orgID, serviceName, namespace string, metadata interface{}) (*domain.Success, error) {
	if !s.ServiceRepository.Exists(orgID, serviceName) {
		return &domain.Success{Success: false}, &ServiceNotFound{Err{errors.New("the Service corresponding to <serviceName> can't be found")}}
	}

	meta, err := s.ServiceRepository.MetadataList(orgID, serviceName)
	if err != nil {
		return &domain.Success{Success: false}, fmt.Errorf("get metadata list: %v", err)
	}

	if len(meta.Metadata) > 50 {
		return &domain.Success{Success: false}, &MetadataLimitExceeded{Err{errors.New("trying to register while exceeding the limit of metadata per service (50 per service)")}}
	}

	b, err := json.Marshal(metadata)
	if err != nil {
		return &domain.Success{Success: false}, fmt.Errorf("marshal: %v", err)
	}

	if len(b) > 100000 {
		return &domain.Success{Success: false}, &MetadataTooLarge{Err{errors.New("the metadata exceeds 100KB")}}
	}

	res, err := s.ServiceRepository.SaveMetadata(orgID, serviceName, namespace, metadata)
	if err != nil {
		return res, fmt.Errorf("save metadata: %v", err)
	}

	return res, nil
}

func (s *ServiceInteractor) DeleteMetadata(orgID, serviceName, namespace string) (*domain.Success, error) {
	if !s.ServiceRepository.Exists(orgID, serviceName) {
		return &domain.Success{Success: false}, &ServiceNotFound{Err{errors.New("the Service corresponding to <serviceName> can't be found")}}
	}

	if !s.ServiceRepository.ExistsMetadata(orgID, serviceName, namespace) {
		return &domain.Success{Success: false}, &ServiceMetadataNotFound{Err{errors.New("the specified metadata does not exist for the host")}}
	}

	res, err := s.ServiceRepository.DeleteMetadata(orgID, serviceName, namespace)
	if err != nil {
		return res, fmt.Errorf("delete metadata: %v", err)
	}

	return res, nil
}

func (s *ServiceInteractor) RoleMetadataList(orgID, serviceName, roleName string) (*domain.RoleMetadataList, error) {
	if !s.ServiceRepository.Exists(orgID, serviceName) {
		return nil, &ServiceNotFound{Err{errors.New("the service does not exist")}}
	}

	if !s.ServiceRepository.ExistsRole(orgID, serviceName, roleName) {
		return nil, &RoleNotFound{Err{errors.New("the role does not exist")}}
	}

	return s.ServiceRepository.RoleMetadataList(orgID, serviceName, roleName)
}

func (s *ServiceInteractor) RoleMetadata(orgID, serviceName, roleName, namespace string) (interface{}, error) {
	if !s.ServiceRepository.Exists(orgID, serviceName) {
		return nil, &ServiceNotFound{Err{errors.New("the service does not exist")}}
	}

	if !s.ServiceRepository.ExistsRole(orgID, serviceName, roleName) {
		return nil, &RoleNotFound{Err{errors.New("the role does not exist")}}
	}

	if !s.ServiceRepository.ExistsRoleMetadata(orgID, serviceName, roleName, namespace) {
		return nil, &RoleMetadataNotFound{Err{errors.New("the metadata specified for the role does not exist")}}
	}

	return s.ServiceRepository.RoleMetadata(orgID, serviceName, roleName, namespace)
}

func (s *ServiceInteractor) SaveRoleMetadata(orgID, serviceName, roleName, namespace string, metadata interface{}) (*domain.Success, error) {
	if !s.ServiceRepository.Exists(orgID, serviceName) {
		return &domain.Success{Success: false}, &ServiceNotFound{Err{errors.New("the service does not exist")}}
	}

	if !s.ServiceRepository.ExistsRole(orgID, serviceName, roleName) {
		return &domain.Success{Success: false}, &RoleNotFound{Err{errors.New("the role does not exist")}}
	}

	meta, err := s.ServiceRepository.RoleMetadataList(orgID, serviceName, roleName)
	if err != nil {
		return &domain.Success{Success: false}, fmt.Errorf("get role metadata list: %v", err)
	}

	if len(meta.Metadata) > 50 {
		return &domain.Success{Success: false}, &MetadataLimitExceeded{Err{errors.New("trying to register while exceeding the limit of metadata per service (50 per service)")}}
	}

	b, err := json.Marshal(metadata)
	if err != nil {
		return &domain.Success{Success: false}, fmt.Errorf("marshal: %v", err)
	}

	if len(b) > 100000 {
		return &domain.Success{Success: false}, &MetadataTooLarge{Err{errors.New("the metadata exceeds 100KB")}}
	}

	res, err := s.ServiceRepository.SaveRoleMetadata(orgID, serviceName, roleName, namespace, metadata)
	if err != nil {
		return res, fmt.Errorf("save role metadata: %v", err)
	}

	return res, nil
}

func (s *ServiceInteractor) DeleteRoleMetadata(orgID, serviceName, roleName, namespace string) (*domain.Success, error) {
	if !s.ServiceRepository.Exists(orgID, serviceName) {
		return &domain.Success{Success: false}, &ServiceNotFound{Err{errors.New("the service does not exist")}}
	}

	if !s.ServiceRepository.ExistsRole(orgID, serviceName, roleName) {
		return &domain.Success{Success: false}, &RoleNotFound{Err{errors.New("the role does not exist")}}
	}

	if !s.ServiceRepository.ExistsRoleMetadata(orgID, serviceName, roleName, namespace) {
		return &domain.Success{Success: false}, &RoleMetadataNotFound{Err{errors.New("the metadata specified for the role does not exist")}}
	}

	res, err := s.ServiceRepository.DeleteRoleMetadata(orgID, serviceName, roleName, namespace)
	if err != nil {
		return res, fmt.Errorf("delete role metadata: %v", err)
	}

	return res, nil
}

func (s *ServiceInteractor) MetricNames(orgID, serviceName string) (*domain.ServiceMetricValueNames, error) {
	if !s.ServiceRepository.Exists(orgID, serviceName) {
		return nil, &ServiceNotFound{Err{errors.New("the Service corresponding to <serviceName> can't be found")}}
	}

	return s.ServiceRepository.MetricNames(orgID, serviceName)
}

func (s *ServiceInteractor) MetricValues(orgID, serviceName, metricName string, from, to int64) (*domain.ServiceMetricValues, error) {
	if !s.ServiceRepository.Exists(orgID, serviceName) {
		return nil, &ServiceNotFound{Err{errors.New("the service does not exist")}}
	}

	if !s.ServiceRepository.ExistsMetric(orgID, serviceName, metricName) {
		return nil, &ServiceMetricNotFound{Err{errors.New("the metric does not exist")}}
	}

	return s.ServiceRepository.MetricValues(orgID, serviceName, metricName, from, to)
}

func (s *ServiceInteractor) SaveMetricValues(orgID, serviceName string, values []domain.ServiceMetricValue) (*domain.Success, error) {
	// TODO
	// When the number of requests per minute is exceeded. Correct this by setting the posting frequency to a 1 minute interval, or posting multiple metrics at once, etc.
	res, err := s.ServiceRepository.SaveMetricValues(orgID, serviceName, values)
	if err != nil {
		return res, fmt.Errorf("save metric values: %v", err)
	}

	return res, nil
}
