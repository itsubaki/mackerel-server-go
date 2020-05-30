package database

import (
	"encoding/json"
	"fmt"

	"github.com/itsubaki/mackerel-server-go/pkg/domain"
)

type DowntimeRepository struct {
	SQLHandler
}

func NewDowntimeRepository(handler SQLHandler) *DowntimeRepository {
	if err := handler.Transact(func(tx Tx) error {
		if _, err := tx.Exec(
			`
			create table if not exists downtimes (
				org_id                 varchar(16)  not null,
				id                     varchar(128) not null primary key,
				name                   varchar(128) not null,
				memo				   text,
				start                  bigint,
				duration               bigint,
				recurrence             text,
				service_scopes         text,
				service_exclude_scopes text,
				role_scopes            text,
				role_exclude_scopes    text,
				monitor_scopes         text,
				monitor_exclude_scopes text
			)
			`,
		); err != nil {
			return fmt.Errorf("create table downtimes: %v", err)
		}

		return nil
	}); err != nil {
		panic(fmt.Errorf("transaction: %v", err))
	}

	return &DowntimeRepository{
		SQLHandler: handler,
	}
}

func (repo *DowntimeRepository) List(orgID string) (*domain.Downtimes, error) {
	downtimes := make([]domain.Downtime, 0)
	if err := repo.Transact(func(tx Tx) error {
		rows, err := tx.Query("select * from downtimes where org_id=?", orgID)
		if err != nil {
			return fmt.Errorf("select * from downtimes: %v", err)
		}
		defer rows.Close()

		for rows.Next() {
			var downtime domain.Downtime
			var recurrence, serviceScopes, serviceExcludeScopes, roleScopes, roleExcludeScopes, monitorScopes, monitorExcludeScopes string
			if err := rows.Scan(
				&downtime.OrgID,
				&downtime.ID,
				&downtime.Name,
				&downtime.Memo,
				&downtime.Start,
				&downtime.Duration,
				&recurrence,
				&serviceScopes,
				&serviceExcludeScopes,
				&roleScopes,
				&roleExcludeScopes,
				&monitorScopes,
				&monitorExcludeScopes,
			); err != nil {
				return fmt.Errorf("select * from downtimes: %v", err)
			}

			if err := json.Unmarshal([]byte(recurrence), &downtime.Recurrence); err != nil {
				return fmt.Errorf("unmarshal downitme.Recurrence: %v", err)
			}

			if err := json.Unmarshal([]byte(serviceScopes), &downtime.ServiceScopes); err != nil {
				return fmt.Errorf("unmarshal downitme.ServiceScopes: %v", err)
			}

			if err := json.Unmarshal([]byte(serviceExcludeScopes), &downtime.ServiceExcludeScopes); err != nil {
				return fmt.Errorf("unmarshal downitme.ServiceExcludeScopes: %v", err)
			}

			if err := json.Unmarshal([]byte(roleScopes), &downtime.RoleScopes); err != nil {
				return fmt.Errorf("unmarshal downitme.RoleScopes: %v", err)
			}

			if err := json.Unmarshal([]byte(roleExcludeScopes), &downtime.RoleExcludeScopes); err != nil {
				return fmt.Errorf("unmarshal downitme.RoleExcludeScopes: %v", err)
			}

			if err := json.Unmarshal([]byte(monitorScopes), &downtime.MonitorScopes); err != nil {
				return fmt.Errorf("unmarshal downitme.MonitorScopes: %v", err)
			}

			if err := json.Unmarshal([]byte(monitorExcludeScopes), &downtime.MonitorExcludeScopes); err != nil {
				return fmt.Errorf("unmarshal downitme.MonitorExcludeScopes: %v", err)
			}

			downtimes = append(downtimes, downtime)
		}

		return nil
	}); err != nil {
		return nil, fmt.Errorf("transaction: %v", err)
	}

	return &domain.Downtimes{Downtimes: downtimes}, nil
}

func (repo *DowntimeRepository) Save(orgID string, downtime *domain.Downtime) (*domain.Downtime, error) {
	if err := repo.Transact(func(tx Tx) error {
		recurrence, err := json.Marshal(downtime.Recurrence)
		if err != nil {
			return fmt.Errorf("marshal downtime.Recurrenc: %v", err)
		}

		serviceScopes, err := json.Marshal(downtime.ServiceScopes)
		if err != nil {
			return fmt.Errorf("marshal downtime.ServiceScopes: %v", err)
		}

		serviceExcludeScopes, err := json.Marshal(downtime.ServiceExcludeScopes)
		if err != nil {
			return fmt.Errorf("marshal downtime.ServiceExcludeScopes: %v", err)
		}

		roleScopes, err := json.Marshal(downtime.RoleScopes)
		if err != nil {
			return fmt.Errorf("marshal downtime.RoleScopes: %v", err)
		}

		roleExcludeScopes, err := json.Marshal(downtime.RoleExcludeScopes)
		if err != nil {
			return fmt.Errorf("marshal downtime.RoleExcludeScopes: %v", err)
		}

		monitorScopes, err := json.Marshal(downtime.MonitorScopes)
		if err != nil {
			return fmt.Errorf("marshal downtime.MonitorScopes: %v", err)
		}

		monitorExcludeScopes, err := json.Marshal(downtime.MonitorExcludeScopes)
		if err != nil {
			return fmt.Errorf("marshal downtime.MonitorExcludeScopes: %v", err)
		}

		if _, err := tx.Exec(
			`
			insert into downtimes (
				org_id,
				id,
				name,
				memo,
				start,
				duration,
				recurrence,
				service_scopes,
				service_exclude_scopes,
				role_scopes,
				role_exclude_scopes,
				monitor_scopes,
				monitor_exclude_scopes
			)
			values (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
			`,
			orgID,
			downtime.ID,
			downtime.Name,
			downtime.Memo,
			downtime.Start,
			downtime.Duration,
			recurrence,
			serviceScopes,
			serviceExcludeScopes,
			roleScopes,
			roleExcludeScopes,
			monitorScopes,
			monitorExcludeScopes,
		); err != nil {
			return fmt.Errorf("insert into downtimes: %v", err)

		}

		return nil
	}); err != nil {
		return nil, fmt.Errorf("transaction: %v", err)
	}

	return downtime, nil
}

func (repo *DowntimeRepository) Update(orgID string, downtime *domain.Downtime) (*domain.Downtime, error) {
	if err := repo.Transact(func(tx Tx) error {
		recurrence, err := json.Marshal(downtime.Recurrence)
		if err != nil {
			return fmt.Errorf("marshal downtime.Recurrence: %v", err)
		}

		serviceScopes, err := json.Marshal(downtime.ServiceScopes)
		if err != nil {
			return fmt.Errorf("marshal downtime.ServiceScopes: %v", err)
		}

		serviceExcludeScopes, err := json.Marshal(downtime.ServiceExcludeScopes)
		if err != nil {
			return fmt.Errorf("marshal downtime.ServiceExcludeScopes: %v", err)
		}

		roleScopes, err := json.Marshal(downtime.RoleScopes)
		if err != nil {
			return fmt.Errorf("marshal downtime.RoleScopes: %v", err)
		}

		roleExcludeScopes, err := json.Marshal(downtime.RoleExcludeScopes)
		if err != nil {
			return fmt.Errorf("marshal downtime.RoleExcludeScopes: %v", err)
		}

		monitorScopes, err := json.Marshal(downtime.MonitorScopes)
		if err != nil {
			return fmt.Errorf("marshal downtime.MonitorScopes: %v", err)
		}

		monitorExcludeScopes, err := json.Marshal(downtime.MonitorExcludeScopes)
		if err != nil {
			return fmt.Errorf("marshal downtime.MonitorExcludeScopes: %v", err)
		}

		if _, err := tx.Exec(
			`
			update downtimes set
				name=?,
				memo=?,
				start=?,
				duration=?,
				recurrence=?,
				service_scopes=?,
				service_exclude_scopes=?,
				role_scopes=?,
				role_exclude_scopes=?,
				monitor_scopes=?,
				monitor_exclude_scopes=?
			where org_id=? and id=?
			`,
			downtime.Name,
			downtime.Memo,
			downtime.Start,
			downtime.Duration,
			recurrence,
			serviceScopes,
			serviceExcludeScopes,
			roleScopes,
			roleExcludeScopes,
			monitorScopes,
			monitorExcludeScopes,
			orgID,
			downtime.ID,
		); err != nil {
			return fmt.Errorf("update downtimes: %v", err)

		}

		return nil
	}); err != nil {
		return nil, fmt.Errorf("transaction: %v", err)
	}

	return downtime, nil
}

func (repo *DowntimeRepository) Downtime(orgID, downtimeID string) (*domain.Downtime, error) {
	var downtime domain.Downtime

	if err := repo.Transact(func(tx Tx) error {
		row := tx.QueryRow("select * from downtimes where org_id=? and id=?", orgID, downtimeID)
		var recurrence, serviceScopes, serviceExcludeScopes, roleScopes, roleExcludeScopes, monitorScopes, monitorExcludeScopes string
		if err := row.Scan(
			&downtime.OrgID,
			&downtime.ID,
			&downtime.Name,
			&downtime.Memo,
			&downtime.Start,
			&downtime.Duration,
			&recurrence,
			&serviceScopes,
			&serviceExcludeScopes,
			&roleScopes,
			&roleExcludeScopes,
			&monitorScopes,
			&monitorExcludeScopes,
		); err != nil {
			return fmt.Errorf("select * from downtimes: %v", err)
		}

		if err := json.Unmarshal([]byte(recurrence), &downtime.Recurrence); err != nil {
			return fmt.Errorf("unmarshal downitme.Recurrence: %v", err)
		}
		if len(downtime.Recurrence.Type) < 1 {
			downtime.Recurrence = nil
		}

		if err := json.Unmarshal([]byte(serviceScopes), &downtime.ServiceScopes); err != nil {
			return fmt.Errorf("unmarshal downitme.ServiceScopes: %v", err)
		}

		if err := json.Unmarshal([]byte(serviceExcludeScopes), &downtime.ServiceExcludeScopes); err != nil {
			return fmt.Errorf("unmarshal downitme.ServiceExcludeScopes: %v", err)
		}

		if err := json.Unmarshal([]byte(roleScopes), &downtime.RoleScopes); err != nil {
			return fmt.Errorf("unmarshal downitme.RoleScopes: %v", err)
		}

		if err := json.Unmarshal([]byte(roleExcludeScopes), &downtime.RoleExcludeScopes); err != nil {
			return fmt.Errorf("unmarshal downitme.RoleExcludeScopes: %v", err)
		}

		if err := json.Unmarshal([]byte(monitorScopes), &downtime.MonitorScopes); err != nil {
			return fmt.Errorf("unmarshal downitme.MonitorScopes: %v", err)
		}

		if err := json.Unmarshal([]byte(monitorExcludeScopes), &downtime.MonitorExcludeScopes); err != nil {
			return fmt.Errorf("unmarshal downitme.MonitorExcludeScopes: %v", err)
		}

		return nil
	}); err != nil {
		return nil, fmt.Errorf("transaction: %v", err)
	}

	return &downtime, nil
}

func (repo *DowntimeRepository) Delete(orgID, downtimeID string) (*domain.Downtime, error) {
	var downtime domain.Downtime

	if err := repo.Transact(func(tx Tx) error {
		row := tx.QueryRow("select * from downtimes where org_id=? and id=?", orgID, downtimeID)
		var recurrence, serviceScopes, serviceExcludeScopes, roleScopes, roleExcludeScopes, monitorScopes, monitorExcludeScopes string
		if err := row.Scan(
			&downtime.OrgID,
			&downtime.ID,
			&downtime.Name,
			&downtime.Memo,
			&downtime.Start,
			&downtime.Duration,
			&recurrence,
			&serviceScopes,
			&serviceExcludeScopes,
			&roleScopes,
			&roleExcludeScopes,
			&monitorScopes,
			&monitorExcludeScopes,
		); err != nil {
			return fmt.Errorf("select * from downtimes: %v", err)
		}

		if err := json.Unmarshal([]byte(recurrence), &downtime.Recurrence); err != nil {
			return fmt.Errorf("unmarshal downitme.Recurrence: %v", err)
		}

		if err := json.Unmarshal([]byte(serviceScopes), &downtime.ServiceScopes); err != nil {
			return fmt.Errorf("unmarshal downitme.ServiceScopes: %v", err)
		}

		if err := json.Unmarshal([]byte(serviceExcludeScopes), &downtime.ServiceExcludeScopes); err != nil {
			return fmt.Errorf("unmarshal downitme.ServiceExcludeScopes: %v", err)
		}

		if err := json.Unmarshal([]byte(roleScopes), &downtime.RoleScopes); err != nil {
			return fmt.Errorf("unmarshal downitme.RoleScopes: %v", err)
		}

		if err := json.Unmarshal([]byte(roleExcludeScopes), &downtime.RoleExcludeScopes); err != nil {
			return fmt.Errorf("unmarshal downitme.RoleExcludeScopes: %v", err)
		}

		if err := json.Unmarshal([]byte(monitorScopes), &downtime.MonitorScopes); err != nil {
			return fmt.Errorf("unmarshal downitme.MonitorScopes: %v", err)
		}

		if err := json.Unmarshal([]byte(monitorExcludeScopes), &downtime.MonitorExcludeScopes); err != nil {
			return fmt.Errorf("unmarshal downitme.MonitorExcludeScopes: %v", err)
		}

		if _, err := tx.Exec("delete from downtimes where org_id=? and id=?", orgID, downtimeID); err != nil {
			return fmt.Errorf("delete from downtimes: %v", err)
		}

		return nil
	}); err != nil {
		return nil, fmt.Errorf("transaction: %v", err)
	}

	return &downtime, nil
}
