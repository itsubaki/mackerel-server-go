package usecase

import (
	"errors"

	"github.com/itsubaki/mackerel-api/pkg/domain"
)

type AlertInteractor struct {
	AlertRepository AlertRepository
}

func (s *AlertInteractor) List(orgID string, withClosed bool, nextID string, limit int) (*domain.Alerts, error) {
	if limit > 100 {
		return nil, &AlertLimitOver{Err{errors.New("`limit` value is larger than maximum allowed value(100)")}}
	}

	return s.AlertRepository.List(orgID, withClosed, nextID, limit)
}

func (s *AlertInteractor) Close(orgID, alertID, reason string) (*domain.Alert, error) {
	if !s.AlertRepository.Exists(orgID, alertID) {
		return nil, &AlertNotFound{Err{errors.New("the <alertId>'s corresponding alert can't be found")}}
	}

	return s.AlertRepository.Close(orgID, alertID, reason)
}
