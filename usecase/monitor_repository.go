package usecase

import "github.com/itsubaki/mackerel-server-go/domain"

type MonitorRepository interface {
	List(orgID string) (*domain.Monitors, error)
	ListHostMetric(orgID string) ([]domain.HostMetricMonitoring, error)
	Save(orgID string, monitor *domain.Monitoring) (interface{}, error)
	Update(orgID string, monitor *domain.Monitoring) (interface{}, error)
	Monitor(orgID, monitorID string) (interface{}, error)
	Delete(orgID, monitorID string) (interface{}, error)
}
