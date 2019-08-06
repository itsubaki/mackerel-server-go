package database

import (
	"fmt"

	"github.com/itsubaki/mackerel-api/pkg/domain"
)

type DashboardRepository struct {
	SQLHandler
}

func NewDashboardRepository(handler SQLHandler) *DashboardRepository {
	if err := handler.Transact(func(tx Tx) error {
		if _, err := tx.Exec(
			`
			create table if not exists dashboards (
				org_id   varchar(16)  not null,
				id       varchar(128) not null primary key,
				title    varchar(128) not null,
				memo     varchar(128) not null default '',
				url_path text,
				created_at bigint,
				updated_at bigint
			)
			`,
		); err != nil {
			return fmt.Errorf("create table dashboards: %v", err)
		}

		if _, err := tx.Exec(
			`
			create table if not exists dashboard_widgets (
				org_id       varchar(16)  not null,
				dashboard_id varchar(128) not null
			)
			`,
		); err != nil {
			return fmt.Errorf("create table dashboard_widgets: %v", err)
		}

		return nil
	}); err != nil {
		panic(fmt.Errorf("transaction: %v", err))
	}

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
