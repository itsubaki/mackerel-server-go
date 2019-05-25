package database

import (
	"fmt"
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
				org        varchar(64) not null,
				id         varchar(16) not null primary key,
				status     enum('OK', 'CRITICAL', 'WARNING', 'UNKNOWN') not null,
				monitor_id varchar(16) not null,
				type       enum('connectivity', 'host', 'service', 'external', 'check', 'expression') not null,
				host_id    varchar(16),
				value      double not null,
				message    text,
				reason     text,
				opened_at  bigint,
				closed_at  bigint
			)
			`,
		); err != nil {
			return fmt.Errorf("create table alerts: %v", err)
		}

		return nil
	}); err != nil {
		panic(fmt.Errorf("transaction: %v", err))
	}

	return &AlertRepository{
		SQLHandler: handler,
	}
}

func (repo *AlertRepository) Exists(org, alertID string) bool {
	rows, err := repo.Query("select * from alerts where org=? and id=? limit 1", org, alertID)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	if rows.Next() {
		return true
	}

	return false
}

func (repo *AlertRepository) List(org string, withClosed bool, nextID string, limit int) (*domain.Alerts, error) {
	status := "UNKNOWN"
	if withClosed {
		status = "OK"
	}

	rows, err := repo.Query(
		`
		select * from alerts
		where org=? and status in ('CRITICAL', 'WARNING', 'UNKNOWN', ?)
		order by opened_at limit ?
		`,
		org,
		status,
		limit+1,
	)
	if err != nil {
		return nil, fmt.Errorf("select * from alerts: %v", err)
	}
	defer rows.Close()

	var alerts []domain.Alert
	for rows.Next() {
		var alert domain.Alert
		var org string
		if err := rows.Scan(
			&org,
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

func (repo *AlertRepository) Close(org, alertID, reason string) (*domain.Alert, error) {
	var alert domain.Alert
	if err := repo.Transact(func(tx Tx) error {
		if _, err := tx.Exec(
			`
			update alerts set reason=?, closed_at=? where org=? and id=?
			`,
			reason,
			time.Now().Unix(),
			org,
			alertID,
		); err != nil {
			return fmt.Errorf("update alerts: %v", err)
		}

		row := tx.QueryRow("select * from alerts where org=? and id=?", alertID)
		var org string
		if err := row.Scan(
			&org,
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

		return nil
	}); err != nil {
		return nil, fmt.Errorf("transaction: %v", err)
	}

	return &alert, nil
}
