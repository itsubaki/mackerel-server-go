package database

import (
	"encoding/json"
	"strings"
	"time"

	"github.com/itsubaki/mackerel-api/pkg/domain"
)

type HostRepository struct {
	SQLHandler
}

func NewHostRepository(handler SQLHandler) *HostRepository {
	err := handler.Transact(func(tx Tx) error {
		if _, err := tx.Exec(
			`
			create table if not exists hosts (
				id varchar(16) not null primary key,
				name varchar(128),
				status varchar(16),
				memo varchar(128),
				display_name varchar(128),
				custom_identifier varchar(128),
				created_at bigint,
				retired_at bigint,
				is_retired boolean,
				roles text,
				role_fullnames text,
				interfaces text,
				checks text,
				meta text
			)
			`,
		); err != nil {
			return err
		}

		if _, err := tx.Exec(
			`
			create table if not exists host_metric_values (
				host_id varchar(16) not null,
				name varchar(128) not null,
				time bigint not null,
				value double not null,
				primary key(host_id,name,time)
			)
			`,
		); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		panic(err)
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
func (repo *HostRepository) List() (*domain.Hosts, error) {
	var hosts []domain.Host

	err := repo.Transact(func(tx Tx) error {
		rows, err := tx.Query("select * from hosts")
		if err != nil {
			return err
		}
		defer rows.Close()

		for rows.Next() {
			var host domain.Host
			var roles, roleFullnames, interfaces, checks, meta string
			if err := rows.Scan(
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
				return err
			}

			if err := json.Unmarshal([]byte(roles), &host.Roles); err != nil {
				return err
			}

			if err := json.Unmarshal([]byte(roleFullnames), &host.RoleFullNames); err != nil {
				return err
			}

			if err := json.Unmarshal([]byte(interfaces), &host.Interfaces); err != nil {
				return err
			}

			if err := json.Unmarshal([]byte(checks), &host.Checks); err != nil {
				return err
			}

			if err := json.Unmarshal([]byte(meta), &host.Meta); err != nil {
				return err
			}

			hosts = append(hosts, host)
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return &domain.Hosts{Hosts: hosts}, nil
}

// insert into hosts values(${name}, ${meta}, ${interfaces}, ${checks}, ${display_name}, ${custom_identifier}, ${created_at}, ${id}, ${status}, ${memo}, ${roles}, ${is_retired}, ${retired_at} )
func (repo *HostRepository) Save(host *domain.Host) (*domain.HostID, error) {
	err := repo.Transact(func(tx Tx) error {
		roles, err := json.Marshal(host.Roles)
		if err != nil {
			return err
		}

		roleFullnames, err := json.Marshal(host.RoleFullNames)
		if err != nil {
			return err
		}

		interfaces, err := json.Marshal(host.Interfaces)
		if err != nil {
			return err
		}

		checks, err := json.Marshal(host.Checks)
		if err != nil {
			return err
		}

		meta, err := json.Marshal(host.Meta)
		if err != nil {
			return err
		}

		if _, err := tx.Exec(
			`
			insert into hosts (
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
			values (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
			on duplicate key update
				name = values(name),
				memo = values(memo),
				display_name = values(display_name),
				custom_identifier = values(custom_identifier),
				roles = values(roles),
				role_fullnames = values(role_fullnames),
				interfaces = values(interfaces),
				checks = values(checks),
				meta = values(meta)
			`,
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
			return err
		}

		return nil
	})
	if err != nil {
		return nil, err
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
func (repo *HostRepository) Host(hostID string) (*domain.Host, error) {
	var host domain.Host

	err := repo.Transact(func(tx Tx) error {
		row := tx.QueryRow("select * from hosts where id=?", hostID)
		var roles, roleFullnames, interfaces, checks, meta string
		if err := row.Scan(
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
			return err
		}

		if err := json.Unmarshal([]byte(roles), &host.Roles); err != nil {
			return err
		}

		if err := json.Unmarshal([]byte(roleFullnames), &host.RoleFullNames); err != nil {
			return err
		}

		if err := json.Unmarshal([]byte(interfaces), &host.Interfaces); err != nil {
			return err
		}

		if err := json.Unmarshal([]byte(checks), &host.Checks); err != nil {
			return err
		}

		if err := json.Unmarshal([]byte(meta), &host.Meta); err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return nil, err
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
func (repo *HostRepository) Exists(hostID string) bool {
	rows, err := repo.Query("select * from hosts where id=? limit 1", hostID)
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
func (repo *HostRepository) Status(hostID, status string) (*domain.Success, error) {
	if err := repo.Transact(func(tx Tx) error {
		if _, err := tx.Exec("update hosts set status=? where id=?", status, hostID); err != nil {
			return err
		}

		return nil
	}); err != nil {
		return &domain.Success{Success: false}, err
	}

	return &domain.Success{Success: true}, nil
}

// update hosts set roles=${roles} where id=${hostID}
func (repo *HostRepository) SaveRoleFullNames(hostID string, names *domain.RoleFullNames) (*domain.Success, error) {
	roleFullnames, err := json.Marshal(names.Names)
	if err != nil {
		return &domain.Success{Success: false}, err
	}

	roles := make(map[string][]string)
	for i := range names.Names {
		svc := strings.Split(names.Names[i], ":")
		if _, ok := roles[svc[0]]; !ok {
			roles[svc[0]] = []string{}
		}
		roles[svc[0]] = append(roles[svc[0]], svc[1])
	}

	broles, err := json.Marshal(roles)
	if err != nil {
		return &domain.Success{Success: false}, err
	}

	if err := repo.Transact(func(tx Tx) error {
		if _, err := tx.Exec(
			"update hosts set role_fullnames=?, roles=? where id=?",
			string(roleFullnames),
			string(broles),
			hostID,
		); err != nil {
			return err
		}

		return nil
	}); err != nil {
		return &domain.Success{Success: false}, err
	}

	return &domain.Success{Success: true}, nil
}

// update hosts set is_retired=true, retired_at=time.Now().Unix() where id=${hostID}
func (repo *HostRepository) Retire(hostID string, retire *domain.HostRetire) (*domain.Success, error) {
	if err := repo.Transact(func(tx Tx) error {
		if _, err := tx.Exec("update hosts set is_retired=?, retired_at=? where id=?", true, time.Now().Unix(), hostID); err != nil {
			return err
		}

		return nil
	}); err != nil {
		return &domain.Success{Success: false}, err
	}

	return &domain.Success{Success: true}, nil
}

// select * from host_metric_values where host_id=${hostID} and name=${name} limit=1
func (repo *HostRepository) ExistsMetric(hostID, name string) bool {
	return false
}

// select distinct name from host_metric_values where host_id=${hostID}
func (repo *HostRepository) MetricNames(hostID string) (*domain.MetricNames, error) {
	return &domain.MetricNames{}, nil
}

// select value from host_metric_values where host_id=${hostID} and name=${name} and ${from} < from and to < ${to}
func (repo *HostRepository) MetricValues(hostID, name string, from, to int64) (*domain.MetricValues, error) {
	return &domain.MetricValues{}, nil
}

// select * from host_metric_values_latest where host_id=${hostID} and name=${name}
func (repo *HostRepository) MetricValuesLatest(hostID, name []string) (*domain.TSDBLatest, error) {
	return &domain.TSDBLatest{}, nil
}

// insert into host_metric_values values(${host_id}, ${name}, ${time}, ${value})
func (repo *HostRepository) SaveMetricValues(values []domain.MetricValue) (*domain.Success, error) {
	if err := repo.Transact(func(tx Tx) error {
		for i := range values {
			if _, err := tx.Exec(
				"insert into host_metric_values values(?, ?, ?, ?)",
				values[i].HostID,
				values[i].Name,
				values[i].Time,
				values[i].Value,
			); err != nil {
				return err
			}
		}

		return nil
	}); err != nil {
		return &domain.Success{Success: false}, err
	}

	return &domain.Success{Success: true}, nil
}

// select * from host_metadata where host_id=${hostID} and namespace=${namespace} limit=1
func (repo *HostRepository) ExistsMetadata(hostID, namespace string) bool {
	return true
}

// select namespace from host_metadata where host_id=${hostID}
func (repo *HostRepository) MetadataList(hostID string) (*domain.HostMetadataList, error) {
	return &domain.HostMetadataList{}, nil
}

// select * from host_metadata where host_id=${hostID} and namespace=${namespace}
func (repo *HostRepository) Metadata(hostID, namespace string) (interface{}, error) {
	return "", nil
}

// insert into host_metadata values(${hostID}, ${namespace}, ${metadata})
func (repo *HostRepository) SaveMetadata(hostID, namespace string, metadata interface{}) (*domain.Success, error) {
	return &domain.Success{Success: true}, nil
}

// delete from host_metadata where host_id=${hostID} and namespace=${namespace}
func (repo *HostRepository) DeleteMetadata(hostID, namespace string) (*domain.Success, error) {
	return &domain.Success{Success: true}, nil
}
