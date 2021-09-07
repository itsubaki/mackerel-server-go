package controller_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/itsubaki/mackerel-server-go/pkg/domain"
	"github.com/itsubaki/mackerel-server-go/pkg/interface/controller"
	"github.com/itsubaki/mackerel-server-go/pkg/usecase"
)

type AlertRepositoryMock struct {
	Alerts []domain.Alert
}

func NewAlertRepositoryMock(alerts []domain.Alert) *AlertRepositoryMock {
	return &AlertRepositoryMock{
		Alerts: alerts,
	}
}

func (r *AlertRepositoryMock) Exists(orgID, alertID string) bool {
	for _, a := range r.Alerts {
		if a.OrgID == orgID && a.ID == alertID {
			return true
		}
	}

	return false
}

func (r *AlertRepositoryMock) Save(orgID string, alert *domain.Alert) (*domain.Alert, error) {
	alert.OrgID = orgID
	r.Alerts = append(r.Alerts, *alert)
	return alert, nil
}

func (r *AlertRepositoryMock) List(orgID string, withClosed bool, nextID string, limit int) (*domain.Alerts, error) {
	alerts := make([]domain.Alert, 0)
	for i, a := range r.Alerts {
		if a.OrgID == orgID {
			alerts = append(alerts, r.Alerts[i])
		}
	}

	return &domain.Alerts{
		Alerts: alerts,
	}, nil
}

func (r *AlertRepositoryMock) Close(orgID, alertID, reason string) (*domain.Alert, error) {
	for i, a := range r.Alerts {
		if a.OrgID == orgID && a.ID == alertID {
			r.Alerts[i].Status = "OK"
			r.Alerts[i].Reason = reason
			r.Alerts[i].ClosedAt = time.Now().Unix()
			return &r.Alerts[i], nil
		}
	}

	return nil, fmt.Errorf("alert=%s/%s not found", orgID, alertID)
}

func TestAlertControllerClose(t *testing.T) {
	cntr := &controller.AlertController{
		Interactor: &usecase.AlertInteractor{
			AlertRepository: &AlertRepositoryMock{[]domain.Alert{
				{OrgID: "foo", ID: "hoge"},
				{OrgID: "bar", ID: domain.NewRandomID()},
			}},
		},
	}

	cases := []struct {
		orgID, alertID string
		reason         string
		want           int
	}{
		{"foo", "hoge", "{\"reason\": \"closed manually\"}", 200},
		{"piy", "hoge", "{\"reason\": \"closed manually\"}", 404},
		{"piy", "hoge", "invalid request body", 400},
	}

	for _, c := range cases {
		ctx := Context()
		ctx.SetParam("alertId", c.alertID)
		ctx.Set("org_id", c.orgID)
		ctx.SetRequestBody([]byte(c.reason))
		cntr.Close(ctx)

		got := ctx.GetStatus()
		if got != c.want {
			t.Fail()
		}
	}
}
