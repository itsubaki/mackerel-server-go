package usecase

import (
	"encoding/json"
	"errors"
	"fmt"
	"regexp"

	"github.com/itsubaki/mackerel-api/pkg/domain"
)

type ServiceInteractor struct {
	NameRule                *regexp.Regexp
	RoleNameRule            *regexp.Regexp
	ServiceRepository       ServiceRepository
	ServiceMetaRepository   ServiceMetaRepository
	ServiceMetricRepository ServiceMetricRepository
	RoleRepository          RoleRepository
	RoleMetaRepository      RoleMetaRepository
}

func (s *ServiceInteractor) List(orgID string) (*domain.Services, error) {
	roles, err := s.RoleRepository.List(orgID)
	if err != nil {
		return nil, fmt.Errorf("list roles: %v", err)
	}

	services, err := s.ServiceRepository.List(orgID)
	if err != nil {
		return nil, fmt.Errorf("list services: %v", err)
	}

	for i := range services.Services {
		for k, v := range roles {
			if services.Services[i].Name != k {
				continue
			}

			services.Services[i].Roles = v
		}
	}

	return services, nil
}

func (s *ServiceInteractor) Save(orgID string, service *domain.Service) (*domain.Service, error) {
	if service.Roles == nil {
		service.Roles = make([]string, 0)
	}

	if !s.NameRule.Match([]byte(service.Name)) {
		return nil, &InvalidServiceName{}
	}

	if err := s.ServiceRepository.Save(orgID, service); err != nil {
		return nil, err
	}

	for i := range service.Roles {
		if !s.RoleRepository.Exists(orgID, service.Name, service.Roles[i]) {
			continue
		}

		if err := s.RoleRepository.Save(orgID, service.Name, &domain.Role{
			OrgID:       orgID,
			ServiceName: service.Name,
			Name:        service.Roles[i],
		}); err != nil {
			return nil, err
		}
	}

	return &domain.Service{
		Name:  service.Name,
		Memo:  service.Memo,
		Roles: service.Roles,
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

	roles, err := s.RoleRepository.ListWith(orgID, serviceName)
	if err != nil {
		return nil, err
	}
	service.Roles = roles.Array()

	if err := s.ServiceRepository.Delete(orgID, serviceName); err != nil {
		return nil, err
	}

	return service, nil
}

func (s *ServiceInteractor) ListRole(orgID, serviceName string) (*domain.Roles, error) {
	if !s.ServiceRepository.Exists(orgID, serviceName) {
		return nil, &ServiceNotFound{Err{errors.New("the Service corresponding to <serviceName> can't be found")}}
	}

	list, err := s.RoleRepository.ListWith(orgID, serviceName)
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

	if err := s.RoleRepository.Save(orgID, serviceName, role); err != nil {
		return nil, err
	}

	return role, nil
}

func (s *ServiceInteractor) DeleteRole(orgID, serviceName, roleName string) (*domain.Role, error) {
	if !s.ServiceRepository.Exists(orgID, serviceName) {
		return nil, &ServiceNotFound{Err{errors.New("the Service corresponding to <serviceName> can't be found")}}
	}

	if !s.RoleRepository.Exists(orgID, serviceName, roleName) {
		return nil, &RoleNotFound{Err{errors.New("the role corresponding to <roleName> can't be found")}}
	}

	r, err := s.RoleRepository.Role(orgID, serviceName, roleName)
	if err != nil {
		return nil, err
	}

	if err := s.RoleRepository.Delete(orgID, serviceName, roleName); err != nil {
		return nil, err
	}

	return r, nil
}

func (s *ServiceInteractor) Metadata(orgID, serviceName, namespace string) (interface{}, error) {
	if !s.ServiceRepository.Exists(orgID, serviceName) {
		return nil, &ServiceNotFound{Err{errors.New("the Service corresponding to <serviceName> can't be found")}}
	}

	if !s.ServiceMetaRepository.Exists(orgID, serviceName, namespace) {
		return nil, &ServiceMetadataNotFound{Err{errors.New("metadata specified for the service does not exist")}}
	}

	return s.ServiceMetaRepository.Metadata(orgID, serviceName, namespace)
}

func (s *ServiceInteractor) ListMetadata(orgID, serviceName string) (*domain.ServiceMetadataList, error) {
	if !s.ServiceRepository.Exists(orgID, serviceName) {
		return nil, &ServiceNotFound{Err{errors.New("the Service corresponding to <serviceName> can't be found")}}
	}

	return s.ServiceMetaRepository.List(orgID, serviceName)
}

func (s *ServiceInteractor) SaveMetadata(orgID, serviceName, namespace string, metadata interface{}) (*domain.Success, error) {
	if !s.ServiceRepository.Exists(orgID, serviceName) {
		return &domain.Success{Success: false}, &ServiceNotFound{Err{errors.New("the Service corresponding to <serviceName> can't be found")}}
	}

	meta, err := s.ServiceMetaRepository.List(orgID, serviceName)
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

	res, err := s.ServiceMetaRepository.Save(orgID, serviceName, namespace, metadata)
	if err != nil {
		return res, fmt.Errorf("save metadata: %v", err)
	}

	return res, nil
}

func (s *ServiceInteractor) DeleteMetadata(orgID, serviceName, namespace string) (*domain.Success, error) {
	if !s.ServiceRepository.Exists(orgID, serviceName) {
		return &domain.Success{Success: false}, &ServiceNotFound{Err{errors.New("the Service corresponding to <serviceName> can't be found")}}
	}

	if !s.ServiceMetaRepository.Exists(orgID, serviceName, namespace) {
		return &domain.Success{Success: false}, &ServiceMetadataNotFound{Err{errors.New("the specified metadata does not exist for the host")}}
	}

	res, err := s.ServiceMetaRepository.Delete(orgID, serviceName, namespace)
	if err != nil {
		return res, fmt.Errorf("delete metadata: %v", err)
	}

	return res, nil
}

func (s *ServiceInteractor) RoleMetadata(orgID, serviceName, roleName, namespace string) (interface{}, error) {
	if !s.ServiceRepository.Exists(orgID, serviceName) {
		return nil, &ServiceNotFound{Err{errors.New("the service does not exist")}}
	}

	if !s.RoleRepository.Exists(orgID, serviceName, roleName) {
		return nil, &RoleNotFound{Err{errors.New("the role does not exist")}}
	}

	if !s.RoleMetaRepository.Exists(orgID, serviceName, roleName, namespace) {
		return nil, &RoleMetadataNotFound{Err{errors.New("the metadata specified for the role does not exist")}}
	}

	return s.RoleMetaRepository.Metadata(orgID, serviceName, roleName, namespace)
}

func (s *ServiceInteractor) ListRoleMetadata(orgID, serviceName, roleName string) (*domain.RoleMetadataList, error) {
	if !s.ServiceRepository.Exists(orgID, serviceName) {
		return nil, &ServiceNotFound{Err{errors.New("the service does not exist")}}
	}

	if !s.RoleRepository.Exists(orgID, serviceName, roleName) {
		return nil, &RoleNotFound{Err{errors.New("the role does not exist")}}
	}

	return s.RoleMetaRepository.List(orgID, serviceName, roleName)
}

func (s *ServiceInteractor) SaveRoleMetadata(orgID, serviceName, roleName, namespace string, metadata interface{}) (*domain.Success, error) {
	if !s.ServiceRepository.Exists(orgID, serviceName) {
		return &domain.Success{Success: false}, &ServiceNotFound{Err{errors.New("the service does not exist")}}
	}

	if !s.RoleRepository.Exists(orgID, serviceName, roleName) {
		return &domain.Success{Success: false}, &RoleNotFound{Err{errors.New("the role does not exist")}}
	}

	meta, err := s.RoleMetaRepository.List(orgID, serviceName, roleName)
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

	res, err := s.RoleMetaRepository.Save(orgID, serviceName, roleName, namespace, metadata)
	if err != nil {
		return res, fmt.Errorf("save role metadata: %v", err)
	}

	return res, nil
}

func (s *ServiceInteractor) DeleteRoleMetadata(orgID, serviceName, roleName, namespace string) (*domain.Success, error) {
	if !s.ServiceRepository.Exists(orgID, serviceName) {
		return &domain.Success{Success: false}, &ServiceNotFound{Err{errors.New("the service does not exist")}}
	}

	if !s.RoleRepository.Exists(orgID, serviceName, roleName) {
		return &domain.Success{Success: false}, &RoleNotFound{Err{errors.New("the role does not exist")}}
	}

	if !s.RoleMetaRepository.Exists(orgID, serviceName, roleName, namespace) {
		return &domain.Success{Success: false}, &RoleMetadataNotFound{Err{errors.New("the metadata specified for the role does not exist")}}
	}

	res, err := s.RoleMetaRepository.Delete(orgID, serviceName, roleName, namespace)
	if err != nil {
		return res, fmt.Errorf("delete role metadata: %v", err)
	}

	return res, nil
}

func (s *ServiceInteractor) MetricNames(orgID, serviceName string) (*domain.ServiceMetricValueNames, error) {
	if !s.ServiceRepository.Exists(orgID, serviceName) {
		return nil, &ServiceNotFound{Err{errors.New("the Service corresponding to <serviceName> can't be found")}}
	}

	return s.ServiceMetricRepository.Names(orgID, serviceName)
}

func (s *ServiceInteractor) MetricValues(orgID, serviceName, metricName string, from, to int64) (*domain.ServiceMetricValues, error) {
	if !s.ServiceRepository.Exists(orgID, serviceName) {
		return nil, &ServiceNotFound{Err{errors.New("the service does not exist")}}
	}

	if !s.ServiceMetricRepository.Exists(orgID, serviceName, metricName) {
		return nil, &ServiceMetricNotFound{Err{errors.New("the metric does not exist")}}
	}

	return s.ServiceMetricRepository.Values(orgID, serviceName, metricName, from, to)
}

func (s *ServiceInteractor) SaveMetricValues(orgID, serviceName string, values []domain.ServiceMetricValue) (*domain.Success, error) {
	//TODO &ServiceMetricLimitExceeded{Err{errors.New("the number of requests per minute is exceeded")}}
	res, err := s.ServiceMetricRepository.Save(orgID, serviceName, values)
	if err != nil {
		return res, fmt.Errorf("save metric values: %v", err)
	}

	return res, nil
}
