package usecase

import (
	"encoding/json"
	"errors"
	"fmt"
	"regexp"

	"github.com/itsubaki/mackerel-server-go/domain"
)

var (
	ErrServiceNotFound           = errors.New("service not found")
	ErrServiceMetadataNotFound   = errors.New("service metadata not found")
	ErrServiceMetricNotFound     = errors.New("service metric not found")
	ErrRoleNotFound              = errors.New("role not found")
	ErrRoleMetadataNotFound      = errors.New("role metadata not found")
	ErrRoleMetadataLimitExceeded = errors.New("role metadata limit exceeded")
	ErrMetadataLimitExceeded     = errors.New("metadata limit exceeded")
	ErrMetadataTooLarge          = errors.New("metadata too large")
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

func (intr *ServiceInteractor) List(orgID string) (*domain.Services, error) {
	roles, err := intr.RoleRepository.List(orgID)
	if err != nil {
		return nil, fmt.Errorf("list roles: %v", err)
	}

	services, err := intr.ServiceRepository.List(orgID)
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

func (intr *ServiceInteractor) Save(orgID string, service *domain.Service) (*domain.Service, error) {
	if service.Roles == nil {
		service.Roles = make([]string, 0)
	}

	if !intr.NameRule.Match([]byte(service.Name)) {
		return nil, &InvalidServiceName{}
	}

	if err := intr.ServiceRepository.Save(orgID, service); err != nil {
		return nil, err
	}

	for i := range service.Roles {
		if !intr.RoleRepository.Exists(orgID, service.Name, service.Roles[i]) {
			continue
		}

		if err := intr.RoleRepository.Save(orgID, service.Name, &domain.Role{
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

func (intr *ServiceInteractor) Delete(orgID, serviceName string) (*domain.Service, error) {
	if !intr.ServiceRepository.Exists(orgID, serviceName) {
		return nil, &ServiceNotFound{
			Err{
				Err: ErrServiceNotFound,
			},
		}
	}

	service, err := intr.ServiceRepository.Service(orgID, serviceName)
	if err != nil {
		return nil, err
	}

	roles, err := intr.RoleRepository.ListWith(orgID, serviceName)
	if err != nil {
		return nil, err
	}
	service.Roles = roles.Array()

	for _, r := range roles.Array() {
		if err := intr.RoleRepository.Delete(orgID, serviceName, r); err != nil {
			return nil, err
		}
	}

	if err := intr.ServiceRepository.Delete(orgID, serviceName); err != nil {
		return nil, err
	}

	return service, nil
}

func (intr *ServiceInteractor) ListRole(orgID, serviceName string) (*domain.Roles, error) {
	if !intr.ServiceRepository.Exists(orgID, serviceName) {
		return nil, &ServiceNotFound{
			Err{
				Err: ErrServiceNotFound,
			},
		}
	}

	list, err := intr.RoleRepository.ListWith(orgID, serviceName)
	if err != nil {
		return nil, err
	}

	return list, nil
}

func (intr *ServiceInteractor) SaveRole(orgID, serviceName string, role *domain.Role) (*domain.Role, error) {
	if !intr.ServiceRepository.Exists(orgID, serviceName) {
		return nil, &ServiceNotFound{
			Err{
				Err: ErrServiceNotFound,
			},
		}
	}

	if !intr.RoleNameRule.Match([]byte(role.Name)) {
		return nil, &InvalidRoleName{}
	}

	if err := intr.RoleRepository.Save(orgID, serviceName, role); err != nil {
		return nil, err
	}

	return role, nil
}

func (intr *ServiceInteractor) DeleteRole(orgID, serviceName, roleName string) (*domain.Role, error) {
	if !intr.ServiceRepository.Exists(orgID, serviceName) {
		return nil, &ServiceNotFound{
			Err{
				Err: ErrServiceNotFound,
			},
		}
	}

	if !intr.RoleRepository.Exists(orgID, serviceName, roleName) {
		return nil, &RoleNotFound{
			Err{
				Err: ErrRoleNotFound,
			},
		}
	}

	r, err := intr.RoleRepository.Role(orgID, serviceName, roleName)
	if err != nil {
		return nil, err
	}

	if err := intr.RoleRepository.Delete(orgID, serviceName, roleName); err != nil {
		return nil, err
	}

	return r, nil
}

func (intr *ServiceInteractor) Metadata(orgID, serviceName, namespace string) (any, error) {
	if !intr.ServiceRepository.Exists(orgID, serviceName) {
		return nil, &ServiceNotFound{
			Err{
				Err: ErrServiceNotFound,
			},
		}
	}

	if !intr.ServiceMetaRepository.Exists(orgID, serviceName, namespace) {
		return nil, &ServiceMetadataNotFound{
			Err{
				Err: ErrServiceMetadataNotFound,
			},
		}
	}

	return intr.ServiceMetaRepository.Metadata(orgID, serviceName, namespace)
}

func (intr *ServiceInteractor) ListMetadata(orgID, serviceName string) (*domain.ServiceMetadataList, error) {
	if !intr.ServiceRepository.Exists(orgID, serviceName) {
		return nil, &ServiceNotFound{
			Err{
				Err: ErrServiceNotFound,
			},
		}
	}

	return intr.ServiceMetaRepository.List(orgID, serviceName)
}

func (intr *ServiceInteractor) SaveMetadata(orgID, serviceName, namespace string, metadata any) (*domain.Success, error) {
	if !intr.ServiceRepository.Exists(orgID, serviceName) {
		return &domain.Success{Success: false}, &ServiceNotFound{
			Err{
				Err: ErrServiceNotFound,
			},
		}
	}

	meta, err := intr.ServiceMetaRepository.List(orgID, serviceName)
	if err != nil {
		return &domain.Success{Success: false}, fmt.Errorf("get metadata list: %v", err)
	}

	if len(meta.Metadata) > 50 {
		return &domain.Success{Success: false}, &MetadataLimitExceeded{
			Err{
				Err: ErrMetadataLimitExceeded,
			},
		}
	}

	b, err := json.Marshal(metadata)
	if err != nil {
		return &domain.Success{Success: false}, fmt.Errorf("marshal: %v", err)
	}

	if len(b) > 100000 {
		return &domain.Success{Success: false}, &MetadataTooLarge{
			Err{
				Err: ErrMetadataTooLarge,
			},
		}
	}

	res, err := intr.ServiceMetaRepository.Save(orgID, serviceName, namespace, metadata)
	if err != nil {
		return res, fmt.Errorf("save metadata: %v", err)
	}

	return res, nil
}

func (intr *ServiceInteractor) DeleteMetadata(orgID, serviceName, namespace string) (*domain.Success, error) {
	if !intr.ServiceRepository.Exists(orgID, serviceName) {
		return &domain.Success{Success: false}, &ServiceNotFound{
			Err{
				Err: ErrServiceNotFound,
			},
		}
	}

	if !intr.ServiceMetaRepository.Exists(orgID, serviceName, namespace) {
		return &domain.Success{Success: false}, &ServiceMetadataNotFound{
			Err{
				Err: ErrServiceMetadataNotFound,
			},
		}
	}

	res, err := intr.ServiceMetaRepository.Delete(orgID, serviceName, namespace)
	if err != nil {
		return res, fmt.Errorf("delete metadata: %v", err)
	}

	return res, nil
}

func (intr *ServiceInteractor) RoleMetadata(orgID, serviceName, roleName, namespace string) (any, error) {
	if !intr.ServiceRepository.Exists(orgID, serviceName) {
		return nil, &ServiceNotFound{
			Err{
				Err: ErrServiceNotFound,
			},
		}
	}

	if !intr.RoleRepository.Exists(orgID, serviceName, roleName) {
		return nil, &RoleNotFound{
			Err{
				Err: ErrRoleNotFound,
			},
		}
	}

	if !intr.RoleMetaRepository.Exists(orgID, serviceName, roleName, namespace) {
		return nil, &RoleMetadataNotFound{
			Err{
				Err: ErrRoleMetadataNotFound,
			},
		}
	}

	return intr.RoleMetaRepository.Metadata(orgID, serviceName, roleName, namespace)
}

func (intr *ServiceInteractor) ListRoleMetadata(orgID, serviceName, roleName string) (*domain.RoleMetadataList, error) {
	if !intr.ServiceRepository.Exists(orgID, serviceName) {
		return nil, &ServiceNotFound{
			Err{
				Err: ErrServiceNotFound,
			},
		}
	}

	if !intr.RoleRepository.Exists(orgID, serviceName, roleName) {
		return nil, &RoleNotFound{
			Err{
				ErrRoleNotFound,
			},
		}
	}

	return intr.RoleMetaRepository.List(orgID, serviceName, roleName)
}

func (intr *ServiceInteractor) SaveRoleMetadata(orgID, serviceName, roleName, namespace string, metadata any) (*domain.Success, error) {
	if !intr.ServiceRepository.Exists(orgID, serviceName) {
		return &domain.Success{Success: false}, &ServiceNotFound{
			Err{
				Err: ErrServiceNotFound,
			},
		}
	}

	if !intr.RoleRepository.Exists(orgID, serviceName, roleName) {
		return &domain.Success{Success: false}, &RoleNotFound{
			Err{
				Err: ErrRoleNotFound,
			},
		}
	}

	meta, err := intr.RoleMetaRepository.List(orgID, serviceName, roleName)
	if err != nil {
		return &domain.Success{Success: false}, fmt.Errorf("get role metadata list: %v", err)
	}

	if len(meta.Metadata) > 50 {
		return &domain.Success{Success: false}, &MetadataLimitExceeded{
			Err{
				Err: ErrRoleMetadataLimitExceeded,
			},
		}
	}

	b, err := json.Marshal(metadata)
	if err != nil {
		return &domain.Success{Success: false}, fmt.Errorf("marshal: %v", err)
	}

	if len(b) > 100000 {
		return &domain.Success{Success: false}, &MetadataTooLarge{
			Err{
				Err: ErrMetadataTooLarge,
			},
		}
	}

	res, err := intr.RoleMetaRepository.Save(orgID, serviceName, roleName, namespace, metadata)
	if err != nil {
		return res, fmt.Errorf("save role metadata: %v", err)
	}

	return res, nil
}

func (intr *ServiceInteractor) DeleteRoleMetadata(orgID, serviceName, roleName, namespace string) (*domain.Success, error) {
	if !intr.ServiceRepository.Exists(orgID, serviceName) {
		return &domain.Success{Success: false}, &ServiceNotFound{
			Err{
				Err: ErrServiceNotFound,
			},
		}
	}

	if !intr.RoleRepository.Exists(orgID, serviceName, roleName) {
		return &domain.Success{Success: false}, &RoleNotFound{
			Err{
				Err: ErrRoleNotFound,
			},
		}
	}

	if !intr.RoleMetaRepository.Exists(orgID, serviceName, roleName, namespace) {
		return &domain.Success{Success: false}, &RoleMetadataNotFound{
			Err{
				Err: ErrRoleMetadataNotFound,
			},
		}
	}

	res, err := intr.RoleMetaRepository.Delete(orgID, serviceName, roleName, namespace)
	if err != nil {
		return res, fmt.Errorf("delete role metadata: %v", err)
	}

	return res, nil
}

func (intr *ServiceInteractor) MetricNames(orgID, serviceName string) (*domain.ServiceMetricValueNames, error) {
	if !intr.ServiceRepository.Exists(orgID, serviceName) {
		return nil, &ServiceNotFound{
			Err{
				Err: ErrServiceNotFound,
			},
		}
	}

	return intr.ServiceMetricRepository.Names(orgID, serviceName)
}

func (intr *ServiceInteractor) MetricValues(orgID, serviceName, metricName string, from, to int64) (*domain.ServiceMetricValues, error) {
	if !intr.ServiceRepository.Exists(orgID, serviceName) {
		return nil, &ServiceNotFound{
			Err{
				Err: ErrServiceNotFound,
			},
		}
	}

	if !intr.ServiceMetricRepository.Exists(orgID, serviceName, metricName) {
		return nil, &ServiceMetricNotFound{
			Err{
				Err: ErrServiceMetricNotFound,
			},
		}
	}

	return intr.ServiceMetricRepository.Values(orgID, serviceName, metricName, from, to)
}

func (intr *ServiceInteractor) SaveMetricValues(orgID, serviceName string, values []domain.ServiceMetricValue) (*domain.Success, error) {
	//TODO &ServiceMetricLimitExceeded{Err{errors.New("the number of requests per minute is exceeded")}}
	res, err := intr.ServiceMetricRepository.Save(orgID, serviceName, values)
	if err != nil {
		return res, fmt.Errorf("save metric values: %v", err)
	}

	return res, nil
}
