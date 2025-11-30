package usecase

import (
	"errors"
	"fmt"

	"github.com/itsubaki/mackerel-server-go/domain"
)

type AlertInteractor struct {
	AlertRepository AlertRepository
}

func (intr *AlertInteractor) List(orgID string, withClosed bool, nextID string, limit int) (*domain.Alerts, error) {
	if limit > 100 {
		return nil, &AlertLimitOver{Err{errors.New("`limit` value is larger than maximum allowed value(100)")}}
	}

	return intr.AlertRepository.List(orgID, withClosed, nextID, limit)
}

func (intr *AlertInteractor) Close(orgID, alertID, reason string) (*domain.Alert, error) {
	if !intr.AlertRepository.Exists(orgID, alertID) {
		return nil, &AlertNotFound{
			Err{
				Err: fmt.Errorf("the <%s>'s corresponding alert can't be found", alertID),
			},
		}
	}

	return intr.AlertRepository.Close(orgID, alertID, reason)
}
