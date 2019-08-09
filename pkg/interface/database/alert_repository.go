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

func (repo *AlertRepository) Save(orgID string, alert *domain.Alert) (*domain.Alert, error) {
	if err := repo.Transact(func(tx Tx) error {
		row := tx.QueryRow(
			`
				select
					alert_id,
					status
				from
					alert_history_latest
				where
					org_id=?  and
					host_id=? and
					monitor_id=?
				`,
			orgID,
			alert.HostID,
			alert.MonitorID,
		)

		var alertID, status string
		if err := row.Scan(&alertID, &status); err != nil && alert.Status == "OK" {
			// no record and no alert
			return nil
		}

		if status == "OK" && alert.Status == "OK" {
			// have record and alert closed
			return nil
		}

		if (len(status) < 1 || status == "OK") && alert.Status != "OK" {
			// new alert
			alertID = alert.ID
		}

		// status != "OK" && reports.Reports[i].Status != "OK"
		// -> continuous alert
		// status != "OK" && reports.Reports[i].Status == "OK"
		// -> close alert

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
			alert.Status,
			alert.MonitorID,
			alert.HostID,
			alert.OpenedAt,
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
			alert.Status,
			alert.MonitorID,
			alert.HostID,
			alert.OpenedAt,
			alert.Message,
		); err != nil {
			return fmt.Errorf("insert into alert_history_latest: %v", err)
		}

		return nil
	}); err != nil {
		log.Printf("transaction: %v\n", err)
		return alert, fmt.Errorf("transaction: %v\n", err)
	}

	if err := repo.Transact(func(tx Tx) error {
		row := tx.QueryRow(
			`
				select
					alert_id,
					status,
					host_id,
					message,
					time
				from
					alert_history_latest
				where
					org_id=?  and
					host_id=? and
					monitor_id=?
				`,
			orgID,
			alert.HostID,
			alert.MonitorID,
		)

		var alertID, status, hostID, message string
		var time int64
		if err := row.Scan(&alertID, &status, &hostID, &message, &time); err != nil {
			// no record
			return nil
		}

		var closedAt int64
		if alert.Status == "OK" {
			closedAt = alert.OpenedAt
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
				)
				values (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
				on duplicate key update
					status = values(status),
					message = values(message),
					closed_at = values(closed_at)
				`,
			orgID,
			alertID,
			status,
			alert.MonitorID,
			alert.Type,
			hostID,
			alert.Value,
			message,
			alert.Reason,
			time,
			closedAt,
		); err != nil {
			return fmt.Errorf("insert into alerts: %v", err)
		}

		return nil
	}); err != nil {
		log.Printf("transaction: %v\n", err)
		return alert, fmt.Errorf("transaction: %v\n", err)
	}

	return alert, nil
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
