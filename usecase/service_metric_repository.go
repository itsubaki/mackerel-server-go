package usecase

import "github.com/itsubaki/mackerel-server-go/domain"

type ServiceMetricRepository interface {
	Exists(orgID, serviceName, metricName string) bool
	Names(orgID, serviceName string) (*domain.ServiceMetricValueNames, error)
	Values(orgID, serviceName, metricName string, from, to int64) (*domain.ServiceMetricValues, error)
	Save(orgID, serviceName string, values []domain.ServiceMetricValue) (*domain.Success, error)
}
