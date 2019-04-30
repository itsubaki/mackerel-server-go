package usecase

import (
	"github.com/itsubaki/mackerel-api/pkg/domain"
	"github.com/itsubaki/mackerel-api/pkg/interfaces/database"
)

func NewHostMetricInteractor() *HostMetricInteractor {
	return &HostMetricInteractor{
		HostMetricRepository: database.NewHostMetricRepository(),
	}
}

type HostMetricInteractor struct {
	HostMetricRepository *database.HostMetricRepository
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
