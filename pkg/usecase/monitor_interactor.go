package usecase

import (
	"github.com/itsubaki/mackerel-api/pkg/domain"
)

type MonitorInteractor struct {
	MonitorRepository MonitorRepository
}

func (s *MonitorInteractor) List(orgID string) (*domain.Monitors, error) {
	return s.MonitorRepository.List(orgID)
}

func (s *MonitorInteractor) Save(orgID string, monitor *domain.Monitoring) (interface{}, error) {
	monitor.ID = domain.NewIDWith(monitor.Name)

	if monitor.Type == "external" && len(monitor.Method) < 1 {
		monitor.Method = "GET"
	}
	return s.MonitorRepository.Save(orgID, monitor)
}

func (s *MonitorInteractor) Update(orgID string, monitor *domain.Monitoring) (interface{}, error) {
	return s.MonitorRepository.Update(orgID, monitor)
}

func (s *MonitorInteractor) Monitor(orgID, monitorID string) (interface{}, error) {
	return s.MonitorRepository.Monitor(orgID, monitorID)
}

func (s *MonitorInteractor) Delete(orgID, monitorID string) (interface{}, error) {
	return s.MonitorRepository.Delete(orgID, monitorID)
}
