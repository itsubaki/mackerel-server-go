package usecase

import "github.com/itsubaki/mackerel-api/pkg/domain"

type AlertInteractor struct {
	AlertRepository AlertRepository
}

func (s *AlertInteractor) List(withClosed bool, nextID string, limit int) (*domain.Alerts, error) {
	return s.AlertRepository.List(withClosed, nextID, limit)
}

func (s *AlertInteractor) Close(alertID, reason string) (*domain.Alert, error) {
	return s.AlertRepository.Close(alertID, reason)
}
