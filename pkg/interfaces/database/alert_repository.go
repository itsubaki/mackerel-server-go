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
				id varchar(16) not null primary key,
				status varchar(16) not null,
				monitor_id varchar(16) not null,
				type varchar(16) not null,
				host_id varchar(16),
				value double not null,
				message text,
				reason text,
				opened_at bigint,
				closed_at bigint,
				seq bigint auto_increment, index(seq)
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

func (repo *AlertRepository) Exists(alertID string) bool {
	rows, err := repo.Query("select * from alerts where id=? limit 1", alertID)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	if rows.Next() {
		return true
	}

	return false
}

func (repo *AlertRepository) List(withClosed bool, nextID string, limit int) (*domain.Alerts, error) {
	// TODO withClosed, nextID
	rows, err := repo.Query(
		`
			select 
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
			from alerts
			limit ?
			`,
		limit,
	)
	if err != nil {
		return nil, fmt.Errorf("select * from alerts: %v", err)
	}
	defer rows.Close()

	var alerts []domain.Alert
	for rows.Next() {
		var alert domain.Alert
		if err := rows.Scan(
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

	return &domain.Alerts{Alerts: alerts}, nil
}

func (repo *AlertRepository) Close(alertID, reason string) (*domain.Alert, error) {
	var alert domain.Alert
	if err := repo.Transact(func(tx Tx) error {
		if _, err := tx.Exec(
			`
			update alerts set reason=?, closed_at=? where id=?
			`,
			reason,
			time.Now().Unix(),
			alertID,
		); err != nil {
			return fmt.Errorf("update alerts: %v", err)
		}

		row := tx.QueryRow(
			`
			select 
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
			from alerts
				where id=?
			for update`,
			alertID,
		)
		if err := row.Scan(
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
