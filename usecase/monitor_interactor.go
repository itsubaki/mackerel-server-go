package usecase

import "github.com/itsubaki/mackerel-server-go/domain"

type MonitorInteractor struct {
	MonitorRepository MonitorRepository
}

func (intr *MonitorInteractor) List(orgID string) (*domain.Monitors, error) {
	return intr.MonitorRepository.List(orgID)
}

func (intr *MonitorInteractor) Save(orgID string, monitor *domain.Monitoring) (any, error) {
	monitor.ID = domain.NewIDWith(monitor.Name)

	if monitor.Type == "external" && monitor.Method == "" {
		monitor.Method = "GET"
	}

	return intr.MonitorRepository.Save(orgID, monitor)
}

func (intr *MonitorInteractor) Update(orgID string, monitor *domain.Monitoring) (any, error) {
	return intr.MonitorRepository.Update(orgID, monitor)
}

func (intr *MonitorInteractor) Monitor(orgID, monitorID string) (any, error) {
	return intr.MonitorRepository.Monitor(orgID, monitorID)
}

func (intr *MonitorInteractor) Delete(orgID, monitorID string) (any, error) {
	return intr.MonitorRepository.Delete(orgID, monitorID)
}
