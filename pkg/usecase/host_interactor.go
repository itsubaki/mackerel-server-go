package usecase

import (
	"encoding/json"
	"errors"

	"github.com/itsubaki/mackerel-api/pkg/domain"
)

type HostInteractor struct {
	HostRepository HostRepository
}

func (s *HostInteractor) List(org string) (*domain.Hosts, error) {
	return s.HostRepository.List(org)
}

func (s *HostInteractor) Save(org string, host *domain.Host) (*domain.HostID, error) {
	// Update
	if len(host.ID) > 0 && !s.HostRepository.Exists(org, host.ID) {
		return nil, &HostNotFound{Err{errors.New("the host that corresponds to the <hostId> can’t be located")}}
	}

	host.Init()
	return s.HostRepository.Save(org, host)
}

func (s *HostInteractor) Host(org, hostID string) (*domain.HostInfo, error) {
	host, err := s.HostRepository.Host(org, hostID)
	if err != nil {
		return nil, err
	}

	return &domain.HostInfo{Host: *host}, nil
}

func (s *HostInteractor) Status(org, hostID, status string) (*domain.Success, error) {
	if !s.HostRepository.Exists(org, hostID) {
		return nil, &HostNotFound{Err{errors.New("the host that corresponds to the <hostId> can’t be located")}}
	}

	return s.HostRepository.Status(org, hostID, status)
}

func (s *HostInteractor) SaveRoleFullNames(org, hostID string, names *domain.RoleFullNames) (*domain.Success, error) {
	if !s.HostRepository.Exists(org, hostID) {
		return nil, &HostNotFound{Err{errors.New("the host that corresponds to the <hostId> can’t be located")}}
	}

	return s.HostRepository.SaveRoleFullNames(org, hostID, names)
}

func (s *HostInteractor) Retire(org, hostID string, retire *domain.HostRetire) (*domain.Success, error) {
	if !s.HostRepository.Exists(org, hostID) {
		return nil, &HostNotFound{Err{errors.New("the host that corresponds to the <hostId> can’t be located")}}
	}

	return s.HostRepository.Retire(org, hostID, retire)
}

func (s *HostInteractor) MetricNames(org, hostID string) (*domain.MetricNames, error) {
	if !s.HostRepository.Exists(org, hostID) {
		return nil, &HostNotFound{Err{errors.New("the host that corresponds to the <hostId> can’t be located")}}
	}

	return s.HostRepository.MetricNames(org, hostID)
}

func (s *HostInteractor) MetricValues(org, hostID, name string, from, to int64) (*domain.MetricValues, error) {
	if !s.HostRepository.Exists(org, hostID) {
		return nil, &HostNotFound{Err{errors.New("the host doesn't exist")}}
	}

	if !s.HostRepository.ExistsMetric(org, hostID, name) {
		return nil, &HostMetricNotFound{Err{errors.New("the metric doesn't exist")}}
	}

	return s.HostRepository.MetricValues(org, hostID, name, from, to)
}

func (s *HostInteractor) MetricValuesLatest(org string, hostId, name []string) (*domain.TSDBLatest, error) {
	return s.HostRepository.MetricValuesLatest(org, hostId, name)
}

func (s *HostInteractor) SaveMetricValues(org string, values []domain.MetricValue) (*domain.Success, error) {
	return s.HostRepository.SaveMetricValues(org, values)
}

func (s *HostInteractor) MetadataList(org, hostID string) (*domain.HostMetadataList, error) {
	if !s.HostRepository.Exists(org, hostID) {
		return nil, &HostNotFound{Err{errors.New("the host does not exist")}}
	}

	h, err := s.HostRepository.Host(org, hostID)
	if err != nil {
		return nil, err
	}

	if h.IsRetired {
		return nil, &HostIsRetired{Err{errors.New("more than a week has passed since the host retired")}}
	}

	return s.HostRepository.MetadataList(org, hostID)
}

func (s *HostInteractor) Metadata(org, hostID, namespace string) (interface{}, error) {
	if !s.HostRepository.Exists(org, hostID) {
		return nil, &HostNotFound{Err{errors.New("the host does not exist")}}
	}

	if !s.HostRepository.ExistsMetadata(org, hostID, namespace) {
		return nil, &HostMetadataNotFound{Err{errors.New("the specified metadata does not exist for the host")}}
	}

	h, err := s.HostRepository.Host(org, hostID)
	if err != nil {
		return nil, err
	}

	if h.IsRetired {
		return nil, &HostIsRetired{Err{errors.New("more than a week has passed since the host retired")}}
	}

	return s.HostRepository.Metadata(org, hostID, namespace)
}

func (s *HostInteractor) SaveMetadata(org, hostID, namespace string, metadata interface{}) (*domain.Success, error) {
	if !s.HostRepository.Exists(org, hostID) {
		return nil, &HostNotFound{Err{errors.New("the host does not exist")}}
	}

	h, err := s.HostRepository.Host(org, hostID)
	if err != nil {
		return nil, err
	}

	if h.IsRetired {
		return nil, &HostIsRetired{Err{errors.New("the host has already been retired")}}
	}

	meta, err := s.HostRepository.MetadataList(org, hostID)
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

	return s.HostRepository.SaveMetadata(org, hostID, namespace, metadata)
}

func (s *HostInteractor) DeleteMetadata(org, hostID, namespace string) (*domain.Success, error) {
	if !s.HostRepository.Exists(org, hostID) {
		return nil, &HostNotFound{Err{errors.New("the host does not exist")}}
	}

	if !s.HostRepository.ExistsMetadata(org, hostID, namespace) {
		return nil, &HostMetadataNotFound{Err{errors.New("the specified metadata does not exist for the host")}}
	}

	h, err := s.HostRepository.Host(org, hostID)
	if err != nil {
		return nil, err
	}

	if h.IsRetired {
		return nil, &HostIsRetired{Err{errors.New("the host has already been retired")}}
	}

	return s.HostRepository.DeleteMetadata(org, hostID, namespace)
}
