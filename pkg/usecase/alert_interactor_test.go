package usecase_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/itsubaki/mackerel-server-go/pkg/domain"
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

func TestAlertInteractorList(t *testing.T) {
	intr := &usecase.AlertInteractor{
		AlertRepository: NewAlertRepositoryMock([]domain.Alert{
			{OrgID: "foo", ID: domain.NewRandomID()},
			{OrgID: "bar", ID: domain.NewRandomID()},
		}),
	}

	cases := []struct {
		limit   int
		message string
	}{
		{10, ""},
		{100, ""},
		{101, "`limit` value is larger than maximum allowed value(100)"},
	}

	for _, c := range cases {
		_, err := intr.List("foo", false, "", c.limit)
		if c.limit < 100 && err != nil {
			t.Fail()
		}

		if c.limit > 100 && err != nil && err.Error() != c.message {
			t.Fail()
		}
	}
}

func TestAlertInteractorClose(t *testing.T) {
	intr := &usecase.AlertInteractor{
		AlertRepository: NewAlertRepositoryMock([]domain.Alert{
			{OrgID: "foo", ID: "12345"},
			{OrgID: "bar", ID: domain.NewRandomID()},
		}),
	}

	cases := []struct {
		orgID   string
		alertID string
		message string
	}{
		{"foo", "12345", ""},
		{"piy", "23456", "the <23456>'s corresponding alert can't be found"},
	}

	for _, c := range cases {
		if _, err := intr.Close(c.orgID, c.alertID, "hoge"); err != nil {
			if err.Error() != c.message {
				t.Fail()
			}
		}
	}
}
