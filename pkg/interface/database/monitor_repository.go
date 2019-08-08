package database

import (
	"encoding/json"
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
				org_id                            varchar(16) not null,
				id                                varchar(16) not null primary key,
				type                              enum('host', 'connectivity', 'service', 'external', 'expression'),
				name                              varchar(128) not null,
				memo                              varchar(128) not null default '',
				notification_interval	          int  not null default 1,
				is_mute                           bool not null default 1,
				duration                          int,
				metric                            varchar(128),
				operator                          enum('>', '<') not null default '<',
				warning                           double,
				critical                          double,
				max_check_attempts                int,
				scopes                            text,
				exclude_scopes                    text,
				missing_duration_warning          int,
				missing_duration_critical         int,
				url                               text,
				method                            enum('GET', 'PUT', 'POST', 'DELETE') not null default 'GET',
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

func (repo *MonitorRepository) List(orgID string) (*domain.Monitors, error) {
	monitors := make([]interface{}, 0)
	if err := repo.Transact(func(tx Tx) error {
		rows, err := tx.Query("select * from monitors where org_id=?", orgID)
		if err != nil {
			return fmt.Errorf("select * from monitors: %v", err)
		}
		defer rows.Close()

		for rows.Next() {
			var monitor domain.Monitoring
			if err := rows.Scan(
				&monitor.OrgID,
				&monitor.ID,
				&monitor.Type,
				&monitor.Name,
				&monitor.Memo,
				&monitor.NotificationInterval,
				&monitor.IsMute,
				&monitor.Duration,
				&monitor.Metric,
				&monitor.Operator,
				&monitor.Warning,
				&monitor.Critical,
				&monitor.MaxCheckAttempts,
				&monitor.Scopes,
				&monitor.ExcludeScopes,
				&monitor.MissingDurationWarning,
				&monitor.MissingDurationCritical,
				&monitor.URL,
				&monitor.Method,
				&monitor.Service,
				&monitor.ResponseTimeWarning,
				&monitor.ResponseTimeCritical,
				&monitor.ResponseTimeDuration,
				&monitor.ContainsString,
				&monitor.CertificationExpirationWarning,
				&monitor.CertificationExpirationCritical,
				&monitor.SkipCertificateVerification,
				&monitor.Headers,
				&monitor.RequestBody,
				&monitor.Expression,
			); err != nil {
				return fmt.Errorf("scan: %v", err)
			}

			monitors = append(monitors, monitor.Cast())
		}

		return nil
	}); err != nil {
		return nil, fmt.Errorf("transaction: %v", err)
	}

	return &domain.Monitors{Monitors: monitors}, nil
}

func (repo *MonitorRepository) Save(orgID string, monitor *domain.Monitoring) (interface{}, error) {
	if err := repo.Transact(func(tx Tx) error {
		scopes, err := json.Marshal(monitor.Scopes)
		if err != nil {
			return fmt.Errorf("marshal monitor.Scopes: %v", err)
		}

		exclude, err := json.Marshal(monitor.ExcludeScopes)
		if err != nil {
			return fmt.Errorf("marshal monitor.ExcludeScopes: %v", err)
		}

		service, err := json.Marshal(monitor.Service)
		if err != nil {
			return fmt.Errorf("marshal monitor.Service: %v", err)
		}

		headers, err := json.Marshal(monitor.Headers)
		if err != nil {
			return fmt.Errorf("marshal monitor.Headers: %v", err)
		}

		body, err := json.Marshal(monitor.RequestBody)
		if err != nil {
			return fmt.Errorf("marshal monitor.RequestBody: %v", err)
		}

		if _, err := tx.Exec(
			`
			insert into monitors (
				org_id,
				id,
				type,
				name,
				memo,
				notification_interval,
				is_mute,
				duration,
				metric,
				operator,
				warning,
				critical,
				max_check_attempts,
				scopes,
				exclude_scopes,
				missing_duration_warning,
				missing_duration_critical,
				url,
				method,
				service,
				response_time_warning,
				response_time_critical,
				response_time_duration,
				contains_string,
				certification_expiration_warning,
				certification_expiration_critical,
				skip_certificate_verification,
				headers,
				request_body,
				expression
			)
			values (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
			`,
			orgID,
			monitor.ID,
			monitor.Type,
			monitor.Name,
			monitor.Memo,
			monitor.NotificationInterval,
			monitor.IsMute,
			monitor.Duration,
			monitor.Metric,
			monitor.Operator,
			monitor.Warning,
			monitor.Critical,
			monitor.MaxCheckAttempts,
			scopes,
			exclude,
			monitor.MissingDurationWarning,
			monitor.MissingDurationCritical,
			monitor.URL,
			monitor.Method,
			service,
			monitor.ResponseTimeWarning,
			monitor.ResponseTimeCritical,
			monitor.ResponseTimeDuration,
			monitor.ContainsString,
			monitor.CertificationExpirationWarning,
			monitor.CertificationExpirationCritical,
			monitor.SkipCertificateVerification,
			headers,
			body,
			monitor.Expression,
		); err != nil {
			return fmt.Errorf("insert into monitors: %v", err)

		}

		return nil
	}); err != nil {
		return nil, fmt.Errorf("transaction: %v", err)
	}

	return monitor, nil
}

func (repo *MonitorRepository) Update(orgID string, monitor *domain.Monitoring) (interface{}, error) {
	if err := repo.Transact(func(tx Tx) error {
		scopes, err := json.Marshal(monitor.Scopes)
		if err != nil {
			return fmt.Errorf("marshal monitor.Scopes: %v", err)
		}

		exclude, err := json.Marshal(monitor.ExcludeScopes)
		if err != nil {
			return fmt.Errorf("marshal monitor.ExcludeScopes: %v", err)
		}

		service, err := json.Marshal(monitor.Service)
		if err != nil {
			return fmt.Errorf("marshal monitor.Service: %v", err)
		}

		headers, err := json.Marshal(monitor.Headers)
		if err != nil {
			return fmt.Errorf("marshal monitor.Headers: %v", err)
		}

		body, err := json.Marshal(monitor.RequestBody)
		if err != nil {
			return fmt.Errorf("marshal monitor.RequestBody: %v", err)
		}

		if _, err := tx.Exec(
			`
			update monitors set
				type=?,
				name=?,
				memo=?,
				notification_interval=?,
				is_mute=?,
				duration=?,
				metric=?,
				operator=?,
				warning=?,
				critical=?,
				max_check_attempts=?,
				scopes=?,
				exclude_scopes=?,
				missing_duration_warning=?,
				missing_duration_critical=?,
				url=?,
				method=?,
				service=?,
				response_time_warning=?,
				response_time_critical=?,
				response_time_duration=?,
				contains_string=?,
				certification_expiration_warning=?,
				certification_expiration_critical=?,
				skip_certificate_verification=?,
				headers=?,
				request_body=?,
				expression=?
			where org_id=? and id=?
			`,
			monitor.Type,
			monitor.Name,
			monitor.Memo,
			monitor.NotificationInterval,
			monitor.IsMute,
			monitor.Duration,
			monitor.Metric,
			monitor.Operator,
			monitor.Warning,
			monitor.Critical,
			monitor.MaxCheckAttempts,
			scopes,
			exclude,
			monitor.MissingDurationWarning,
			monitor.MissingDurationCritical,
			monitor.URL,
			monitor.Method,
			service,
			monitor.ResponseTimeWarning,
			monitor.ResponseTimeCritical,
			monitor.ResponseTimeDuration,
			monitor.ContainsString,
			monitor.CertificationExpirationWarning,
			monitor.CertificationExpirationCritical,
			monitor.SkipCertificateVerification,
			headers,
			body,
			monitor.Expression,
			orgID,
			monitor.ID,
		); err != nil {
			return fmt.Errorf("insert into monitors: %v", err)

		}

		return nil
	}); err != nil {
		return nil, fmt.Errorf("transaction: %v", err)
	}

	return monitor.Cast(), nil
}

func (repo *MonitorRepository) Monitor(orgID, monitorID string) (interface{}, error) {
	var monitor domain.Monitoring
	if err := repo.Transact(func(tx Tx) error {
		row := tx.QueryRow("select * from monitors where org_id=? and id=?", orgID, monitorID)
		var scopes, exclude, service, headers, body string
		if err := row.Scan(
			&monitor.OrgID,
			&monitor.ID,
			&monitor.Type,
			&monitor.Name,
			&monitor.Memo,
			&monitor.NotificationInterval,
			&monitor.IsMute,
			&monitor.Duration,
			&monitor.Metric,
			&monitor.Operator,
			&monitor.Warning,
			&monitor.Critical,
			&monitor.MaxCheckAttempts,
			&scopes,
			&exclude,
			&monitor.MissingDurationWarning,
			&monitor.MissingDurationCritical,
			&monitor.URL,
			&monitor.Method,
			&service,
			&monitor.ResponseTimeWarning,
			&monitor.ResponseTimeCritical,
			&monitor.ResponseTimeDuration,
			&monitor.ContainsString,
			&monitor.CertificationExpirationWarning,
			&monitor.CertificationExpirationCritical,
			&monitor.SkipCertificateVerification,
			&headers,
			&body,
			&monitor.Expression,
		); err != nil {
			return fmt.Errorf("scan: %v", err)
		}

		if err := json.Unmarshal([]byte(scopes), &monitor.Scopes); err != nil {
			return fmt.Errorf("unmarshal monitor.Scopes: %v", err)
		}

		if err := json.Unmarshal([]byte(exclude), &monitor.ExcludeScopes); err != nil {
			return fmt.Errorf("unmarshal monitor.ExcludeScopes: %v", err)
		}

		if err := json.Unmarshal([]byte(service), &monitor.Service); err != nil {
			return fmt.Errorf("unmarshal monitor.Service: %v", err)
		}

		if err := json.Unmarshal([]byte(headers), &monitor.Headers); err != nil {
			return fmt.Errorf("unmarshal monitor.Headers: %v", err)
		}

		if err := json.Unmarshal([]byte(body), &monitor.RequestBody); err != nil {
			return fmt.Errorf("unmarshal monitor.RequestBody: %v", err)
		}

		return nil
	}); err != nil {
		return nil, fmt.Errorf("transaction: %v", err)
	}

	return monitor.Cast(), nil
}

func (repo *MonitorRepository) Delete(orgID, monitorID string) (interface{}, error) {
	var monitor domain.Monitoring
	if err := repo.Transact(func(tx Tx) error {
		row := tx.QueryRow("select * from monitors where org_id=? and id=?", orgID, monitorID)
		var scopes, exclude, service, headers, body string
		if err := row.Scan(
			&monitor.OrgID,
			&monitor.ID,
			&monitor.Type,
			&monitor.Name,
			&monitor.Memo,
			&monitor.NotificationInterval,
			&monitor.IsMute,
			&monitor.Duration,
			&monitor.Metric,
			&monitor.Operator,
			&monitor.Warning,
			&monitor.Critical,
			&monitor.MaxCheckAttempts,
			&scopes,
			&exclude,
			&monitor.MissingDurationWarning,
			&monitor.MissingDurationCritical,
			&monitor.URL,
			&monitor.Method,
			&service,
			&monitor.ResponseTimeWarning,
			&monitor.ResponseTimeCritical,
			&monitor.ResponseTimeDuration,
			&monitor.ContainsString,
			&monitor.CertificationExpirationWarning,
			&monitor.CertificationExpirationCritical,
			&monitor.SkipCertificateVerification,
			&headers,
			&body,
			&monitor.Expression,
		); err != nil {
			return fmt.Errorf("scan: %v", err)
		}

		if err := json.Unmarshal([]byte(scopes), &monitor.Scopes); err != nil {
			return fmt.Errorf("unmarshal monitor.Scopes: %v", err)
		}

		if err := json.Unmarshal([]byte(exclude), &monitor.ExcludeScopes); err != nil {
			return fmt.Errorf("unmarshal monitor.ExcludeScopes: %v", err)
		}

		if err := json.Unmarshal([]byte(service), &monitor.Service); err != nil {
			return fmt.Errorf("unmarshal monitor.Service: %v", err)
		}

		if err := json.Unmarshal([]byte(headers), &monitor.Headers); err != nil {
			return fmt.Errorf("unmarshal monitor.Headers: %v", err)
		}

		if err := json.Unmarshal([]byte(body), &monitor.RequestBody); err != nil {
			return fmt.Errorf("unmarshal monitor.RequestBody: %v", err)
		}

		if _, err := tx.Exec("delete from monitors where org_id=? and id=?", orgID, monitorID); err != nil {
			return fmt.Errorf("delete from monitors: %v", err)
		}

		return nil
	}); err != nil {
		return nil, fmt.Errorf("transaction: %v", err)
	}

	return monitor.Cast(), nil
}
