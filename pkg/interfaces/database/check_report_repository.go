package database

import (
	"fmt"
	"strconv"
	"time"

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
				org_id                varchar(64)  not null,
				host_id               varchar(16)  not null,
				type                  enum('host') not null,
				name                  varchar(128) not null,
				status                enum('OK', 'CRITICAL', 'WARNING', 'UNKNOWN') not null,
				message               text,
				occurred_at           bigint,
				notification_interval bigint,
				max_check_attempts    bigint,
				primary key(host_id, name),
				index (org_id, status)
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

func (repo *CheckReportRepository) CheckReport(orgID string) (*domain.CheckReports, error) {
	reports := make([]domain.CheckReport, 0)
	if err := repo.Transact(func(tx Tx) error {
		rows, err := tx.Query("select * from check_reports where org_id=? and status not in('OK')", orgID)
		if err != nil {
			return fmt.Errorf("select * from check_reports: %v", err)
		}

		for rows.Next() {
			var report domain.CheckReport
			var trash string
			if err := rows.Scan(
				&trash,
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

func (repo *CheckReportRepository) Save(orgID string, reports *domain.CheckReports) (*domain.Success, error) {
	if err := repo.Transact(func(tx Tx) error {
		for i := range reports.Reports {
			if _, err := tx.Exec(
				`
				insert into check_reports (
					org_id,
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
					occurred_at = values(occurred_at),
					notification_interval = values(notification_interval),
					max_check_attempts = values(max_check_attempts)
				`,
				orgID,
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

	// update alert
	if err := repo.Transact(func(tx Tx) error {
		for i := range reports.Reports {
			row := tx.QueryRow(
				`select 1 from alerts where org_id=? and monitor_id=? and status!='OK'`,
				orgID,
				domain.NewMonitorID(
					reports.Reports[i].Source.HostID,
					reports.Reports[i].Name,
				),
			)

			var trash int
			if err := row.Scan(&trash); err != nil {
				continue
			}

			var closedAt int64
			if reports.Reports[i].Status == "OK" {
				closedAt = time.Now().Unix()
			}

			if _, err := tx.Exec(`update alerts set status=?, message=?, closed_at=? where org_id=? and monitor_id=?`,
				reports.Reports[i].Status,
				reports.Reports[i].Message,
				closedAt,
				orgID,
				domain.NewMonitorID(
					reports.Reports[i].Source.HostID,
					reports.Reports[i].Name,
				),
			); err != nil {
				return fmt.Errorf("update alerts: %v", err)
			}
		}

		return nil
	}); err != nil {
		return &domain.Success{Success: false}, nil
	}

	// insert alert
	if err := repo.Transact(func(tx Tx) error {
		for i := range reports.Reports {
			if reports.Reports[i].Status == "OK" {
				continue
			}

			row := tx.QueryRow(
				`select 1 from alerts where org_id=? and monitor_id=? and status!='OK'`,
				orgID,
				domain.NewMonitorID(
					reports.Reports[i].Source.HostID,
					reports.Reports[i].Name,
				),
			)

			var trash int
			if err := row.Scan(&trash); err != nil {
				continue
			}

			if _, err := tx.Exec(
				`
					insert into alerts (
						org_id,
						id,
						status,
						monitor_id,
						type,
						host_id,
						value,
						message,
						reason,
						opened_at,
						closed_at
					) values(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
					`,
				orgID,
				domain.NewAlertID(
					reports.Reports[i].Source.HostID,
					reports.Reports[i].Name,
					strconv.FormatInt(reports.Reports[i].OccurredAt, 10),
				),
				reports.Reports[i].Status,
				domain.NewMonitorID(
					reports.Reports[i].Source.HostID,
					reports.Reports[i].Name,
				),
				"check",
				reports.Reports[i].Source.HostID,
				0,
				reports.Reports[i].Message[:len(reports.Reports[i].Message)-1], // remove \n
				"",
				reports.Reports[i].OccurredAt,
				0,
			); err != nil {
				return fmt.Errorf("insert into alerts: %v", err)
			}
		}

		return nil
	}); err != nil {
		return &domain.Success{Success: false}, nil
	}

	return &domain.Success{Success: true}, nil
}
