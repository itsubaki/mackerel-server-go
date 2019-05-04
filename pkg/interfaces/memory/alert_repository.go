package memory

import (
	"fmt"
	"time"

	"github.com/itsubaki/mackerel-api/pkg/domain"
)

type AlertRepository struct {
	Alerts []domain.Alert
}

// select * from alerts where id=${alertID} limit=1
func (repo *AlertRepository) Exists(alertID string) bool {
	for _, a := range repo.Alerts {
		if a.ID == alertID {
			return true
		}
	}

	return false
}

// select * from alerts where closed=${withClosed} order by openedAt desc limit ${limit}
func (repo *AlertRepository) List(withClosed bool, nextID string, limit int) (*domain.Alerts, error) {
	return &domain.Alerts{
		Alerts: repo.Alerts,
		NextID: "",
	}, nil
}

// update alerts set reason=${reason}, status=OK, closedAt=time.Now().Unix() where id=${alertID}
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
