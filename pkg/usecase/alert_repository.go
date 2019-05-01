package usecase

import "github.com/itsubaki/mackerel-api/pkg/domain"

type AlertRepository interface {
	Exists(alertID string) bool
	List(withClosed bool, nextID string, limit int) (*domain.Alerts, error)
	Close(alertID, reason string) (*domain.Alert, error)
}
