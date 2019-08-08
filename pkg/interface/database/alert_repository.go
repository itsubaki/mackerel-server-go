package database

import (
	"fmt"
	"log"
	"time"

	"github.com/itsubaki/mackerel-api/pkg/domain"
)

type AlertRepository struct {
	SQLHandler
}

func NewAlertRepository(handler SQLHandler) *AlertRepository {
	if err := handler.Transact(func(tx Tx) error {
		if _, err := tx.Exec(
			`
			create table if not exists alerts (
				org_id     varchar(16) not null,
				id         varchar(16) not null primary key,
				status     enum('OK', 'CRITICAL', 'WARNING', 'UNKNOWN') not null,
				monitor_id varchar(16) not null,
				type       enum('connectivity', 'host', 'service', 'external', 'check', 'expression') not null,
				host_id    varchar(16),
				value      double,
				message    text,
				reason     text,
				opened_at  bigint,
				closed_at  bigint,
				index(org_id, opened_at desc)
			)
			`,
		); err != nil {
			return fmt.Errorf("create table alerts: %v", err)
		}

		if _, err := tx.Exec(
			`
			create table if not exists alert_history (
				org_id     varchar(16) not null,
				alert_id   varchar(16) not null,
				status     enum('OK', 'CRITICAL', 'WARNING', 'UNKNOWN') not null,
				monitor_id varchar(16) not null,
				host_id    varchar(16),
				time       bigint      not null,
				message    text,
				primary key(alert_id, monitor_id, time desc)
			)
			`,
		); err != nil {
			return fmt.Errorf("create table alert_history: %v", err)
		}

		if _, err := tx.Exec(
			`
			create table if not exists alert_history_latest (
				org_id     varchar(16) not null,
				alert_id   varchar(16) not null,
				status     enum('OK', 'CRITICAL', 'WARNING', 'UNKNOWN') not null,
				monitor_id varchar(16) not null,
				host_id    varchar(16),
				time       bigint      not null,
				message    text,
				primary key(alert_id, monitor_id)
			)
			`,
		); err != nil {
			return fmt.Errorf("create table alert_history_latest: %v", err)
		}

		return nil
	}); err != nil {
		panic(fmt.Errorf("transaction: %v", err))
	}

	return &AlertRepository{
		SQLHandler: handler,
	}
}

func (repo *AlertRepository) Exists(orgID, alertID string) bool {
	rows, err := repo.Query("select 1 from alerts where org_id=? and id=?", orgID, alertID)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	if rows.Next() {
		return true
	}

	return false
}

func (repo *AlertRepository) List(orgID string, withClosed bool, nextID string, limit int) (*domain.Alerts, error) {
	status := "UNKNOWN"
	if withClosed {
		status = "OK"
	}

	rows, err := repo.Query(
		`
		select * from alerts where org_id=? and status in ('CRITICAL', 'WARNING', 'UNKNOWN', ?) order by opened_at desc limit ?
		`,
		orgID,
		status,
		limit+1,
	)
	if err != nil {
		return nil, fmt.Errorf("select * from alerts: %v", err)
	}
	defer rows.Close()

	alerts := make([]domain.Alert, 0)
	for rows.Next() {
		var alert domain.Alert
		if err := rows.Scan(
			&alert.OrgID,
			&alert.ID,
			&alert.Status,
			&alert.MonitorID,
			&alert.Type,
			&alert.HostID,
			&alert.Value,
			&alert.Message,
			&alert.Reason,
			&alert.OpenedAt,
			&alert.ClosedAt,
		); err != nil {
			return nil, fmt.Errorf("scan: %v", err)
		}

		alerts = append(alerts, alert)
	}

	if len(alerts) > limit {
		return &domain.Alerts{
			Alerts: alerts[:len(alerts)-1],
			NextID: alerts[len(alerts)-1].ID,
		}, nil
	}

	return &domain.Alerts{Alerts: alerts}, nil
}

func (repo *AlertRepository) Close(orgID, alertID, reason string) (*domain.Alert, error) {
	var alert domain.Alert
	if err := repo.Transact(func(tx Tx) error {
		if _, err := tx.Exec(
			`
			update alerts set status='OK', reason=?, closed_at=? where org_id=? and id=?
			`,
			reason,
			time.Now().Unix(),
			orgID,
			alertID,
		); err != nil {
			return fmt.Errorf("update alerts: %v", err)
		}

		row := tx.QueryRow("select * from alerts where org_id=? and id=?", orgID, alertID)
		if err := row.Scan(
			&alert.OrgID,
			&alert.ID,
			&alert.Status,
			&alert.MonitorID,
			&alert.Type,
			&alert.HostID,
			&alert.Value,
			&alert.Message,
			&alert.Reason,
			&alert.OpenedAt,
			&alert.ClosedAt,
		); err != nil {
			return fmt.Errorf("scan: %v", err)
		}

		if _, err := tx.Exec(
			`
				insert into alert_history (
					org_id,
					alert_id,
					status,
					monitor_id,
					host_id,
					time,
					message
				) values (?, ?, ?, ?, ?, ?, ?)
				`,
			orgID,
			alertID,
			"OK",
			alert.MonitorID,
			alert.HostID,
			alert.ClosedAt,
			alert.Message,
		); err != nil {
			return fmt.Errorf("insert into alert_history: %v", err)
		}

		if _, err := tx.Exec(
			`
				insert into alert_history_latest (
					org_id,
					alert_id,
					status,
					monitor_id,
					host_id,
					time,
					message
				) values (?, ?, ?, ?, ?, ?, ?)
				on duplicate key update
					status = values(status),
					time = values(time),
					message = values(message)
				`,
			orgID,
			alertID,
			"OK",
			alert.MonitorID,
			alert.HostID,
			alert.ClosedAt,
			alert.Message,
		); err != nil {
			return fmt.Errorf("insert into alert_history_latest: %v", err)
		}

		return nil
	}); err != nil {
		log.Printf("transaction: %v\n", err)
		return nil, fmt.Errorf("transaction: %v", err)
	}

	return &alert, nil
}
