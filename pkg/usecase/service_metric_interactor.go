package usecase

import "github.com/itsubaki/mackerel-api/pkg/domain"

type ServiceMetricInteractor struct {
	ServiceMetricRepository ServiceMetricRepository
}

func (s *ServiceMetricInteractor) Save(values domain.ServiceMetricValues) error {
	return s.ServiceMetricRepository.Save(values)
}

func (s *ServiceMetricInteractor) Find(serviceName, metricName string, from, to int64) (domain.ServiceMetricValues, error) {
	return s.ServiceMetricRepository.FindBy(serviceName, metricName, from, to)
}

func (s *ServiceMetricInteractor) MetricNames(serviceName string) ([]string, error) {
	v, err := s.ServiceMetricRepository.FindAll()
	if err != nil {
		return nil, err
	}

	return v.MetricNames(), nil
}
