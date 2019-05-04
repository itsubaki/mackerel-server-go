package memory

import (
	"fmt"
	"time"

	"github.com/itsubaki/mackerel-api/pkg/domain"
)

type AlertRepository struct {
	Alerts []domain.Alert
}

func (repo *AlertRepository) Exists(alertID string) bool {
	for _, a := range repo.Alerts {
		if a.ID == alertID {
			return true
		}
	}

	return false
}

func (repo *AlertRepository) List(withClosed bool, nextID string, limit int) (*domain.Alerts, error) {
	return &domain.Alerts{
		Alerts: repo.Alerts,
		NextID: "",
	}, nil
}

func (repo *AlertRepository) Close(alertID, reason string) (*domain.Alert, error) {
	for i := range repo.Alerts {
		if repo.Alerts[i].ID == alertID {
			repo.Alerts[i].Reason = reason
			repo.Alerts[i].Status = "OK"
			repo.Alerts[i].ClosedAt = time.Now().Unix()
			return &repo.Alerts[i], nil
		}
	}

	return nil, fmt.Errorf("alert not found")
}
