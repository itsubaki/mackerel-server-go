package database

import (
	"fmt"

	"github.com/jinzhu/gorm"

	"github.com/itsubaki/mackerel-server-go/pkg/domain"
)

type CheckReportRepository struct {
	DB *gorm.DB
}

type CheckReport struct {
	OrgID                string `gorm:"column:org_id;     type:varchar(16);  not null"`
	HostID               string `gorm:"column:host_id;    type:varchar(16);  not null; primary key"`
	Type                 string `gorm:"column:type;       type:enum('host'); not null"`
	Name                 string `gorm:"column:name;       type:varchar(128); not null; primary key"`
	Status               string `gorm:"column:status;     type:enum('OK', 'CRITICAL', 'WARNING', 'UNKNOWN'); not null"`
	Message              string `gorm:"column:message;    type:text;"`
	OccurredAt           int64  `gorm:"column:occurred_at;           type:bigint;"`
	NotificationInterval int64  `gorm:"column:notification_interval; type:bigint;"`
	MaxCheckAttempts     int64  `gorm:"column:max_check_attempts;    type:bigint;"`
}

func NewCheckReportRepository(handler SQLHandler) *CheckReportRepository {
	db, err := gorm.Open(handler.Dialect(), handler.Raw())
	if err != nil {
		panic(err)
	}
	db.LogMode(handler.IsDebugging())

	if err := db.AutoMigrate(&CheckReport{}).Error; err != nil {
		panic(fmt.Errorf("auto migrate check_report: %v", err))
	}

	return &CheckReportRepository{
		DB: db,
	}
}

func (repo *CheckReportRepository) CheckReport(orgID string) (*domain.CheckReports, error) {
	result := make([]CheckReport, 0)
	if err := repo.DB.Where(&CheckReport{OrgID: orgID}).Not("status", []string{"OK"}).Find(&result).Error; err != nil {
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

func (repo *CheckReportRepository) Save(orgID string, reports *domain.CheckReports) (*domain.Success, error) {
	if err := repo.DB.Transaction(func(tx *gorm.DB) error {
		for _, r := range reports.Reports {
			where := CheckReport{
				OrgID:  orgID,
				HostID: r.Source.HostID,
				Type:   r.Source.Type,
				Name:   r.Name,
			}

			update := CheckReport{
				Status:               r.Status,
				Message:              r.Message,
				OccurredAt:           r.OccurredAt,
				NotificationInterval: r.NotificationInterval,
				MaxCheckAttempts:     r.MaxCheckAttempts,
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
