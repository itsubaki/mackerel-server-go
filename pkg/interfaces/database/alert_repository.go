package database

import "github.com/itsubaki/mackerel-api/pkg/domain"

type AlertRepository struct {
	Internal domain.Alerts
}

func NewAlertRepository() *AlertRepository {
	return &AlertRepository{
		Internal: domain.Alerts{},
	}
}
