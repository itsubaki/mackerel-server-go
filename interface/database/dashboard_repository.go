package database

import (
	"fmt"

	"github.com/itsubaki/mackerel-server-go/domain"
	"github.com/itsubaki/mackerel-server-go/usecase"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var _ usecase.DashboardRepository = (*DashboardRepository)(nil)

type DashboardRepository struct {
	DB *gorm.DB
}

type Dashboard struct {
	OrgID     string `gorm:"column:org_id;     type:varchar(16);  not null"`
	ID        string `gorm:"column:id;         type:varchar(128); not null; primary_key"`
	Title     string `gorm:"column:title;      type:varchar(128); not null"`
	Memo      string `gorm:"column:memo;       type:varchar(128); not null; default:''"`
	URLPath   string `gorm:"column:url_path;   type:text"`
	CreatedAt int64  `gorm:"column:created_at; type:bigint"`
	UpdatedAt int64  `gorm:"column:updated_at; type:bigint"`
}

func NewDashboardRepository(handler SQLHandler) *DashboardRepository {
	db, err := gorm.Open(mysql.New(mysql.Config{
		Conn: handler.Raw().(gorm.ConnPool),
	}), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	if handler.IsDebugMode() {
		db.Logger.LogMode(4)
	}

	if err := db.AutoMigrate(&Dashboard{}); err != nil {
		panic(fmt.Errorf("auto migrate dashboard: %v", err))
	}

	return &DashboardRepository{
		DB: db,
	}
}

func (r *DashboardRepository) List(orgID string) (*domain.Dashboards, error) {
	return nil, nil
}

func (r *DashboardRepository) Save(orgID string, dashboard *domain.Dashboard) (*domain.Dashboard, error) {
	return nil, nil
}

func (r *DashboardRepository) Dashboard(orgID, dashboardID string) (*domain.Dashboard, error) {
	return nil, nil
}

func (r *DashboardRepository) Update(orgID string, dashboard *domain.Dashboard) (*domain.Dashboard, error) {
	return nil, nil
}

func (r *DashboardRepository) Exists(orgID, dashboardID string) bool {
	return true
}

func (r *DashboardRepository) Delete(orgID, dashboardID string) (*domain.Dashboard, error) {
	return nil, nil
}
