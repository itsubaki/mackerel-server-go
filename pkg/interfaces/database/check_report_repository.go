package database

import (
	"fmt"

	"github.com/itsubaki/mackerel-api/pkg/domain"
)

type CheckReportRepository struct {
	SQLHandler
}

func NewCheckReportRepository(handler SQLHandler) *CheckReportRepository {
	if err := handler.Transact(func(tx Tx) error {
		if _, err := tx.Exec(
			`
			create table if not exists check_reports (
				org                   varchar(64)  not null,
				host_id               varchar(16)  not null,
				type                  enum('host') not null,
				name                  varchar(128) not null,
				status                enum('OK', 'CRITICAL', 'WARNING', 'UNKNOWN') not null,
				message               text,
				occurred_at           bigint,
				notification_interval bigint,
				max_check_attempts    bigint,
				primary key(host_id, name)
			)
			`,
		); err != nil {
			return fmt.Errorf("create table check_reports: %v", err)
		}

		return nil
	}); err != nil {
		panic(fmt.Errorf("transaction: %v", err))
	}

	return &CheckReportRepository{
		SQLHandler: handler,
	}
}

func (repo *CheckReportRepository) CheckReport(org string) (*domain.CheckReports, error) {
	reports := make([]domain.CheckReport, 0)
	if err := repo.Transact(func(tx Tx) error {
		rows, err := tx.Query("select * from check_reports where org=? and status not in('OK')", org)
		if err != nil {
			return fmt.Errorf("select * from check_reports: %v", err)
		}

		for rows.Next() {
			var report domain.CheckReport
			var org string
			if err := rows.Scan(
				&org,
				&report.Source.HostID,
				&report.Source.Type,
				&report.Name,
				&report.Status,
				&report.Message,
				&report.OccurredAt,
				&report.NotificationInterval,
				&report.MaxCheckAttempts,
			); err != nil {
				return fmt.Errorf("scan: %v", err)
			}

			reports = append(reports, report)
		}

		return nil
	}); err != nil {
		return nil, fmt.Errorf("transaction: %v", err)
	}

	return &domain.CheckReports{Reports: reports}, nil
}

func (repo *CheckReportRepository) Save(org string, reports *domain.CheckReports) (*domain.Success, error) {
	if err := repo.Transact(func(tx Tx) error {
		for i := range reports.Reports {
			if _, err := tx.Exec(
				`
				insert into check_reports (
					org,
					host_id,
					type,
					name,
					status,
					message,
					occurred_at,
					notification_interval,
					max_check_attempts
				)
				values (?, ?, ?, ?, ?, ?, ?, ?, ?)
				on duplicate key update
					status = values(status),
					message = values(message),
					notification_interval = values(notification_interval),
					max_check_attempts = values(max_check_attempts)
				`,
				org,
				reports.Reports[i].Source.HostID,
				reports.Reports[i].Source.Type,
				reports.Reports[i].Name,
				reports.Reports[i].Status,
				reports.Reports[i].Message[:len(reports.Reports[i].Message)-1], // remove \n
				reports.Reports[i].OccurredAt,
				reports.Reports[i].NotificationInterval,
				reports.Reports[i].MaxCheckAttempts,
			); err != nil {
				return fmt.Errorf("insert into check_reports: %v", err)
			}
		}

		return nil
	}); err != nil {
		return &domain.Success{Success: false}, nil
	}

	return &domain.Success{Success: true}, nil
}
