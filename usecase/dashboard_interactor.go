package usecase

import (
	"errors"

	"github.com/itsubaki/mackerel-server-go/domain"
)

type DashboardInteractor struct {
	DashboardRepository DashboardRepository
}

func (intr *DashboardInteractor) List(orgID string) (*domain.Dashboards, error) {
	return intr.DashboardRepository.List(orgID)
}

func (intr *DashboardInteractor) Save(orgID string, dashboard *domain.Dashboard) (*domain.Dashboard, error) {
	dashboard.ID = domain.NewRandomID()
	return intr.DashboardRepository.Save(orgID, dashboard)
}

func (intr *DashboardInteractor) Dashboard(orgID, dashboardID string) (*domain.Dashboard, error) {
	return intr.DashboardRepository.Dashboard(orgID, dashboardID)
}

func (intr *DashboardInteractor) Update(orgID string, dashboard *domain.Dashboard) (*domain.Dashboard, error) {
	return intr.DashboardRepository.Update(orgID, dashboard)
}

func (intr *DashboardInteractor) Exists(orgID, dashboardID string) bool {
	return intr.DashboardRepository.Exists(orgID, dashboardID)
}

func (intr *DashboardInteractor) Delete(orgID, dashboardID string) (*domain.Dashboard, error) {
	if !intr.DashboardRepository.Exists(orgID, dashboardID) {
		return nil, &DashboardNotFound{Err{errors.New("when the dashboard corresponding to the designated ID can't be found")}}
	}

	return intr.DashboardRepository.Delete(orgID, dashboardID)
}
