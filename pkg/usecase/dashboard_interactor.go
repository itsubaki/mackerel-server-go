package usecase

import (
	"errors"

	"github.com/itsubaki/mackerel-api/pkg/domain"
)

type DashboardInteractor struct {
	DashboardRepository DashboardRepository
}

func (s *DashboardInteractor) List(orgID string) (*domain.Dashboards, error) {
	return s.DashboardRepository.List(orgID)
}

func (s *DashboardInteractor) Save(orgID string, dashboard *domain.Dashboard) (*domain.Dashboard, error) {
	dashboard.ID = domain.NewID()
	return s.DashboardRepository.Save(orgID, dashboard)
}

func (s *DashboardInteractor) Dashboard(orgID, dashboardID string) (*domain.Dashboard, error) {
	return s.DashboardRepository.Dashboard(orgID, dashboardID)
}

func (s *DashboardInteractor) Update(orgID string, dashboard *domain.Dashboard) (*domain.Dashboard, error) {
	return s.DashboardRepository.Update(orgID, dashboard)
}

func (s *DashboardInteractor) Exists(orgID, dashboardID string) bool {
	return s.DashboardRepository.Exists(orgID, dashboardID)
}

func (s *DashboardInteractor) Delete(orgID, dashboardID string) (*domain.Dashboard, error) {
	if !s.DashboardRepository.Exists(orgID, dashboardID) {
		return nil, &DashboardNotFound{Err{errors.New("when the dashboard corresponding to the designated ID can't be found")}}
	}

	return s.DashboardRepository.Delete(orgID, dashboardID)
}
