package usecase

import "github.com/itsubaki/mackerel-api/pkg/domain"

type HostMetricInteractor struct {
	HostMetricRepository HostMetricRepository
}

func (s *HostMetricInteractor) Save(v domain.HostMetricValues) error {
	return s.HostMetricRepository.Save(v)
}

func (s *HostMetricInteractor) Find(hostID, metricName string, from, to int64) (domain.HostMetricValues, error) {
	return s.HostMetricRepository.FindBy(hostID, metricName, from, to)
}

func (s *HostMetricInteractor) Latest(hostID, metricName []string) (domain.HostMetricValues, error) {
	return s.HostMetricRepository.Latest(hostID, metricName)
}

func (s *HostMetricInteractor) MetricNames(hostID string) ([]string, error) {
	v, err := s.HostMetricRepository.FindByID(hostID)
	if err != nil {
		return nil, err
	}

	return v.MetricNames(), nil
}
