package usecase

import "github.com/itsubaki/mackerel-api/pkg/domain"

type AlertRepository interface {
	Exists(org, alertID string) bool
	List(org string, withClosed bool, nextID string, limit int) (*domain.Alerts, error)
	Close(org, alertID, reason string) (*domain.Alert, error)
}
