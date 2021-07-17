package database

import (
	"fmt"

	"github.com/itsubaki/mackerel-server-go/pkg/domain"
	"github.com/itsubaki/mackerel-server-go/pkg/usecase"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var _ usecase.CheckReportRepository = (*CheckReportRepository)(nil)

type CheckReportRepository struct {
	DB *gorm.DB
}

type CheckReport struct {
	OrgID                string `gorm:"column:org_id;     type:varchar(16);  not null"`
	HostID               string `gorm:"column:host_id;    type:varchar(16);  not null; primary_key"`
	Type                 string `gorm:"column:type;       type:enum('host'); not null"`
	Name                 string `gorm:"column:name;       type:varchar(128); not null; primary_key"`
	Status               string `gorm:"column:status;     type:enum('OK', 'CRITICAL', 'WARNING', 'UNKNOWN'); not null"`
	Message              string `gorm:"column:message;    type:text;"`
	OccurredAt           int64  `gorm:"column:occurred_at;           type:bigint;"`
	NotificationInterval int64  `gorm:"column:notification_interval; type:bigint;"`
	MaxCheckAttempts     int64  `gorm:"column:max_check_attempts;    type:bigint;"`
}

func NewCheckReportRepository(handler SQLHandler) *CheckReportRepository {
	db, err := gorm.Open(mysql.New(mysql.Config{
		Conn: handler.Raw().(gorm.ConnPool),
	}), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	if handler.IsDebugMode() {
		db.Logger.LogMode(4)
	}

	if err := db.AutoMigrate(&CheckReport{}); err != nil {
		panic(fmt.Errorf("auto migrate check_report: %v", err))
	}

	return &CheckReportRepository{
		DB: db,
	}
}

func (r *CheckReportRepository) CheckReport(orgID string) (*domain.CheckReports, error) {
	result := make([]CheckReport, 0)
	if err := r.DB.Where(&CheckReport{OrgID: orgID}).Not("status", []string{"OK"}).Find(&result).Error; err != nil {
		return nil, fmt.Errorf("select * from check_reports: %v", err)
	}

	out := make([]domain.CheckReport, 0)
	for _, r := range result {
		out = append(out, domain.CheckReport{
			OrgID: r.OrgID,
			Source: domain.Source{
				HostID: r.HostID,
				Type:   r.Type,
			},
			Name:                 r.Name,
			Status:               r.Status,
			Message:              r.Message,
			OccurredAt:           r.OccurredAt,
			NotificationInterval: r.NotificationInterval,
			MaxCheckAttempts:     r.MaxCheckAttempts,
		})
	}

	return &domain.CheckReports{Reports: out}, nil
}

func (r *CheckReportRepository) Save(orgID string, reports *domain.CheckReports) (*domain.Success, error) {
	if err := r.DB.Transaction(func(tx *gorm.DB) error {
		for _, report := range reports.Reports {
			where := CheckReport{
				OrgID:  orgID,
				HostID: report.Source.HostID,
				Type:   report.Source.Type,
				Name:   report.Name,
			}

			update := CheckReport{
				Status:               report.Status,
				Message:              report.Message,
				OccurredAt:           report.OccurredAt,
				NotificationInterval: report.NotificationInterval,
				MaxCheckAttempts:     report.MaxCheckAttempts,
			}

			if err := tx.Where(&where).Assign(&update).FirstOrCreate(&CheckReport{}).Error; err != nil {
				return fmt.Errorf("firts or create: %v", err)
			}
		}

		return nil
	}); err != nil {
		return &domain.Success{Success: false}, fmt.Errorf("transaction: %v", err)
	}

	return &domain.Success{Success: true}, nil
}
