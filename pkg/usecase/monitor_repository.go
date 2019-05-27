package usecase

import "github.com/itsubaki/mackerel-api/pkg/domain"

type MonitorRepository interface {
	List(org string) (*domain.Monitors, error)
	Save(org string, monitor *domain.Monitoring) (interface{}, error)
	Update(org string, monitor *domain.Monitoring) (interface{}, error)
	Monitor(org, monitorID string) (interface{}, error)
	Delete(org, monitorID string) (interface{}, error)
}
