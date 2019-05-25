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
				org               varchar(64)  not null,
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
				org       varchar(64)  not null,
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
				org     varchar(64)  not null,
				host_id varchar(16)  not null,
				name    varchar(128) not null,
				time    bigint not null,
				value   double not null,
				primary key(host_id, name, time)
			)
			`,
		); err != nil {
			return fmt.Errorf("create table host_metric_values: %v", err)
		}

		if _, err := tx.Exec(
			`
			create table if not exists host_metric_values_latest (
				org     varchar(64)  not null,
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

// mysql> explain select * from hosts;
// +----+-------------+-------+------------+------+---------------+------+---------+------+------+----------+-------+
// | id | select_type | table | partitions | type | possible_keys | key  | key_len | ref  | rows | filtered | Extra |
// +----+-------------+-------+------------+------+---------------+------+---------+------+------+----------+-------+
// |  1 | SIMPLE      | hosts | NULL       | ALL  | NULL          | NULL | NULL    | NULL |    4 |   100.00 | NULL  |
// +----+-------------+-------+------------+------+---------------+------+---------+------+------+----------+-------+
// 1 row in set, 1 warning (0.01 sec)
func (repo *HostRepository) List(org string) (*domain.Hosts, error) {
	hosts := make([]domain.Host, 0)

	if err := repo.Transact(func(tx Tx) error {
		rows, err := tx.Query("select * from hosts where org=?", org)
		if err != nil {
			return fmt.Errorf("select * from hosts: %v", err)
		}
		defer rows.Close()

		for rows.Next() {
			var host domain.Host
			var org, roles, roleFullnames, interfaces, checks, meta string
			if err := rows.Scan(
				&org,
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

// insert into hosts values(${name}, ${meta}, ${interfaces}, ${checks}, ${display_name}, ${custom_identifier}, ${created_at}, ${id}, ${status}, ${memo}, ${roles}, ${is_retired}, ${retired_at} )
func (repo *HostRepository) Save(org string, host *domain.Host) (*domain.HostID, error) {
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
				org,
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
			org,
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
					org,
					name
				)
				select ?, ? where not exists (select 1 from services where org=? and name=?)
				`,
				org,
				svc,
				org,
				svc,
			); err != nil {
				return fmt.Errorf("insert into services: %v", err)
			}

			for i := range role {
				if _, err := tx.Exec(
					`
					insert into roles (
						org,
						service_name,
						name
					)
					select ?, ?, ? where not exists (select 1 from roles where org=? and service_name=? and name=?)
					`,
					org,
					svc,
					role[i],
					org,
					svc,
					role[i],
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
func (repo *HostRepository) Host(org, hostID string) (*domain.Host, error) {
	var host domain.Host

	if err := repo.Transact(func(tx Tx) error {
		row := tx.QueryRow("select * from hosts where org=? and id=?", org, hostID)
		var org, roles, roleFullnames, interfaces, checks, meta string
		if err := row.Scan(
			&org,
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
func (repo *HostRepository) Exists(org, hostID string) bool {
	rows, err := repo.Query("select 1 from hosts where org=? and id=? limit 1", org, hostID)
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
func (repo *HostRepository) Status(org, hostID, status string) (*domain.Success, error) {
	if err := repo.Transact(func(tx Tx) error {
		if _, err := tx.Exec("update hosts set status=? where org=? and id=?", status, org, hostID); err != nil {
			return fmt.Errorf("update hosts: %v", err)
		}

		return nil
	}); err != nil {
		return &domain.Success{Success: false}, nil
	}

	return &domain.Success{Success: true}, nil
}

// update hosts set roles=${roles} where id=${hostID}
func (repo *HostRepository) SaveRoleFullNames(org, hostID string, names *domain.RoleFullNames) (*domain.Success, error) {
	roles := names.Roles()

	mroles, err := json.Marshal(roles)
	if err != nil {
		return &domain.Success{Success: false}, nil
	}

	roleFullnames, err := json.Marshal(names.Names)
	if err != nil {
		return &domain.Success{Success: false}, nil
	}

	if err := repo.Transact(func(tx Tx) error {
		if _, err := tx.Exec(
			"update hosts set role_fullnames=?, roles=? where org=? and id=?",
			string(roleFullnames),
			string(mroles),
			org,
			hostID,
		); err != nil {
			return fmt.Errorf("update hosts: %v", err)
		}

		for svc, role := range roles {
			if _, err := tx.Exec(
				`
				insert into services (
					org,
					name
				)
				select ?, ? where not exists (select 1 from services where org=? and name=?)
				`,
				org,
				svc,
				org,
				svc,
			); err != nil {
				return fmt.Errorf("insert into services: %v", err)
			}

			for i := range role {
				if _, err := tx.Exec(
					`
					insert into roles (
						org,
						service_name,
						name
					)
					select ?, ?, ? where not exists (select 1 from roles where org=? and service_name=? and name=?)
					`,
					org,
					svc,
					role[i],
					org,
					svc,
					role[i],
				); err != nil {
					return fmt.Errorf("insert into roles: %v", err)
				}

			}
		}

		return nil
	}); err != nil {
		return &domain.Success{Success: false}, nil
	}

	return &domain.Success{Success: true}, nil
}

// update hosts set is_retired=true, retired_at=time.Now().Unix() where id=${hostID}
func (repo *HostRepository) Retire(org, hostID string, retire *domain.HostRetire) (*domain.Success, error) {
	if err := repo.Transact(func(tx Tx) error {
		if _, err := tx.Exec("update hosts set is_retired=?, retired_at=? where org=? and id=?", true, time.Now().Unix(), org, hostID); err != nil {
			return fmt.Errorf("update hosts: %v", err)
		}

		return nil
	}); err != nil {
		return &domain.Success{Success: false}, nil
	}

	return &domain.Success{Success: true}, nil
}

// mysql> explain select * from host_metric_values where host_id='f8775dfd1af' and name='loadavg5' limit 1;
// +----+-------------+--------------------+------------+------+---------------+---------+---------+-------------+------+----------+-------+
// | id | select_type | table              | partitions | type | possible_keys | key     | key_len | ref         | rows | filtered | Extra |
// +----+-------------+--------------------+------------+------+---------------+---------+---------+-------------+------+----------+-------+
// |  1 | SIMPLE      | host_metric_values | NULL       | ref  | PRIMARY       | PRIMARY | 580     | const,const |    5 |   100.00 | NULL  |
// +----+-------------+--------------------+------------+------+---------------+---------+---------+-------------+------+----------+-------+
// 1 row in set, 1 warning (0.00 sec)
func (repo *HostRepository) ExistsMetric(org, hostID, name string) bool {
	rows, err := repo.Query("select 1 from host_metric_values where org=? and host_id=? and name=? limit 1", org, hostID, name)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	if rows.Next() {
		return true
	}

	return false
}

// mysql> explain select distinct name from host_metric_values where host_id='f8775dfd1af';
// +----+-------------+--------------------+------------+------+---------------+---------+---------+-------+------+----------+-------------+
// | id | select_type | table              | partitions | type | possible_keys | key     | key_len | ref   | rows | filtered | Extra       |
// +----+-------------+--------------------+------------+------+---------------+---------+---------+-------+------+----------+-------------+
// |  1 | SIMPLE      | host_metric_values | NULL       | ref  | PRIMARY       | PRIMARY | 66      | const |  170 |   100.00 | Using index |
// +----+-------------+--------------------+------------+------+---------------+---------+---------+-------+------+----------+-------------+
// 1 row in set, 1 warning (0.01 sec)
func (repo *HostRepository) MetricNames(org, hostID string) (*domain.MetricNames, error) {
	names := make([]string, 0)
	if err := repo.Transact(func(tx Tx) error {
		rows, err := tx.Query("select distinct name from host_metric_values where org=? and host_id=?", org, hostID)
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
func (repo *HostRepository) MetricValues(org, hostID, name string, from, to int64) (*domain.MetricValues, error) {
	values := make([]domain.MetricValue, 0)
	if err := repo.Transact(func(tx Tx) error {
		rows, err := tx.Query("select time, value from host_metric_values where org=? and host_id=? and name=? and ? < time and time < ?", org, hostID, name, from, to)
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
func (repo *HostRepository) MetricValuesLatest(org string, hostID, name []string) (*domain.TSDBLatest, error) {
	latest := make(map[string]map[string]float64)
	if err := repo.Transact(func(tx Tx) error {
		rows, err := tx.Query("select * from host_metric_values_latest where org=?", org)
		// TODO multiple ?
		//	rows, err := tx.Query("select * from host_metric_values_latest where host_id in(?) and name in(?)", hostID[0], name[0])
		if err != nil {
			return fmt.Errorf("select * from host_metric_value_latest: %v", err)
		}
		defer rows.Close()

		for rows.Next() {
			var hostID, name string
			var value float64
			if err := rows.Scan(&hostID, &name, &value); err != nil {
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
func (repo *HostRepository) SaveMetricValues(org string, values []domain.MetricValue) (*domain.Success, error) {
	if err := repo.Transact(func(tx Tx) error {
		for i := range values {
			if _, err := tx.Exec(
				"insert into host_metric_values values(?, ?, ?, ?, ?)",
				org,
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
					org,
					host_id,
					name,
					value
				)
				values (?, ?, ?, ?)
				on duplicate key update
					value = values(value)
				`,
				org,
				values[i].HostID,
				values[i].Name,
				values[i].Value,
			); err != nil {
				return fmt.Errorf("insert into host_metric_values_latest: %v", err)
			}
		}

		return nil
	}); err != nil {
		return &domain.Success{Success: false}, nil
	}

	return &domain.Success{Success: true}, nil
}

// select * from host_metadata where host_id=${hostID} and namespace=${namespace} limit=1
func (repo *HostRepository) ExistsMetadata(org, hostID, namespace string) bool {
	rows, err := repo.Query("select 1 from host_meta where org=? and host_id=? and namespace=?", org, hostID, namespace)
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
func (repo *HostRepository) MetadataList(org, hostID string) (*domain.HostMetadataList, error) {
	values := make([]domain.Namespace, 0)
	if err := repo.Transact(func(tx Tx) error {
		rows, err := tx.Query("select namespace from host_meta where org=? and host_id=?", org, hostID)
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
func (repo *HostRepository) Metadata(org, hostID, namespace string) (interface{}, error) {
	var meta string
	if err := repo.Transact(func(tx Tx) error {
		row := tx.QueryRow("select meta from host_meta where org=? and host_id=? and namespace=?", org, hostID, namespace)
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
func (repo *HostRepository) SaveMetadata(org, hostID, namespace string, metadata interface{}) (*domain.Success, error) {
	meta, err := json.Marshal(metadata)
	if err != nil {
		return &domain.Success{Success: false}, nil
	}

	if err := repo.Transact(func(tx Tx) error {
		if _, err := tx.Exec(
			"insert into host_meta values(?, ?, ?, ?)",
			org,
			hostID,
			namespace,
			string(meta),
		); err != nil {
			return fmt.Errorf("insert into host_meta: %v", err)
		}

		return nil
	}); err != nil {
		return &domain.Success{Success: false}, nil
	}

	return &domain.Success{Success: true}, nil
}

// delete from host_metadata where host_id=${hostID} and namespace=${namespace}
func (repo *HostRepository) DeleteMetadata(org, hostID, namespace string) (*domain.Success, error) {
	if err := repo.Transact(func(tx Tx) error {
		if _, err := tx.Exec(
			"delete from host_meta where org=? and host_id=? and namespace=?",
			hostID,
			namespace,
		); err != nil {
			return fmt.Errorf("delete from host_meta: %v", err)
		}

		return nil
	}); err != nil {
		return &domain.Success{Success: false}, nil
	}

	return &domain.Success{Success: true}, nil
}
