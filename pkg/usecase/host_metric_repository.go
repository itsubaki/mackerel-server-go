package usecase

import "github.com/itsubaki/mackerel-api/pkg/domain"

type HostMetricRepository interface {
	Exists(orgID, hostID, name string) bool
	Names(orgID, hostID string) (*domain.MetricNames, error)
	Values(orgID, hostID, name string, from, to int64) (*domain.MetricValues, error)
	ValuesLimit(orgID, hostID, name string, limit int) (*domain.MetricValues, error)
	ValuesLatest(orgID string, hostID, name []string) (*domain.TSDBLatest, error)
	ValuesAverage(orgID, hostID, name string, duration int) (*domain.MetricValueAverage, error)
	Save(orgID string, values []domain.MetricValue) (*domain.Success, error)
}
