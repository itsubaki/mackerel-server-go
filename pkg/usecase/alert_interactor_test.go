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

func (repo *AlertRepositoryMock) Exists(orgID, alertID string) bool {
	for _, a := range repo.Alerts {
		if a.OrgID == orgID && a.ID == alertID {
			return true
		}
	}

	return false
}

func (repo *AlertRepositoryMock) Save(orgID string, alert *domain.Alert) (*domain.Alert, error) {
	alert.OrgID = orgID
	repo.Alerts = append(repo.Alerts, *alert)
	return alert, nil
}

func (repo *AlertRepositoryMock) List(orgID string, withClosed bool, nextID string, limit int) (*domain.Alerts, error) {
	alerts := make([]domain.Alert, 0)
	for i, a := range repo.Alerts {
		if a.OrgID == orgID {
			alerts = append(alerts, repo.Alerts[i])
		}
	}

	return &domain.Alerts{
		Alerts: alerts,
	}, nil
}

func (repo *AlertRepositoryMock) Close(orgID, alertID, reason string) (*domain.Alert, error) {
	for i, a := range repo.Alerts {
		if a.OrgID == orgID && a.ID == alertID {
			repo.Alerts[i].Status = "OK"
			repo.Alerts[i].Reason = reason
			repo.Alerts[i].ClosedAt = time.Now().Unix()
			return &repo.Alerts[i], nil
		}
	}

	return nil, fmt.Errorf("alert=%s/%s not found", orgID, alertID)
}

func TestAlertInteractorList(t *testing.T) {
	ai := &usecase.AlertInteractor{
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
		_, err := ai.List("foo", false, "", c.limit)
		if c.limit < 100 && err != nil {
			t.Fail()
		}

		if c.limit > 100 && err != nil && err.Error() != c.message {
			t.Fail()
		}
	}
}

func TestAlertInteractorClose(t *testing.T) {
	ai := &usecase.AlertInteractor{
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
		if _, err := ai.Close(c.orgID, c.alertID, "hoge"); err != nil {
			if err.Error() != c.message {
				t.Fail()
			}
		}
	}
}
