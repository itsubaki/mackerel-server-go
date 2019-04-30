package usecase

import "github.com/itsubaki/mackerel-api/pkg/domain"

type ServiceMetricRepository interface {
	FindAll() (domain.ServiceMetricValues, error)
	FindBy(serviceName, metricName string, from, to int64) (domain.ServiceMetricValues, error)
	Save(v domain.ServiceMetricValues) error
}
