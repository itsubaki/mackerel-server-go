package usecase

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/itsubaki/mackerel-api/pkg/domain"
)

type HostInteractor struct {
	HostRepository       HostRepository
	HostMetaRepository   HostMetaRepository
	HostMetricRepository HostMetricRepository
	ServiceRepository    ServiceRepository
	RoleRepository       RoleRepository
}

func (s *HostInteractor) List(orgID string) (*domain.Hosts, error) {
	return s.HostRepository.List(orgID)
}

func (s *HostInteractor) Save(orgID string, host *domain.Host) (*domain.HostID, error) {
	// Update
	if len(host.ID) > 0 && !s.HostRepository.Exists(orgID, host.ID) {
		return nil, &HostNotFound{Err{errors.New(fmt.Sprintf("the host that corresponds to the <%s> can’t be located", host.ID))}}
	}
	if len(host.Status) < 1 && s.HostRepository.Exists(orgID, host.ID) {
		exists, err := s.HostRepository.Host(orgID, host.ID)
		if err != nil {
			return nil, &HostNotFound{Err{errors.New(fmt.Sprintf("the host that corresponds to the <%s> can’t be located", host.ID))}}
		}
		host.Status = exists.Status
	}

	// Create
	if len(host.ID) < 1 {
		host.ID = domain.NewHostID()
		host.CreatedAt = time.Now().Unix()
		host.RetiredAt = 0
		host.IsRetired = false
		host.Checks = []domain.Check{}
		if len(host.Status) < 1 {
			host.Status = "working"
		}
	}

	// role_fullnames -> roles
	host.Roles = make(map[string][]string)
	for i := range host.RoleFullNames {
		svc := strings.Split(host.RoleFullNames[i], ":")
		if _, ok := host.Roles[svc[0]]; !ok {
			host.Roles[svc[0]] = make([]string, 0)
		}

		host.Roles[svc[0]] = append(host.Roles[svc[0]], svc[1])
	}

	res, err := s.HostRepository.Save(orgID, host)
	if err != nil {
		return nil, fmt.Errorf("save host: %v", err)
	}

	for svc := range host.Roles {
		if s.ServiceRepository.Exists(orgID, svc) {
			continue
		}

		if err := s.ServiceRepository.Save(orgID, &domain.Service{
			OrgID: orgID,
			Name:  svc,
		}); err != nil {
			return nil, fmt.Errorf("save service<%s>: %v", svc, err)
		}
	}

	for svc, roles := range host.Roles {
		for i := range roles {
			if s.RoleRepository.Exists(orgID, svc, roles[i]) {
				continue
			}

			if err := s.RoleRepository.Save(orgID, svc, &domain.Role{
				OrgID:       orgID,
				ServiceName: svc,
				Name:        roles[i],
			}); err != nil {
				return nil, fmt.Errorf("save role<%s:%s>: %v", svc, roles[i], err)
			}
		}
	}

	return res, nil
}

func (s *HostInteractor) Host(orgID, hostID string) (*domain.HostInfo, error) {
	host, err := s.HostRepository.Host(orgID, hostID)
	if err != nil {
		return nil, fmt.Errorf("get host<%s>: %v", hostID, err)
	}

	return &domain.HostInfo{Host: *host}, nil
}

func (s *HostInteractor) Status(orgID, hostID, status string) (*domain.Success, error) {
	if !s.HostRepository.Exists(orgID, hostID) {
		return nil, &HostNotFound{Err{errors.New(fmt.Sprintf("the host that corresponds to the <%s> can’t be located", hostID))}}
	}

	res, err := s.HostRepository.Status(orgID, hostID, status)
	if err != nil {
		return res, fmt.Errorf("host status: %v", err)
	}

	return res, nil
}

func (s *HostInteractor) SaveRoleFullNames(orgID, hostID string, names *domain.RoleFullNames) (*domain.Success, error) {
	if !s.HostRepository.Exists(orgID, hostID) {
		return nil, &HostNotFound{Err{errors.New(fmt.Sprintf("the host that corresponds to the <%s> can’t be located", hostID))}}
	}

	res, err := s.HostRepository.SaveRoleFullNames(orgID, hostID, names)
	if err != nil {
		return res, fmt.Errorf("save role fullnames: %v", err)
	}

	for svc, roles := range names.Roles() {
		if err := s.ServiceRepository.Save(orgID, &domain.Service{
			OrgID: orgID,
			Name:  svc,
		}); err != nil {
			return nil, fmt.Errorf("save service<%s>: %v", svc, err)
		}

		for i := range roles {
			if err := s.RoleRepository.Save(orgID, svc, &domain.Role{
				OrgID:       orgID,
				ServiceName: svc,
				Name:        roles[i],
			}); err != nil {
				return nil, fmt.Errorf("save role<%s:%s>: %v", svc, roles[i], err)
			}
		}
	}

	return res, nil
}

func (s *HostInteractor) Retire(orgID, hostID string, retire *domain.HostRetire) (*domain.Success, error) {
	if !s.HostRepository.Exists(orgID, hostID) {
		return nil, &HostNotFound{Err{errors.New(fmt.Sprintf("the host that corresponds to the <%s> can’t be located", hostID))}}
	}

	res, err := s.HostRepository.Retire(orgID, hostID, retire)
	if err != nil {
		return res, fmt.Errorf("retire host<%s>: %v", hostID, err)
	}

	return res, nil
}

func (s *HostInteractor) MetricNames(orgID, hostID string) (*domain.MetricNames, error) {
	if !s.HostRepository.Exists(orgID, hostID) {
		return nil, &HostNotFound{Err{errors.New(fmt.Sprintf("the host that corresponds to the <%s> can’t be located", hostID))}}
	}

	return s.HostMetricRepository.Names(orgID, hostID)
}

func (s *HostInteractor) MetricValues(orgID, hostID, name string, from, to int64) (*domain.MetricValues, error) {
	if !s.HostRepository.Exists(orgID, hostID) {
		return nil, &HostNotFound{Err{errors.New(fmt.Sprintf("the host<%s> doesn't exist", hostID))}}
	}

	if !s.HostMetricRepository.Exists(orgID, hostID, name) {
		return nil, &HostMetricNotFound{Err{errors.New(fmt.Sprintf("the metric<%s:%s> doesn't exist", hostID, name))}}
	}

	return s.HostMetricRepository.Values(orgID, hostID, name, from, to)
}

func (s *HostInteractor) MetricValuesLatest(orgID string, hostId, name []string) (*domain.TSDBLatest, error) {
	return s.HostMetricRepository.ValuesLatest(orgID, hostId, name)
}

func (s *HostInteractor) SaveMetricValues(orgID string, values []domain.MetricValue) (*domain.Success, error) {
	res, err := s.HostMetricRepository.Save(orgID, values)
	if err != nil {
		return res, fmt.Errorf("save host metric values: %v", err)
	}

	return res, nil
}

func (s *HostInteractor) MetadataList(orgID, hostID string) (*domain.HostMetadataList, error) {
	if !s.HostRepository.Exists(orgID, hostID) {
		return nil, &HostNotFound{Err{errors.New(fmt.Sprintf("the host<%s> does not exist", hostID))}}
	}

	h, err := s.HostRepository.Host(orgID, hostID)
	if err != nil {
		return nil, err
	}

	if h.IsRetired {
		return nil, &HostIsRetired{Err{errors.New(fmt.Sprintf("more than a week has passed since the host<%s> retired", hostID))}}
	}

	return s.HostMetaRepository.List(orgID, hostID)
}

func (s *HostInteractor) Metadata(orgID, hostID, namespace string) (interface{}, error) {
	if !s.HostRepository.Exists(orgID, hostID) {
		return nil, &HostNotFound{Err{errors.New(fmt.Sprintf("the host<%s> does not exist", hostID))}}
	}

	if !s.HostMetaRepository.Exists(orgID, hostID, namespace) {
		return nil, &HostMetadataNotFound{Err{errors.New(fmt.Sprintf("the specified metadata does not exist for the host<%s:%s>", hostID, namespace))}}
	}

	h, err := s.HostRepository.Host(orgID, hostID)
	if err != nil {
		return nil, err
	}

	if h.IsRetired {
		return nil, &HostIsRetired{Err{errors.New(fmt.Sprintf("more than a week has passed since the host<%s> retired", hostID))}}
	}

	return s.HostMetaRepository.Metadata(orgID, hostID, namespace)
}

func (s *HostInteractor) SaveMetadata(orgID, hostID, namespace string, metadata interface{}) (*domain.Success, error) {
	if !s.HostRepository.Exists(orgID, hostID) {
		return nil, &HostNotFound{Err{errors.New(fmt.Sprintf("the host<%s> does not exist", hostID))}}
	}

	h, err := s.HostRepository.Host(orgID, hostID)
	if err != nil {
		return &domain.Success{Success: false}, fmt.Errorf("get host<%s>: %v", hostID, err)
	}

	if h.IsRetired {
		return &domain.Success{Success: false}, &HostIsRetired{Err{errors.New(fmt.Sprintf("the host<%s> has already been retired", hostID))}}
	}

	meta, err := s.HostMetaRepository.List(orgID, hostID)
	if err != nil {
		return &domain.Success{Success: false}, fmt.Errorf("get metadata list: %v", err)
	}

	if len(meta.Metadata) > 50 {
		return &domain.Success{Success: false}, &MetadataLimitExceeded{Err{errors.New(fmt.Sprintf("metadata<%s:%s> limit (50 per 1 host) has been exceeded and you try to register", hostID, namespace))}}
	}

	b, err := json.Marshal(metadata)
	if err != nil {
		return &domain.Success{Success: false}, fmt.Errorf("marshal: %v", err)
	}

	if len(b) > 100000 {
		return &domain.Success{Success: false}, &MetadataTooLarge{Err{errors.New(fmt.Sprintf("the metadata<%s:%s> exceeds 100KB", hostID, namespace))}}
	}

	res, err := s.HostMetaRepository.Save(orgID, hostID, namespace, metadata)
	if err != nil {
		return res, fmt.Errorf("save metadata: %v", err)
	}

	return res, nil
}

func (s *HostInteractor) DeleteMetadata(orgID, hostID, namespace string) (*domain.Success, error) {
	if !s.HostRepository.Exists(orgID, hostID) {
		return &domain.Success{Success: false}, &HostNotFound{Err{errors.New(fmt.Sprintf("the host<%s> does not exist", hostID))}}
	}

	if !s.HostMetaRepository.Exists(orgID, hostID, namespace) {
		return &domain.Success{Success: false}, &HostMetadataNotFound{Err{errors.New(fmt.Sprintf("the specified metadata does not exist for the host<%s:%s>", hostID, namespace))}}
	}

	h, err := s.HostRepository.Host(orgID, hostID)
	if err != nil {
		return &domain.Success{Success: false}, fmt.Errorf("get host<%s>: %v", hostID, err)
	}

	if h.IsRetired {
		return &domain.Success{Success: false}, &HostIsRetired{Err{errors.New(fmt.Sprintf("the host<%s> has already been retired", hostID))}}
	}

	res, err := s.HostMetaRepository.Delete(orgID, hostID, namespace)
	if err != nil {
		return res, fmt.Errorf("delete metadata: %v", err)
	}

	return res, nil
}
