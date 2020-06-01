package database

import (
	"fmt"

	"github.com/jinzhu/gorm"

	"github.com/itsubaki/mackerel-server-go/pkg/domain"
)

type DashboardRepository struct {
	DB *gorm.DB
}

type Dashboard struct {
	OrgID     string `gorm:"column:org_id;     type:varchar(16);  not null"`
	ID        string `gorm:"column:id;         type:varchar(128); not null; primary key"`
	Title     string `gorm:"column:title;      type:varchar(128); not null"`
	Memo      string `gorm:"column:memo;       type:varchar(128); not null; default:''"`
	URLPath   string `gorm:"column:url_path;   type:text"`
	CreatedAt int64  `gorm:"column:created_at; type:bigint"`
	UpdatedAt int64  `gorm:"column:updated_at; type:bigint"`
}

func NewDashboardRepository(handler SQLHandler) *DashboardRepository {
	db, err := gorm.Open(handler.Dialect(), handler.Raw())
	if err != nil {
		panic(err)
	}
	db.LogMode(handler.IsDebugging())

	if err := db.AutoMigrate(&Dashboard{}).Error; err != nil {
		panic(fmt.Errorf("auto migrate dashboard: %v", err))
	}

	return &DashboardRepository{
		DB: db,
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
