package usecase

import "github.com/itsubaki/mackerel-api/pkg/domain"

type AlertRepository interface {
	List(withClosed bool, nextID string, limit int) (*domain.Alerts, error)
	Exists(alertID string) bool
	Close(alertID, reason string) (*domain.Alert, error)
}
