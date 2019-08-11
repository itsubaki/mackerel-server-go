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

		if _, err := tx.Exec(
			`
			create table if not exists host_meta (
				org_id    varchar(16)  not null,
				host_id   varchar(16)  not null,
				namespace varchar(128) not null,
				meta      text,
				primary key(host_id, namespace)
			)
			`,
		); err != nil {
			return fmt.Errorf("create table host_meta: %v", err)
		}

		if _, err := tx.Exec(
			`
			create table if not exists host_metric_values (
				org_id  varchar(16)  not null,
				host_id varchar(16)  not null,
				name    varchar(128) not null,
				time    bigint not null,
				value   double not null,
				primary key(host_id, name, time desc)
			)
			`,
		); err != nil {
			return fmt.Errorf("create table host_metric_values: %v", err)
		}

		if _, err := tx.Exec(
			`
			create table if not exists host_metric_values_latest (
				org_id  varchar(16)  not null,
				host_id varchar(16)  not null,
				name    varchar(128) not null,
				value   double       not null,
				primary key(host_id, name)
			)
			`,
		); err != nil {
			return fmt.Errorf("create table host_metric_values_latest: %v", err)
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

// mysql> explain select * from host_metric_values where org_id='default' and host_id='bc754968495' and name='loadavg5' limit 1;
// +----+-------------+--------------------+------------+------+---------------+---------+---------+-------------+------+----------+-------------+
// | id | select_type | table              | partitions | type | possible_keys | key     | key_len | ref         | rows | filtered | Extra       |
// +----+-------------+--------------------+------------+------+---------------+---------+---------+-------------+------+----------+-------------+
// |  1 | SIMPLE      | host_metric_values | NULL       | ref  | PRIMARY       | PRIMARY | 580     | const,const |   14 |    10.00 | Using where |
// +----+-------------+--------------------+------------+------+---------------+---------+---------+-------------+------+----------+-------------+
// 1 row in set, 1 warning (0.00 sec)
func (repo *HostRepository) ExistsMetric(orgID, hostID, name string) bool {
	rows, err := repo.Query("select 1 from host_metric_values where org_id=? and host_id=? and name=? limit 1", orgID, hostID, name)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	if rows.Next() {
		return true
	}

	return false
}

// mysql> explain select distinct name from host_metric_values where org_id='default' and host_id='bc754968495';
// +----+-------------+--------------------+------------+------+---------------+---------+---------+-------+------+----------+-------------+
// | id | select_type | table              | partitions | type | possible_keys | key     | key_len | ref   | rows | filtered | Extra       |
// +----+-------------+--------------------+------------+------+---------------+---------+---------+-------+------+----------+-------------+
// |  1 | SIMPLE      | host_metric_values | NULL       | ref  | PRIMARY       | PRIMARY | 66      | const |  570 |    10.00 | Using where |
// +----+-------------+--------------------+------------+------+---------------+---------+---------+-------+------+----------+-------------+
// 1 row in set, 1 warning (0.00 sec)
func (repo *HostRepository) MetricNames(orgID, hostID string) (*domain.MetricNames, error) {
	names := make([]string, 0)
	if err := repo.Transact(func(tx Tx) error {
		rows, err := tx.Query("select distinct name from host_metric_values where org_id=? and host_id=?", orgID, hostID)
		if err != nil {
			return fmt.Errorf("select distinct name from host_metric_values: %v", err)
		}
		defer rows.Close()

		for rows.Next() {
			var name string
			if err := rows.Scan(&name); err != nil {
				return fmt.Errorf("scan: %v", err)
			}
			names = append(names, name)
		}

		return nil
	}); err != nil {
		return nil, fmt.Errorf("transaction: %v", err)
	}

	return &domain.MetricNames{Names: names}, nil
}

// mysql> explain select value from host_metric_values where host_id='f8775dfd1af' and name='loadavg5' and 1558236780 > time and time > 1558236660;
// +----+-------------+--------------------+------------+-------+---------------+---------+---------+------+------+----------+-------------+
// | id | select_type | table              | partitions | type  | possible_keys | key     | key_len | ref  | rows | filtered | Extra       |
// +----+-------------+--------------------+------------+-------+---------------+---------+---------+------+------+----------+-------------+
// |  1 | SIMPLE      | host_metric_values | NULL       | range | PRIMARY       | PRIMARY | 588     | NULL |    1 |   100.00 | Using where |
// +----+-------------+--------------------+------------+-------+---------------+---------+---------+------+------+----------+-------------+
// 1 row in set, 1 warning (0.01 sec)
func (repo *HostRepository) MetricValues(orgID, hostID, name string, from, to int64) (*domain.MetricValues, error) {
	values := make([]domain.MetricValue, 0)
	if err := repo.Transact(func(tx Tx) error {
		rows, err := tx.Query("select time, value from host_metric_values where org_id=? and host_id=? and name=? and ? < time and time < ?", orgID, hostID, name, from, to)
		if err != nil {
			return fmt.Errorf("select time, value from host_metric_values: %v", err)
		}
		defer rows.Close()

		for rows.Next() {
			var time int64
			var value float64
			if err := rows.Scan(&time, &value); err != nil {
				return fmt.Errorf("scan: %v", err)
			}

			values = append(values, domain.MetricValue{
				HostID: hostID,
				Name:   name,
				Time:   time,
				Value:  value,
			})
		}

		return nil
	}); err != nil {
		return nil, fmt.Errorf("transaction: %v", err)
	}

	return &domain.MetricValues{Metrics: values}, nil
}

// mysql> explain select time, value from host_metric_values where org_id='default' and host_id='ceb7c8b51c0' and name='loadavg5' order by time desc limit 3;
// +----+-------------+--------------------+------------+------+---------------+---------+---------+-------------+------+----------+-------------+
// | id | select_type | table              | partitions | type | possible_keys | key     | key_len | ref         | rows | filtered | Extra       |
// +----+-------------+--------------------+------------+------+---------------+---------+---------+-------------+------+----------+-------------+
// |  1 | SIMPLE      | host_metric_values | NULL       | ref  | PRIMARY       | PRIMARY | 580     | const,const |    4 |    10.00 | Using where |
// +----+-------------+--------------------+------------+------+---------------+---------+---------+-------------+------+----------+-------------+
// 1 row in set, 1 warning (0.01 sec)
func (repo *HostRepository) MetricValuesLimit(orgID, hostID, name string, limit int) (*domain.MetricValues, error) {
	values := make([]domain.MetricValue, 0)
	if err := repo.Transact(func(tx Tx) error {
		rows, err := tx.Query("select time, value from host_metric_values where org_id=? and host_id=? and name=? order by time desc limit ?", orgID, hostID, name, limit)
		if err != nil {
			return fmt.Errorf("select time, value from host_metric_values: %v", err)
		}
		defer rows.Close()

		for rows.Next() {
			var time int64
			var value float64
			if err := rows.Scan(&time, &value); err != nil {
				return fmt.Errorf("scan: %v", err)
			}

			values = append(values, domain.MetricValue{
				HostID: hostID,
				Name:   name,
				Time:   time,
				Value:  value,
			})
		}

		return nil
	}); err != nil {
		return nil, fmt.Errorf("transaction: %v", err)
	}

	return &domain.MetricValues{Metrics: values}, nil
}

// mysql> explain select * from host_metric_values_latest where host_id in('27b9dad3197') and name in('loadavg5', 'loadavg15');
// +----+-------------+---------------------------+------------+-------+---------------+---------+---------+------+------+----------+-------------+
// | id | select_type | table                     | partitions | type  | possible_keys | key     | key_len | ref  | rows | filtered | Extra       |
// +----+-------------+---------------------------+------------+-------+---------------+---------+---------+------+------+----------+-------------+
// |  1 | SIMPLE      | host_metric_values_latest | NULL       | range | PRIMARY       | PRIMARY | 580     | NULL |    2 |   100.00 | Using where |
// +----+-------------+---------------------------+------------+-------+---------------+---------+---------+------+------+----------+-------------+
// 1 row in set, 1 warning (0.01 sec)
// mysql> explain select * from host_metric_values where host_id='84a14dbbdc7' and (name='loadavg5' or name='loadavg15') and time in (select max(time) from host_metric_values group by host_id, name);
// +----+-------------+--------------------+------------+-------+---------------+---------+---------+------+------+----------+--------------------------+
// | id | select_type | table              | partitions | type  | possible_keys | key     | key_len | ref  | rows | filtered | Extra                    |
// +----+-------------+--------------------+------------+-------+---------------+---------+---------+------+------+----------+--------------------------+
// |  1 | PRIMARY     | host_metric_values | NULL       | range | PRIMARY       | PRIMARY | 580     | NULL |   62 |   100.00 | Using where              |
// |  2 | SUBQUERY    | host_metric_values | NULL       | range | PRIMARY       | PRIMARY | 580     | NULL |  122 |   100.00 | Using index for group-by |
// +----+-------------+--------------------+------------+-------+---------------+---------+---------+------+------+----------+--------------------------+
// 2 rows in set, 1 warning (0.01 sec)
//
// mysql> explain select host_id, name, value from host_metric_values where time in (select max(time) from host_metric_values group by host_id, name);
// +----+-------------+--------------------+------------+-------+---------------+---------+---------+------+------+----------+--------------------------+
// | id | select_type | table              | partitions | type  | possible_keys | key     | key_len | ref  | rows | filtered | Extra                    |
// +----+-------------+--------------------+------------+-------+---------------+---------+---------+------+------+----------+--------------------------+
// |  1 | PRIMARY     | host_metric_values | NULL       | ALL   | NULL          | NULL    | NULL    | NULL | 3078 |   100.00 | Using where              |
// |  2 | SUBQUERY    | host_metric_values | NULL       | range | PRIMARY       | PRIMARY | 580     | NULL |  116 |   100.00 | Using index for group-by |
// +----+-------------+--------------------+------------+-------+---------------+---------+---------+------+------+----------+--------------------------+
// 2 rows in set, 1 warning (0.00 sec)
func (repo *HostRepository) MetricValuesLatest(orgID string, hostID, name []string) (*domain.TSDBLatest, error) {
	latest := make(map[string]map[string]float64)
	if err := repo.Transact(func(tx Tx) error {
		rows, err := tx.Query("select * from host_metric_values_latest where org_id=?", orgID)
		// TODO multiple ?
		//	rows, err := tx.Query("select * from host_metric_values_latest where host_id in(?) and name in(?)", hostID[0], name[0])
		if err != nil {
			return fmt.Errorf("select * from host_metric_value_latest: %v", err)
		}
		defer rows.Close()

		for rows.Next() {
			var orgID, hostID, name string
			var value float64
			if err := rows.Scan(&orgID, &hostID, &name, &value); err != nil {
				return fmt.Errorf("scan: %v", err)
			}

			if _, ok := latest[hostID]; !ok {
				latest[hostID] = make(map[string]float64)
			}

			latest[hostID][name] = value
		}

		return nil
	}); err != nil {
		return nil, fmt.Errorf("transaction: %v", err)
	}

	return &domain.TSDBLatest{TSDBLatest: latest}, nil
}

// insert into host_metric_values values(${host_id}, ${name}, ${time}, ${value})
func (repo *HostRepository) SaveMetricValues(orgID string, values []domain.MetricValue) (*domain.Success, error) {
	if err := repo.Transact(func(tx Tx) error {
		for i := range values {
			if _, err := tx.Exec(
				"insert into host_metric_values values(?, ?, ?, ?, ?)",
				orgID,
				values[i].HostID,
				values[i].Name,
				values[i].Time,
				values[i].Value,
			); err != nil {
				return fmt.Errorf("insert into host_metric_values: %v", err)
			}

			if _, err := tx.Exec(
				`
				insert into host_metric_values_latest (
					org_id,
					host_id,
					name,
					value
				)
				values (?, ?, ?, ?)
				on duplicate key update
					value = values(value)
				`,
				orgID,
				values[i].HostID,
				values[i].Name,
				values[i].Value,
			); err != nil {
				return fmt.Errorf("insert into host_metric_values_latest: %v", err)
			}
		}

		return nil
	}); err != nil {
		return &domain.Success{Success: false}, fmt.Errorf("transaction: %v", err)
	}

	return &domain.Success{Success: true}, nil
}

func (repo *HostRepository) MetricValuesAverage(orgID, hostID, name string, duration int) (*domain.MetricValueAverage, error) {
	avg := &domain.MetricValueAverage{
		OrgID:    orgID,
		HostID:   hostID,
		Name:     name,
		Duration: duration,
	}

	if err := repo.Transact(func(tx Tx) error {

		row := tx.QueryRow(`
			select
				max(latest.time),
				avg(latest.value)
			from (
				select
					time,
					value
				from
					host_metric_values
				where
					org_id=?  and
					host_id=? and
					name=?
				order by
					time desc
				limit ?
				) as latest
			`,
			orgID,
			hostID,
			name,
			duration,
		)

		if err := row.Scan(
			&avg.Time,
			&avg.Value,
		); err != nil {
			return fmt.Errorf("scan: %v", err)
		}

		return nil
	}); err != nil {
		return avg, fmt.Errorf("transaction: %v", err)
	}

	return avg, nil
}

// select * from host_metadata where host_id=${hostID} and namespace=${namespace} limit=1
func (repo *HostRepository) ExistsMetadata(orgID, hostID, namespace string) bool {
	rows, err := repo.Query("select 1 from host_meta where org_id=? and host_id=? and namespace=?", orgID, hostID, namespace)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	if rows.Next() {
		return true
	}

	return false
}

// select namespace from host_metadata where host_id=${hostID}
func (repo *HostRepository) MetadataList(orgID, hostID string) (*domain.HostMetadataList, error) {
	values := make([]domain.Namespace, 0)
	if err := repo.Transact(func(tx Tx) error {
		rows, err := tx.Query("select namespace from host_meta where org_id=? and host_id=?", orgID, hostID)
		if err != nil {
			return fmt.Errorf("select namespace from host_meta: %v", err)
		}
		defer rows.Close()

		for rows.Next() {
			var namespace string
			if err := rows.Scan(
				&namespace,
			); err != nil {
				return fmt.Errorf("scan: %v", err)
			}

			values = append(values, domain.Namespace{Namespace: namespace})
		}

		return nil
	}); err != nil {
		return nil, fmt.Errorf("transaction: %v", err)
	}

	return &domain.HostMetadataList{Metadata: values}, nil
}

// select * from host_metadata where host_id=${hostID} and namespace=${namespace}
func (repo *HostRepository) Metadata(orgID, hostID, namespace string) (interface{}, error) {
	var meta string
	if err := repo.Transact(func(tx Tx) error {
		row := tx.QueryRow("select meta from host_meta where org_id=? and host_id=? and namespace=?", orgID, hostID, namespace)
		if err := row.Scan(&meta); err != nil {
			return fmt.Errorf("scan: %v", err)
		}

		return nil
	}); err != nil {
		return nil, fmt.Errorf("transaction: %v", err)
	}

	var out interface{}
	if err := json.Unmarshal([]byte(meta), &out); err != nil {
		return nil, fmt.Errorf("unmarshal: %v", err)
	}

	return out, nil
}

// insert into host_metadata values(${hostID}, ${namespace}, ${metadata})
func (repo *HostRepository) SaveMetadata(orgID, hostID, namespace string, metadata interface{}) (*domain.Success, error) {
	meta, err := json.Marshal(metadata)
	if err != nil {
		return &domain.Success{Success: false}, fmt.Errorf("marshal: %v", err)
	}

	if err := repo.Transact(func(tx Tx) error {
		if _, err := tx.Exec(
			"insert into host_meta values(?, ?, ?, ?)",
			orgID,
			hostID,
			namespace,
			string(meta),
		); err != nil {
			return fmt.Errorf("insert into host_meta: %v", err)
		}

		return nil
	}); err != nil {
		return &domain.Success{Success: false}, fmt.Errorf("transaction: %v", err)
	}

	return &domain.Success{Success: true}, nil
}

// delete from host_metadata where host_id=${hostID} and namespace=${namespace}
func (repo *HostRepository) DeleteMetadata(orgID, hostID, namespace string) (*domain.Success, error) {
	if err := repo.Transact(func(tx Tx) error {
		if _, err := tx.Exec(
			"delete from host_meta where org_id=? and host_id=? and namespace=?",
			orgID,
			hostID,
			namespace,
		); err != nil {
			return fmt.Errorf("delete from host_meta: %v", err)
		}

		return nil
	}); err != nil {
		return &domain.Success{Success: false}, fmt.Errorf("transaction: %v", err)
	}

	return &domain.Success{Success: true}, nil
}
