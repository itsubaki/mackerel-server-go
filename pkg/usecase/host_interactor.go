package usecase

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/itsubaki/mackerel-api/pkg/domain"
)

type HostInteractor struct {
	HostRepository HostRepository
}

func (s *HostInteractor) List() (*domain.Hosts, error) {
	return s.HostRepository.List()
}

func (s *HostInteractor) Save(host *domain.Host) (*domain.HostID, error) {
	sha := sha256.Sum256([]byte(uuid.Must(uuid.NewRandom()).String()))
	hash := hex.EncodeToString(sha[:])

	host.ID = hash[:11]
	host.CreatedAt = time.Now().Unix()
	host.RetiredAt = 0
	host.IsRetired = false
	host.Checks = []domain.Check{}
	if len(host.Status) == 0 {
		host.Status = "working"
	}

	return s.HostRepository.Save(host)
}

func (s *HostInteractor) Update(host *domain.Host) (*domain.HostID, error) {
	if !s.HostRepository.Exists(host.ID) {
		return nil, &HostNotFound{Err{errors.New("the host that corresponds to the <hostId> can’t be located")}}
	}

	return s.HostRepository.Save(host)
}

func (s *HostInteractor) Host(hostID string) (*domain.HostInfo, error) {
	host, err := s.HostRepository.Host(hostID)
	if err != nil {
		return nil, err
	}

	return &domain.HostInfo{Host: *host}, nil
}

func (s *HostInteractor) Status(hostID, status string) (*domain.Success, error) {
	if !s.HostRepository.Exists(hostID) {
		return nil, &HostNotFound{Err{errors.New("the host that corresponds to the <hostId> can’t be located")}}
	}

	return s.HostRepository.Status(hostID, status)
}

func (s *HostInteractor) SaveRoleFullNames(hostID string, names *domain.RoleFullNames) (*domain.Success, error) {
	if !s.HostRepository.Exists(hostID) {
		return nil, &HostNotFound{Err{errors.New("the host that corresponds to the <hostId> can’t be located")}}
	}

	return s.HostRepository.SaveRoleFullNames(hostID, names)
}

func (s *HostInteractor) Retire(hostID string, retire *domain.HostRetire) (*domain.Success, error) {
	if !s.HostRepository.Exists(hostID) {
		return nil, &HostNotFound{Err{errors.New("the host that corresponds to the <hostId> can’t be located")}}
	}

	return s.HostRepository.Retire(hostID, retire)
}

func (s *HostInteractor) MetricNames(hostID string) (*domain.MetricNames, error) {
	if !s.HostRepository.Exists(hostID) {
		return nil, &HostNotFound{Err{errors.New("the host that corresponds to the <hostId> can’t be located")}}
	}

	return s.HostRepository.MetricNames(hostID)
}

func (s *HostInteractor) MetricValues(hostID, name string, from, to int64) (*domain.MetricValues, error) {
	if !s.HostRepository.Exists(hostID) {
		return nil, &HostNotFound{Err{errors.New("the host doesn't exist")}}
	}

	if !s.HostRepository.ExistsMetric(hostID, name) {
		return nil, &HostMetricNotFound{Err{errors.New("the metric doesn't exist")}}
	}

	return s.HostRepository.MetricValues(hostID, name, from, to)
}

func (s *HostInteractor) MetricValuesLatest(hostId, name []string) (*domain.TSDBLatest, error) {
	return s.HostRepository.MetricValuesLatest(hostId, name)
}

func (s *HostInteractor) SaveMetricValues(values []domain.MetricValue) (*domain.Success, error) {
	return s.HostRepository.SaveMetricValues(values)
}

func (s *HostInteractor) MetadataList(hostID string) (*domain.HostMetadataList, error) {
	if !s.HostRepository.Exists(hostID) {
		return nil, &HostNotFound{Err{errors.New("the host does not exist")}}
	}

	h, err := s.HostRepository.Host(hostID)
	if err != nil {
		return nil, err
	}

	if h.IsRetired {
		return nil, &HostIsRetired{Err{errors.New("more than a week has passed since the host retired")}}
	}

	return s.HostRepository.MetadataList(hostID)
}

func (s *HostInteractor) Metadata(hostID, namespace string) (interface{}, error) {
	if !s.HostRepository.Exists(hostID) {
		return nil, &HostNotFound{Err{errors.New("the host does not exist")}}
	}

	if !s.HostRepository.ExistsMetadata(hostID, namespace) {
		return nil, &HostMetadataNotFound{Err{errors.New("the specified metadata does not exist for the host")}}
	}

	h, err := s.HostRepository.Host(hostID)
	if err != nil {
		return nil, err
	}

	if h.IsRetired {
		return nil, &HostIsRetired{Err{errors.New("more than a week has passed since the host retired")}}
	}

	return s.HostRepository.Metadata(hostID, namespace)
}

func (s *HostInteractor) SaveMetadata(hostID, namespace string, metadata interface{}) (*domain.Success, error) {
	if !s.HostRepository.Exists(hostID) {
		return nil, &HostNotFound{Err{errors.New("the host does not exist")}}
	}

	h, err := s.HostRepository.Host(hostID)
	if err != nil {
		return nil, err
	}

	if h.IsRetired {
		return nil, &HostIsRetired{Err{errors.New("the host has already been retired")}}
	}

	meta, err := s.HostRepository.MetadataList(hostID)
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

	return s.HostRepository.SaveMetadata(hostID, namespace, metadata)
}

func (s *HostInteractor) DeleteMetadata(hostID, namespace string) (*domain.Success, error) {
	if !s.HostRepository.Exists(hostID) {
		return nil, &HostNotFound{Err{errors.New("the host does not exist")}}
	}

	if !s.HostRepository.ExistsMetadata(hostID, namespace) {
		return nil, &HostMetadataNotFound{Err{errors.New("the specified metadata does not exist for the host")}}
	}

	h, err := s.HostRepository.Host(hostID)
	if err != nil {
		return nil, err
	}

	if h.IsRetired {
		return nil, &HostIsRetired{Err{errors.New("the host has already been retired")}}
	}

	return s.HostRepository.DeleteMetadata(hostID, namespace)
}
