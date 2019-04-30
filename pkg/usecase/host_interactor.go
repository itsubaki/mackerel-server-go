package usecase

import "github.com/itsubaki/mackerel-api/pkg/domain"

type HostInteractor struct {
	HostRepository HostRepository
}

func (s *HostInteractor) Find(serviceName, hostName, status, customIdentifier string, roleName []string) (domain.Hosts, error) {
	return domain.Hosts{}, nil
}

func (s *HostInteractor) FindByID(hostID string) (*domain.Host, error) {
	host, err := s.HostRepository.FindByID(hostID)
	return &host, err
}

func (s *HostInteractor) Save(host *domain.Host) (string, error) {
	if err := s.HostRepository.Save(*host); err != nil {
		return host.ID, err
	}

	return host.ID, nil
}

func (s *HostInteractor) Delete(hostID string) error {
	return s.HostRepository.DeleteByID(hostID)
}
