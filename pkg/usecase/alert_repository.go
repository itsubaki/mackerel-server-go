package usecase

import "github.com/itsubaki/mackerel-api/pkg/domain"

type AlertRepository interface {
	Exists(orgID, alertID string) bool
	Save(orgID string, alert *domain.Alert) (*domain.Alert, error)
	Alert(orgID, hostID, monitorID string) (*domain.Alert, error)
	List(orgID string, withClosed bool, nextID string, limit int) (*domain.Alerts, error)
	Close(orgID, alertID, reason string) (*domain.Alert, error)
}
