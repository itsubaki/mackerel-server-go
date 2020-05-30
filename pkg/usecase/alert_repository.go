package usecase

import "github.com/itsubaki/mackerel-server-go/pkg/domain"

type AlertRepository interface {
	Exists(orgID, alertID string) bool
	Save(orgID string, alert *domain.Alert) (*domain.Alert, error)
	List(orgID string, withClosed bool, nextID string, limit int) (*domain.Alerts, error)
	Close(orgID, alertID, reason string) (*domain.Alert, error)
}
