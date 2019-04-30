package usecase

import (
	"github.com/itsubaki/mackerel-api/pkg/domain"
	"github.com/itsubaki/mackerel-api/pkg/interfaces/database"
)

func NewAlertInteractor() *AlertInteractor {
	return &AlertInteractor{
		AlertRepository: database.NewAlertRepository(),
	}
}

type AlertInteractor struct {
	AlertRepository *database.AlertRepository
}

func (s *AlertInteractor) FindBy(withClosed bool, nextID string, limit int64) (domain.Alerts, string, error) {
	return domain.Alerts{}, "nextID", nil
}

func (s *AlertInteractor) Save(alertID, reason string) (*domain.Alert, error) {
	return &domain.Alert{}, nil
}
