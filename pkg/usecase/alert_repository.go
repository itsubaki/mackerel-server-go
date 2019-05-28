package usecase

import "github.com/itsubaki/mackerel-api/pkg/domain"

type AlertRepository interface {
	Exists(orgID, alertID string) bool
	List(orgID string, withClosed bool, nextID string, limit int) (*domain.Alerts, error)
	Close(orgID, alertID, reason string) (*domain.Alert, error)
}
