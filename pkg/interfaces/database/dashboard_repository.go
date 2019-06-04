package database

import "github.com/itsubaki/mackerel-api/pkg/domain"

type DashboardRepository struct {
	SQLHandler
}

func NewDashboardRepository(handler SQLHandler) *DashboardRepository {
	return &DashboardRepository{
		SQLHandler: handler,
	}
}

func (repo *DashboardRepository) List(orgID string) (*domain.Dashboards, error) {
	return nil, nil
}

func (repo *DashboardRepository) Save(orgID string, dashboard *domain.Dashboard) (*domain.Dashboard, error) {
	return nil, nil
}

func (repo *DashboardRepository) Dashboard(orgID, dashboardID string) (*domain.Dashboard, error) {
	return nil, nil
}

func (repo *DashboardRepository) Update(orgID string, dashboard *domain.Dashboard) (*domain.Dashboard, error) {
	return nil, nil
}

func (repo *DashboardRepository) Exists(orgID, dashboardID string) bool {
	return true
}

func (repo *DashboardRepository) Delete(orgID, dashboardID string) (*domain.Dashboard, error) {
	return nil, nil
}
