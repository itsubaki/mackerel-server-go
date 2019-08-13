package database

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/itsubaki/mackerel-api/pkg/domain"
)

type HostRepository struct {
	SQLHandler
}

func NewHostRepository(handler SQLHandler) *HostRepository {
	if err := handler.Transact(func(tx Tx) error {
		if _, err := tx.Exec(
			`
			create table if not exists hosts (
				org_id            varchar(16)  not null,
				id                varchar(16)  not null primary key,
				name              varchar(128) not null,
				status            enum('working', 'standby', 'maintenance', 'poweroff') not null,
				memo              varchar(128) not null default '',
				display_name      varchar(128),
				custom_identifier varchar(128),
				created_at        bigint,
				retired_at        bigint,
				is_retired        boolean,
				roles             text,
				role_fullnames    text,
				interfaces        text,
				checks            text,
				meta              text
			)
			`,
		); err != nil {
			return fmt.Errorf("create table hosts: %v", err)
		}

		return nil
	}); err != nil {
		panic(fmt.Errorf("transaction: %v", err))
	}

	return &HostRepository{
		SQLHandler: handler,
	}
}

func (repo *HostRepository) ActiveList(orgID string) (*domain.Hosts, error) {
	hosts := make([]domain.Host, 0)

	if err := repo.Transact(func(tx Tx) error {
		rows, err := tx.Query("select * from hosts where org_id=? and is_retired=0", orgID)
		if err != nil {
			return fmt.Errorf("select * from hosts: %v", err)
		}
		defer rows.Close()

		for rows.Next() {
			var host domain.Host
			var roles, roleFullnames, interfaces, checks, meta string
			if err := rows.Scan(
				&host.OrgID,
				&host.ID,
				&host.Name,
				&host.Status,
				&host.Memo,
				&host.DisplayName,
				&host.CustomIdentifier,
				&host.CreatedAt,
				&host.RetiredAt,
				&host.IsRetired,
				&roles,
				&roleFullnames,
				&interfaces,
				&checks,
				&meta,
			); err != nil {
				return fmt.Errorf("scan: %v", err)
			}

			if err := json.Unmarshal([]byte(roles), &host.Roles); err != nil {
				return fmt.Errorf("unmarshal host.Roles: %v", err)
			}

			if err := json.Unmarshal([]byte(roleFullnames), &host.RoleFullNames); err != nil {
				return fmt.Errorf("unmarshal host.RoleFullNames: %v", err)
			}

			if err := json.Unmarshal([]byte(interfaces), &host.Interfaces); err != nil {
				return fmt.Errorf("unmarshal host.Interfaces: %v", err)
			}

			if err := json.Unmarshal([]byte(checks), &host.Checks); err != nil {
				return fmt.Errorf("unmarshal host.Checks: %v", err)
			}

			if err := json.Unmarshal([]byte(meta), &host.Meta); err != nil {
				return fmt.Errorf("unmarshal host.Meta: %v", err)
			}

			hosts = append(hosts, host)
		}

		return nil
	}); err != nil {
		return nil, fmt.Errorf("transaction: %v", err)
	}

	return &domain.Hosts{Hosts: hosts}, nil
}

// mysql> explain select * from hosts;
// +----+-------------+-------+------------+------+---------------+------+---------+------+------+----------+-------+
// | id | select_type | table | partitions | type | possible_keys | key  | key_len | ref  | rows | filtered | Extra |
// +----+-------------+-------+------------+------+---------------+------+---------+------+------+----------+-------+
// |  1 | SIMPLE      | hosts | NULL       | ALL  | NULL          | NULL | NULL    | NULL |    4 |   100.00 | NULL  |
// +----+-------------+-------+------------+------+---------------+------+---------+------+------+----------+-------+
// 1 row in set, 1 warning (0.01 sec)
func (repo *HostRepository) List(orgID string) (*domain.Hosts, error) {
	hosts := make([]domain.Host, 0)

	if err := repo.Transact(func(tx Tx) error {
		rows, err := tx.Query("select * from hosts where org_id=?", orgID)
		if err != nil {
			return fmt.Errorf("select * from hosts: %v", err)
		}
		defer rows.Close()

		for rows.Next() {
			var host domain.Host
			var roles, roleFullnames, interfaces, checks, meta string
			if err := rows.Scan(
				&host.OrgID,
				&host.ID,
				&host.Name,
				&host.Status,
				&host.Memo,
				&host.DisplayName,
				&host.CustomIdentifier,
				&host.CreatedAt,
				&host.RetiredAt,
				&host.IsRetired,
				&roles,
				&roleFullnames,
				&interfaces,
				&checks,
				&meta,
			); err != nil {
				return fmt.Errorf("scan: %v", err)
			}

			if err := json.Unmarshal([]byte(roles), &host.Roles); err != nil {
				return fmt.Errorf("unmarshal host.Roles: %v", err)
			}

			if err := json.Unmarshal([]byte(roleFullnames), &host.RoleFullNames); err != nil {
				return fmt.Errorf("unmarshal host.RoleFullNames: %v", err)
			}

			if err := json.Unmarshal([]byte(interfaces), &host.Interfaces); err != nil {
				return fmt.Errorf("unmarshal host.Interfaces: %v", err)
			}

			if err := json.Unmarshal([]byte(checks), &host.Checks); err != nil {
				return fmt.Errorf("unmarshal host.Checks: %v", err)
			}

			if err := json.Unmarshal([]byte(meta), &host.Meta); err != nil {
				return fmt.Errorf("unmarshal host.Meta: %v", err)
			}

			hosts = append(hosts, host)
		}

		return nil
	}); err != nil {
		return nil, fmt.Errorf("transaction: %v", err)
	}

	return &domain.Hosts{Hosts: hosts}, nil
}

// insert into hosts values(${name}, ${meta}, ${interface}, ${checks}, ${display_name}, ${custom_identifier}, ${created_at}, ${id}, ${status}, ${memo}, ${roles}, ${is_retired}, ${retired_at} )
func (repo *HostRepository) Save(orgID string, host *domain.Host) (*domain.HostID, error) {
	if err := repo.Transact(func(tx Tx) error {
		roles, err := json.Marshal(host.Roles)
		if err != nil {
			return fmt.Errorf("marshal host.Roles: %v", err)
		}

		roleFullnames, err := json.Marshal(host.RoleFullNames)
		if err != nil {
			return fmt.Errorf("marshal host.RoleFullNames: %v", err)
		}

		interfaces, err := json.Marshal(host.Interfaces)
		if err != nil {
			return fmt.Errorf("marshal host.Interfaces: %v", err)
		}

		checks, err := json.Marshal(host.Checks)
		if err != nil {
			return fmt.Errorf("marshal host.Checks: %v", err)
		}

		meta, err := json.Marshal(host.Meta)
		if err != nil {
			return fmt.Errorf("marshal host.Meta: %v", err)
		}

		if _, err := tx.Exec(
			`
			insert into hosts (
				org_id,
				id,
				name,
				status,
				memo,
				display_name,
				custom_identifier,
				created_at,
				retired_at,
				is_retired,
				roles,
				role_fullnames,
				interfaces,
				checks,
				meta
			)
			values (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
			on duplicate key update
				name = values(name),
				status = values(status),
				memo = values(memo),
				display_name = values(display_name),
				custom_identifier = values(custom_identifier),
				roles = values(roles),
				role_fullnames = values(role_fullnames),
				interfaces = values(interfaces),
				checks = values(checks),
				meta = values(meta)
			`,
			orgID,
			host.ID,
			host.Name,
			host.Status,
			host.Memo,
			host.DisplayName,
			host.CustomIdentifier,
			host.CreatedAt,
			host.RetiredAt,
			host.IsRetired,
			string(roles),
			string(roleFullnames),
			string(interfaces),
			string(checks),
			string(meta),
		); err != nil {
			return fmt.Errorf("insert into hosts: %v", err)

		}

		for svc, role := range host.Roles {
			if _, err := tx.Exec(
				`
				insert into services (
					org_id,
					name
				)
				select ?, ? where not exists (select 1 from services where org_id=? and name=?)
				`,
				orgID, svc,
				orgID, svc,
			); err != nil {
				return fmt.Errorf("insert into services: %v", err)
			}

			for i := range role {
				if _, err := tx.Exec(
					`
					insert into roles (
						org_id,
						service_name,
						name
					)
					select ?, ?, ? where not exists (select 1 from roles where org_id=? and service_name=? and name=?)
					`,
					orgID, svc, role[i],
					orgID, svc, role[i],
				); err != nil {
					return fmt.Errorf("insert into roles: %v", err)
				}

			}
		}

		return nil
	}); err != nil {
		return nil, fmt.Errorf("transaction: %v", err)
	}

	return &domain.HostID{ID: host.ID}, nil
}

// mysql> explain select * from hosts where id='de3d16e34dc';
// +----+-------------+-------+------------+-------+---------------+---------+---------+-------+------+----------+-------+
// | id | select_type | table | partitions | type  | possible_keys | key     | key_len | ref   | rows | filtered | Extra |
// +----+-------------+-------+------------+-------+---------------+---------+---------+-------+------+----------+-------+
// |  1 | SIMPLE      | hosts | NULL       | const | PRIMARY       | PRIMARY | 66      | const |    1 |   100.00 | NULL  |
// +----+-------------+-------+------------+-------+---------------+---------+---------+-------+------+----------+-------+
// 1 row in set, 1 warning (0.01 sec)
func (repo *HostRepository) Host(orgID, hostID string) (*domain.Host, error) {
	var host domain.Host

	if err := repo.Transact(func(tx Tx) error {
		row := tx.QueryRow("select * from hosts where org_id=? and id=?", orgID, hostID)
		var roles, roleFullnames, interfaces, checks, meta string
		if err := row.Scan(
			&host.OrgID,
			&host.ID,
			&host.Name,
			&host.Status,
			&host.Memo,
			&host.DisplayName,
			&host.CustomIdentifier,
			&host.CreatedAt,
			&host.RetiredAt,
			&host.IsRetired,
			&roles,
			&roleFullnames,
			&interfaces,
			&checks,
			&meta,
		); err != nil {
			return fmt.Errorf("select * from hosts: %v", err)
		}

		if err := json.Unmarshal([]byte(roles), &host.Roles); err != nil {
			return fmt.Errorf("unmarshal host.Roles: %v", err)
		}

		if err := json.Unmarshal([]byte(roleFullnames), &host.RoleFullNames); err != nil {
			return fmt.Errorf("unmarshal host.RoleFullNames: %v", err)
		}

		if err := json.Unmarshal([]byte(interfaces), &host.Interfaces); err != nil {
			return fmt.Errorf("unmarshal host.Interfaces: %v", err)
		}

		if err := json.Unmarshal([]byte(checks), &host.Checks); err != nil {
			return fmt.Errorf("unmarshal host.Checks: %v", err)
		}

		if err := json.Unmarshal([]byte(meta), &host.Meta); err != nil {
			return fmt.Errorf("unmarshal host.Meta: %v", err)
		}

		return nil
	}); err != nil {
		return nil, fmt.Errorf("transaction: %v", err)
	}

	return &host, nil
}

// mysql> explain select * from hosts where id='de3d16e34dc' limit 1;
// +----+-------------+-------+------------+-------+---------------+---------+---------+-------+------+----------+-------+
// | id | select_type | table | partitions | type  | possible_keys | key     | key_len | ref   | rows | filtered | Extra |
// +----+-------------+-------+------------+-------+---------------+---------+---------+-------+------+----------+-------+
// |  1 | SIMPLE      | hosts | NULL       | const | PRIMARY       | PRIMARY | 66      | const |    1 |   100.00 | NULL  |
// +----+-------------+-------+------------+-------+---------------+---------+---------+-------+------+----------+-------+
// 1 row in set, 1 warning (0.00 sec)
func (repo *HostRepository) Exists(orgID, hostID string) bool {
	rows, err := repo.Query("select 1 from hosts where org_id=? and id=? limit 1", orgID, hostID)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	if rows.Next() {
		return true
	}

	return false
}

// update hosts set status=${status} where id=${hostID}
func (repo *HostRepository) Status(orgID, hostID, status string) (*domain.Success, error) {
	if err := repo.Transact(func(tx Tx) error {
		if _, err := tx.Exec("update hosts set status=? where org_id=? and id=?", status, orgID, hostID); err != nil {
			return fmt.Errorf("update hosts: %v", err)
		}

		return nil
	}); err != nil {
		return &domain.Success{Success: false}, fmt.Errorf("transaction: %v", err)
	}

	return &domain.Success{Success: true}, nil
}

// update hosts set roles=${roles} where id=${hostID}
func (repo *HostRepository) SaveRoleFullNames(orgID, hostID string, names *domain.RoleFullNames) (*domain.Success, error) {
	roles := names.Roles()

	mroles, err := json.Marshal(roles)
	if err != nil {
		return &domain.Success{Success: false}, fmt.Errorf("marshal: %v", err)
	}

	roleFullnames, err := json.Marshal(names.Names)
	if err != nil {
		return &domain.Success{Success: false}, fmt.Errorf("marshal: %v", err)
	}

	if err := repo.Transact(func(tx Tx) error {
		if _, err := tx.Exec(
			"update hosts set role_fullnames=?, roles=? where org_id=? and id=?",
			string(roleFullnames),
			string(mroles),
			orgID,
			hostID,
		); err != nil {
			return fmt.Errorf("update hosts: %v", err)
		}

		for svc, role := range roles {
			if _, err := tx.Exec(
				`
				insert into services (
					org_id,
					name
				)
				select ?, ? where not exists (select 1 from services where org_id=? and name=?)
				`,
				orgID, svc,
				orgID, svc,
			); err != nil {
				return fmt.Errorf("insert into services: %v", err)
			}

			for i := range role {
				if _, err := tx.Exec(
					`
					insert into roles (
						org_id,
						service_name,
						name
					)
					select ?, ?, ? where not exists (select 1 from roles where org_id=? and service_name=? and name=?)
					`,
					orgID, svc, role[i],
					orgID, svc, role[i],
				); err != nil {
					return fmt.Errorf("insert into roles: %v", err)
				}

			}
		}

		return nil
	}); err != nil {
		return &domain.Success{Success: false}, fmt.Errorf("transaction: %v", err)
	}

	return &domain.Success{Success: true}, nil
}

// update hosts set is_retired=true, retired_at=time.Now().Unix() where id=${hostID}
func (repo *HostRepository) Retire(orgID, hostID string, retire *domain.HostRetire) (*domain.Success, error) {
	if err := repo.Transact(func(tx Tx) error {
		if _, err := tx.Exec("update hosts set is_retired=?, retired_at=? where org_id=? and id=?", true, time.Now().Unix(), orgID, hostID); err != nil {
			return fmt.Errorf("update hosts: %v", err)
		}

		return nil
	}); err != nil {
		return &domain.Success{Success: false}, fmt.Errorf("transaction: %v", err)
	}

	return &domain.Success{Success: true}, nil
}
