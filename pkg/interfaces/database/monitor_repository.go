package database

import (
	"fmt"

	"github.com/itsubaki/mackerel-api/pkg/domain"
)

type MonitorRepository struct {
	SQLHandler
}

func NewMonitorRepository(handler SQLHandler) *MonitorRepository {
	if err := handler.Transact(func(tx Tx) error {
		if _, err := tx.Exec(
			`
			create table if not exists monitors (
				org                               varchar(64) not null,
				id                                varchar(16) not null primary key,
				type                              enum('host', 'connectivity', 'service', 'external', 'expression'),
				name                              varchar(128) not null,
				Memo                              text,
				notification_interval	          int  not null default 1,
				is_mute                           bool not null default 1,
				duration                          int  not null default 1,
				metric                            varchar(128),
				operator                          enum('>', '<'),
				warning                           double,
				critical                          double,
				max_check_attempts                int,
				scopes                            text,
				exclude_scopes                    text,
				missing_duration_warning          int,
				missing_duration_critical         int,
				url                               text,
				method                            enum('GET', 'PUT', 'POST', 'DELETE'),
				service                           text,
				response_time_warning             int,
				response_time_critical            int,
				response_time_duration            int,
				contains_string                   text,
				certification_expiration_warning  int,
				certification_expiration_critical int,
				skip_certificate_verification     bool,
				headers                           text,
				request_body                      text,
				expression                        text
			)
			`,
		); err != nil {
			return fmt.Errorf("create table monitors: %v", err)
		}

		return nil
	}); err != nil {
		panic(fmt.Errorf("transaction: %v", err))
	}

	return &MonitorRepository{
		SQLHandler: handler,
	}
}

func (repo *MonitorRepository) List(org string) (*domain.Monitors, error) {
	monitors := make([]interface{}, 0)
	if err := repo.Transact(func(tx Tx) error {
		rows, err := tx.Query("select * from monitors where org=?", org)
		if err != nil {
			return fmt.Errorf("select * from monitors: %v", err)
		}
		defer rows.Close()

		for rows.Next() {
			var mon domain.Monitoring
			var org string
			if err := rows.Scan(
				&org,
				&mon.ID,
			); err != nil {
				return fmt.Errorf("scan: %v", err)
			}

			monitors = append(monitors, mon)
		}

		return nil
	}); err != nil {
		return nil, fmt.Errorf("transaction: %v", err)
	}

	return &domain.Monitors{Monitors: monitors}, nil
}

func (repo *MonitorRepository) Save(org string, monitor *domain.Monitoring) (interface{}, error) {
	return nil, nil
}

func (repo *MonitorRepository) Update(org string, monitor *domain.Monitoring) (interface{}, error) {
	return nil, nil
}

func (repo *MonitorRepository) Monitor(org, monitorID string) (interface{}, error) {
	return nil, nil
}

func (repo *MonitorRepository) Delete(org, monitorID string) (interface{}, error) {
	return nil, nil
}
