package usecase

import (
	"errors"

	"github.com/itsubaki/mackerel-api/pkg/domain"
)

type HostInteractor struct {
	HostRepository HostRepository
}

func (s *HostInteractor) List() (*domain.Hosts, error) {
	return s.HostRepository.List()
}

func (s *HostInteractor) Save(host *domain.Host) (*domain.HostID, error) {
	host.ID = "genHostID"
	return s.HostRepository.Save(host)
}

func (s *HostInteractor) Update(host *domain.Host) (*domain.HostID, error) {
	if !s.HostRepository.Exists(host.ID) {
		return nil, &HostNotFound{Err{errors.New("the host that corresponds to the <hostId> can’t be located")}}
	}

	return s.HostRepository.Save(host)
}

func (s *HostInteractor) Host(hostID string) (*domain.Host, error) {
	return s.HostRepository.Host(hostID)
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

func (s *HostInteractor) MetricValues(hostID, name string, from, to int) (*domain.MetricValues, error) {
	return s.HostRepository.MetricValues(hostID, name, from, to)
}

func (s *HostInteractor) MetadataList(hostID string) (*domain.HostMetadata, error) {
	return s.HostRepository.MetadataList(hostID)
}

func (s *HostInteractor) Metadata(hostID, namespace string) (interface{}, error) {
	return s.HostRepository.Metadata(hostID, namespace)
}

func (s *HostInteractor) SaveMetadata(hostID, namespace string, metadata interface{}) (*domain.Success, error) {
	return s.HostRepository.SaveMetadata(hostID, namespace, metadata)
}

func (s *HostInteractor) DeleteMetadata(hostID, namespace string) (*domain.Success, error) {
	return s.HostRepository.DeleteMetadata(hostID, namespace)
}
