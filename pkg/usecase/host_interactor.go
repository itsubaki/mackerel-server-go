package usecase

import "github.com/itsubaki/mackerel-api/pkg/domain"

type HostInteractor struct {
	HostRepository HostRepository
}

func (s *HostInteractor) SaveCustomGraphDefs(defs domain.CustomGraphDefs) error {
	return s.HostRepository.SaveCustomGraphDefs(defs)
}

func (s *HostInteractor) MetadataList(hostID string) (domain.HostMetadataList, error) {
	return domain.HostMetadataList{}, nil
}

func (s *HostInteractor) Metadata(hostID, namespace string) (interface{}, error) {
	return nil, nil
}

func (s *HostInteractor) SaveMetadata(hostID, namespace string, metadata interface{}) error {
	return nil
}

func (s *HostInteractor) DeleteMetadata(hostID, namespace string) error {
	return nil
}

func (s *HostInteractor) MetricNames(hostID string) ([]string, error) {
	return s.HostRepository.MetricNames()
}

func (s *HostInteractor) SaveMetricValues(v domain.HostMetricValues) error {
	return s.HostRepository.SaveMetricValues(v)
}

func (s *HostInteractor) MetricValues(hostID, metricName string, from, to int64) (domain.HostMetricValues, error) {
	return s.HostRepository.MetricValues(hostID, metricName, from, to)
}

func (s *HostInteractor) MetricValuesLatest(hostID, metricName []string) (domain.HostMetricValues, error) {
	return s.HostRepository.MetricValuesLatest(hostID, metricName)
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
