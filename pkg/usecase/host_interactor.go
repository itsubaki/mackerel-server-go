package usecase

import (
	"encoding/json"
	"errors"
	"strings"
	"time"

	"github.com/itsubaki/mackerel-api/pkg/domain"
)

type HostInteractor struct {
	HostRepository HostRepository
}

func (s *HostInteractor) List(orgID string) (*domain.Hosts, error) {
	return s.HostRepository.List(orgID)
}

func (s *HostInteractor) Save(orgID string, host *domain.Host) (*domain.HostID, error) {
	// Update
	if len(host.ID) > 0 && !s.HostRepository.Exists(orgID, host.ID) {
		return nil, &HostNotFound{Err{errors.New("the host that corresponds to the <hostId> can’t be located")}}
	}
	if len(host.Status) < 1 && s.HostRepository.Exists(orgID, host.ID) {
		exists, err := s.HostRepository.Host(orgID, host.ID)
		if err != nil {
			return nil, &HostNotFound{Err{errors.New("the host that corresponds to the <hostId> can’t be located")}}
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

	return s.HostRepository.Save(orgID, host)
}

func (s *HostInteractor) Host(orgID, hostID string) (*domain.HostInfo, error) {
	host, err := s.HostRepository.Host(orgID, hostID)
	if err != nil {
		return nil, err
	}

	return &domain.HostInfo{Host: *host}, nil
}

func (s *HostInteractor) Status(orgID, hostID, status string) (*domain.Success, error) {
	if !s.HostRepository.Exists(orgID, hostID) {
		return nil, &HostNotFound{Err{errors.New("the host that corresponds to the <hostId> can’t be located")}}
	}

	return s.HostRepository.Status(orgID, hostID, status)
}

func (s *HostInteractor) SaveRoleFullNames(orgID, hostID string, names *domain.RoleFullNames) (*domain.Success, error) {
	if !s.HostRepository.Exists(orgID, hostID) {
		return nil, &HostNotFound{Err{errors.New("the host that corresponds to the <hostId> can’t be located")}}
	}

	return s.HostRepository.SaveRoleFullNames(orgID, hostID, names)
}

func (s *HostInteractor) Retire(orgID, hostID string, retire *domain.HostRetire) (*domain.Success, error) {
	if !s.HostRepository.Exists(orgID, hostID) {
		return nil, &HostNotFound{Err{errors.New("the host that corresponds to the <hostId> can’t be located")}}
	}

	return s.HostRepository.Retire(orgID, hostID, retire)
}

func (s *HostInteractor) MetricNames(orgID, hostID string) (*domain.MetricNames, error) {
	if !s.HostRepository.Exists(orgID, hostID) {
		return nil, &HostNotFound{Err{errors.New("the host that corresponds to the <hostId> can’t be located")}}
	}

	return s.HostRepository.MetricNames(orgID, hostID)
}

func (s *HostInteractor) MetricValues(orgID, hostID, name string, from, to int64) (*domain.MetricValues, error) {
	if !s.HostRepository.Exists(orgID, hostID) {
		return nil, &HostNotFound{Err{errors.New("the host doesn't exist")}}
	}

	if !s.HostRepository.ExistsMetric(orgID, hostID, name) {
		return nil, &HostMetricNotFound{Err{errors.New("the metric doesn't exist")}}
	}

	return s.HostRepository.MetricValues(orgID, hostID, name, from, to)
}

func (s *HostInteractor) MetricValuesLatest(orgID string, hostId, name []string) (*domain.TSDBLatest, error) {
	return s.HostRepository.MetricValuesLatest(orgID, hostId, name)
}

func (s *HostInteractor) SaveMetricValues(orgID string, values []domain.MetricValue) (*domain.Success, error) {
	return s.HostRepository.SaveMetricValues(orgID, values)
}

func (s *HostInteractor) MetadataList(orgID, hostID string) (*domain.HostMetadataList, error) {
	if !s.HostRepository.Exists(orgID, hostID) {
		return nil, &HostNotFound{Err{errors.New("the host does not exist")}}
	}

	h, err := s.HostRepository.Host(orgID, hostID)
	if err != nil {
		return nil, err
	}

	if h.IsRetired {
		return nil, &HostIsRetired{Err{errors.New("more than a week has passed since the host retired")}}
	}

	return s.HostRepository.MetadataList(orgID, hostID)
}

func (s *HostInteractor) Metadata(orgID, hostID, namespace string) (interface{}, error) {
	if !s.HostRepository.Exists(orgID, hostID) {
		return nil, &HostNotFound{Err{errors.New("the host does not exist")}}
	}

	if !s.HostRepository.ExistsMetadata(orgID, hostID, namespace) {
		return nil, &HostMetadataNotFound{Err{errors.New("the specified metadata does not exist for the host")}}
	}

	h, err := s.HostRepository.Host(orgID, hostID)
	if err != nil {
		return nil, err
	}

	if h.IsRetired {
		return nil, &HostIsRetired{Err{errors.New("more than a week has passed since the host retired")}}
	}

	return s.HostRepository.Metadata(orgID, hostID, namespace)
}

func (s *HostInteractor) SaveMetadata(orgID, hostID, namespace string, metadata interface{}) (*domain.Success, error) {
	if !s.HostRepository.Exists(orgID, hostID) {
		return nil, &HostNotFound{Err{errors.New("the host does not exist")}}
	}

	h, err := s.HostRepository.Host(orgID, hostID)
	if err != nil {
		return nil, err
	}

	if h.IsRetired {
		return nil, &HostIsRetired{Err{errors.New("the host has already been retired")}}
	}

	meta, err := s.HostRepository.MetadataList(orgID, hostID)
	if err != nil {
		return nil, err
	}

	if len(meta.Metadata) > 50 {
		return nil, &MetadataLimitExceeded{Err{errors.New("metadata limit (50 per 1 host) has been exceeded and you try to register")}}
	}

	b, err := json.Marshal(metadata)
	if err != nil {
		return nil, err
	}

	if len(b) > 100000 {
		return nil, &MetadataTooLarge{Err{errors.New("the metadata exceeds 100KB")}}
	}

	return s.HostRepository.SaveMetadata(orgID, hostID, namespace, metadata)
}

func (s *HostInteractor) DeleteMetadata(orgID, hostID, namespace string) (*domain.Success, error) {
	if !s.HostRepository.Exists(orgID, hostID) {
		return nil, &HostNotFound{Err{errors.New("the host does not exist")}}
	}

	if !s.HostRepository.ExistsMetadata(orgID, hostID, namespace) {
		return nil, &HostMetadataNotFound{Err{errors.New("the specified metadata does not exist for the host")}}
	}

	h, err := s.HostRepository.Host(orgID, hostID)
	if err != nil {
		return nil, err
	}

	if h.IsRetired {
		return nil, &HostIsRetired{Err{errors.New("the host has already been retired")}}
	}

	return s.HostRepository.DeleteMetadata(orgID, hostID, namespace)
}
