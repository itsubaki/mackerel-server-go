package usecase

import "github.com/itsubaki/mackerel-server-go/domain"

type DashboardRepository interface {
	List(orgID string) (*domain.Dashboards, error)
	Save(orgID string, dashboard *domain.Dashboard) (*domain.Dashboard, error)
	Dashboard(orgID, dashboardID string) (*domain.Dashboard, error)
	Update(orgID string, dashboard *domain.Dashboard) (*domain.Dashboard, error)
	Exists(orgID, dashboardID string) bool
	Delete(orgID, dashboardID string) (*domain.Dashboard, error)
}
