package usecase

import (
	"github.com/itsubaki/mackerel-api/pkg/domain"
)

type MonitorInteractor struct {
	MonitorRepository MonitorRepository
}

func (s *MonitorInteractor) List(org string) (*domain.Monitors, error) {
	return s.MonitorRepository.List(org)
}

func (s *MonitorInteractor) Save(org string, monitor *domain.Monitoring) (interface{}, error) {
	monitor.ID = domain.NewMonitorID(monitor.Name)
	return s.MonitorRepository.Save(org, monitor)
}

func (s *MonitorInteractor) Update(org string, monitor *domain.Monitoring) (interface{}, error) {
	return s.MonitorRepository.Update(org, monitor)
}

func (s *MonitorInteractor) Monitor(org, monitorID string) (interface{}, error) {
	return s.MonitorRepository.Monitor(org, monitorID)
}

func (s *MonitorInteractor) Delete(org, monitorID string) (interface{}, error) {
	return s.MonitorRepository.Delete(org, monitorID)
}
