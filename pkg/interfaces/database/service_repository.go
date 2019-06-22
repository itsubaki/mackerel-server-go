package database

import (
	"encoding/json"
	"fmt"

	"github.com/itsubaki/mackerel-api/pkg/domain"
)

type ServiceRepository struct {
	SQLHandler
}

func NewServiceRepository(handler SQLHandler) *ServiceRepository {
	if err := handler.Transact(func(tx Tx) error {
		if _, err := tx.Exec(
			`
			create table if not exists services (
				org_id varchar(64)  not null,
				name   varchar(128) not null,
				memo   varchar(128) not null default '',
				primary key(org_id, name)
			)
			`,
		); err != nil {
			return fmt.Errorf("create table services: %v", err)
		}

		if _, err := tx.Exec(
			`
			create table if not exists roles (
				org_id       varchar(64)  not null,
				service_name varchar(128) not null,
				name         varchar(128) not null,
				memo         varchar(128) not null default '',
				primary key(org_id, service_name, name),
				foreign key fk_service_name (org_id, service_name) references services(org_id, name) on delete cascade on update cascade
			)
			`,
		); err != nil {
			return fmt.Errorf("create table roles: %v", err)
		}

		if _, err := tx.Exec(
			`
			create table if not exists service_metric_values (
				org_id       varchar(64)  not null,
				service_name varchar(16)  not null,
				name         varchar(128) not null,
				time         bigint not null,
				value        double not null,
				primary key(org_id, service_name, name, time desc)
			)
			`,
		); err != nil {
			return fmt.Errorf("create table service_metric_values: %v", err)
		}

		if _, err := tx.Exec(
			`
			create table if not exists service_meta (
				org_id       varchar(64)  not null,
				service_name varchar(16)  not null,
				namespace    varchar(128) not null,
				meta         text,
				primary key(org_id, service_name, namespace)
			)
			`,
		); err != nil {
			return fmt.Errorf("create table service_meta: %v", err)
		}

		if _, err := tx.Exec(
			`
			create table if not exists role_meta (
				org_id       varchar(64)  not null,
				service_name varchar(16)  not null,
				role_name    varchar(16)  not null,
				namespace    varchar(128) not null,
				meta         text,
				primary key(org_id, role_name, namespace)
			)
			`,
		); err != nil {
			return fmt.Errorf("create table role_meta: %v", err)
		}

		return nil
	}); err != nil {
		panic(fmt.Errorf("transaction: %v", err))
	}

	return &ServiceRepository{
		SQLHandler: handler,
	}
}

// mysql> explain select service_name, name from roles;
// +----+-------------+-------+------------+-------+---------------+---------+---------+------+------+----------+-------------+
// | id | select_type | table | partitions | type  | possible_keys | key     | key_len | ref  | rows | filtered | Extra       |
// +----+-------------+-------+------------+-------+---------------+---------+---------+------+------+----------+-------------+
// |  1 | SIMPLE      | roles | NULL       | index | NULL          | PRIMARY | 1028    | NULL |    4 |   100.00 | Using index |
// +----+-------------+-------+------------+-------+---------------+---------+---------+------+------+----------+-------------+
// 1 row in set, 1 warning (0.00 sec)
// mysql> explain select memo from services where name='My-Service';
// +----+-------------+----------+------------+-------+---------------+---------+---------+-------+------+----------+-------+
// | id | select_type | table    | partitions | type  | possible_keys | key     | key_len | ref   | rows | filtered | Extra |
// +----+-------------+----------+------------+-------+---------------+---------+---------+-------+------+----------+-------+
// |  1 | SIMPLE      | services | NULL       | const | PRIMARY       | PRIMARY | 514     | const |    1 |   100.00 | NULL  |
// +----+-------------+----------+------------+-------+---------------+---------+---------+-------+------+----------+-------+
// 1 row in set, 1 warning (0.01 sec)
func (repo *ServiceRepository) List(orgID string) (*domain.Services, error) {
	services := make([]domain.Service, 0)

	if err := repo.Transact(func(tx Tx) error {
		rows, err := tx.Query("select service_name, name from roles where org_id=?", orgID)
		if err != nil {
			return fmt.Errorf("select service_name, name from roles: %v", err)
		}
		defer rows.Close()

		roles := make(map[string][]string)
		for rows.Next() {
			var service, role string
			if err := rows.Scan(
				&service,
				&role,
			); err != nil {
				return fmt.Errorf("scan roles: %v", err)
			}

			if _, ok := roles[service]; !ok {
				roles[service] = make([]string, 0)
			}

			roles[service] = append(roles[service], role)
		}

		for svc := range roles {
			service := domain.Service{
				Name:  svc,
				Roles: roles[svc],
			}

			row := tx.QueryRow("select memo from services where org_id=? and name=?", orgID, svc)
			if err := row.Scan(
				&service.Memo,
			); err != nil {
				return fmt.Errorf("scan services: %v", err)
			}

			services = append(services, service)
		}

		return nil
	}); err != nil {
		return nil, fmt.Errorf("transaction: %v", err)
	}

	return &domain.Services{Services: services}, nil
}

// insert into services values()
func (repo *ServiceRepository) Save(orgID string, s *domain.Service) error {
	if err := repo.Transact(func(tx Tx) error {
		if _, err := tx.Exec(
			`
			insert into services (
				org_id,
				name,
				memo
			)
			values (?, ?, ?)
			on duplicate key update
				name = values(memo)
			`,
			orgID, s.Name, s.Memo,
		); err != nil {
			return fmt.Errorf("insert into services: %v", err)
		}

		for i := range s.Roles {
			if _, err := tx.Exec(
				`
				insert into roles (
					org_id,
					service_name,
					name
				)
				select ?, ?, ? where not exists (select 1 from roles where org=? and service_name=? and name=?)
				`,
				orgID, s.Name, s.Roles[i],
				orgID, s.Name, s.Roles[i],
			); err != nil {
				return fmt.Errorf("insert into roles: %v", err)
			}
		}

		return nil
	}); err != nil {
		return fmt.Errorf("transaction: %v", err)
	}

	return nil
}

// mysql> explain select name from roles where service_name='My-Service';
// +----+-------------+-------+------------+------+---------------+---------+---------+-------+------+----------+-------------+
// | id | select_type | table | partitions | type | possible_keys | key     | key_len | ref   | rows | filtered | Extra       |
// +----+-------------+-------+------------+------+---------------+---------+---------+-------+------+----------+-------------+
// |  1 | SIMPLE      | roles | NULL       | ref  | PRIMARY       | PRIMARY | 514     | const |    2 |   100.00 | Using index |
// +----+-------------+-------+------------+------+---------------+---------+---------+-------+------+----------+-------------+
// 1 row in set, 1 warning (0.00 sec)
func (repo *ServiceRepository) Service(orgID, serviceName string) (*domain.Service, error) {
	service := domain.Service{
		Roles: []string{},
	}

	if err := repo.Transact(func(tx Tx) error {
		row := tx.QueryRow("select name, memo from services where org_id=? and name=?", orgID, serviceName)
		if err := row.Scan(
			&service.Name,
			&service.Memo,
		); err != nil {
			return fmt.Errorf("scan: %v", err)
		}

		rows, err := tx.Query("select name from roles where org_id=? and service_name=?", orgID, serviceName)
		if err != nil {
			return fmt.Errorf("select name from roles: %v", err)
		}
		defer rows.Close()

		for rows.Next() {
			var name string
			if err := rows.Scan(&name); err != nil {
				return fmt.Errorf("scan: %v", err)
			}
			service.Roles = append(service.Roles, name)
		}

		return nil
	}); err != nil {
		return nil, fmt.Errorf("transaction: %v", err)
	}

	return &service, nil
}

// mysql> explain select * from services where name='My-Service' limit 1;
// +----+-------------+----------+------------+-------+---------------+---------+---------+-------+------+----------+-------+
// | id | select_type | table    | partitions | type  | possible_keys | key     | key_len | ref   | rows | filtered | Extra |
// +----+-------------+----------+------------+-------+---------------+---------+---------+-------+------+----------+-------+
// |  1 | SIMPLE      | services | NULL       | const | PRIMARY       | PRIMARY | 514     | const |    1 |   100.00 | NULL  |
// +----+-------------+----------+------------+-------+---------------+---------+---------+-------+------+----------+-------+
// 1 row in set, 1 warning (0.00 sec)
func (repo *ServiceRepository) Exists(orgID, serviceName string) bool {
	rows, err := repo.Query("select 1 from services where org_id=? and name=? limit 1", orgID, serviceName)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	if rows.Next() {
		return true
	}

	return false
}

// delete from services where service_name=${serviceName}
func (repo *ServiceRepository) Delete(orgID, serviceName string) error {
	if err := repo.Transact(func(tx Tx) error {
		if _, err := tx.Exec("delete from services where org_id=? and name=?", orgID, serviceName); err != nil {
			return fmt.Errorf("delete from services: %v", err)
		}

		return nil
	}); err != nil {
		return fmt.Errorf("transaction: %v", err)
	}

	return nil
}

// select * from service_roles where service_name=${serviceName}
func (repo *ServiceRepository) RoleList(orgID, serviceName string) (*domain.Roles, error) {
	var roles []domain.Role
	if err := repo.Transact(func(tx Tx) error {
		rows, err := tx.Query("select name, memo from roles where org_id=? and service_name=?", orgID, serviceName)
		if err != nil {
			return fmt.Errorf("select name, memo from roles: %v", err)
		}
		defer rows.Close()

		for rows.Next() {
			var role domain.Role
			if err := rows.Scan(
				&role.Name,
				&role.Memo,
			); err != nil {
				return fmt.Errorf("scan: %v", err)
			}

			roles = append(roles, role)
		}

		return nil
	}); err != nil {
		return nil, fmt.Errorf("transaction: %v", err)
	}

	return &domain.Roles{Roles: roles}, nil
}

// insert into service_roles values(${serviceName}, ${roleName}, ${Memo})
func (repo *ServiceRepository) SaveRole(orgID, serviceName string, r *domain.Role) error {
	if err := repo.Transact(func(tx Tx) error {
		if _, err := tx.Exec(
			`
			insert into roles (
				org_id,
				service_name,
				name,
				memo
			)
			values (?, ?, ?, ?)
			on duplicate key update
				service_name = values(service_name),
				name = values(name),
				memo = values(memo)
			`,
			orgID,
			serviceName,
			r.Name,
			r.Memo,
		); err != nil {
			return fmt.Errorf("insert into roles: %v", err)
		}

		return nil
	}); err != nil {
		return fmt.Errorf("transaction: %v", err)
	}

	return nil
}

// select * from service_roles where service_name=${serviceName} and role_name=${roleName}
func (repo *ServiceRepository) ExistsRole(orgID, serviceName, roleName string) bool {
	rows, err := repo.Query("select 1 from roles where org_id=? and service_name=? and name=? limit 1", orgID, serviceName, roleName)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	if rows.Next() {
		return true
	}

	return false
}

// select * from service_roles where service_name=${serviceName} and role_name=${roleName}
func (repo *ServiceRepository) Role(orgID, serviceName, roleName string) (*domain.Role, error) {
	role := domain.Role{
		ServiceName: serviceName,
		Name:        roleName,
	}

	if err := repo.Transact(func(tx Tx) error {
		row := repo.QueryRow("select memo from roles where org_id=? and service_name=? and name=?", orgID, serviceName, roleName)
		if err := row.Scan(
			&role.Memo,
		); err != nil {
			return fmt.Errorf("scan: %v", err)
		}

		return nil
	}); err != nil {
		return nil, fmt.Errorf("transaction: %v", err)
	}

	return &role, nil
}

// delete from service_roles where service_name=${serviceName} and role_name=${roleName}
func (repo *ServiceRepository) DeleteRole(orgID, serviceName, roleName string) error {
	if err := repo.Transact(func(tx Tx) error {
		if _, err := tx.Exec("delete from roles where org_id=? and service_name=? and name=?", orgID, serviceName, roleName); err != nil {
			return fmt.Errorf("delete from roles: %v", err)
		}

		return nil
	}); err != nil {
		return fmt.Errorf("transaction: %v", err)
	}

	return nil
}

// select * from service_metrics where service_name=${serviceName} and name=${metricName}
func (repo *ServiceRepository) ExistsMetric(orgID, serviceName, metricName string) bool {
	rows, err := repo.Query("select 1 from service_metric_values where org_id=? and service_name=? and name=? limit 1", orgID, serviceName, metricName)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	if rows.Next() {
		return true
	}

	return false
}

// select distinct name from service_metrics where service_name=${serviceName}
func (repo *ServiceRepository) MetricNames(orgID, serviceName string) (*domain.ServiceMetricValueNames, error) {
	names := make([]string, 0)
	if err := repo.Transact(func(tx Tx) error {
		rows, err := repo.Query("select distinct name from service_metric_values where org_id=? and service_name=?", orgID, serviceName)
		if err != nil {
			return fmt.Errorf("select distinct name from service_metric_values: %v", err)
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

	return &domain.ServiceMetricValueNames{Names: names}, nil
}

// select * from service_metrics where service_name=${serviceName} and name=${metricName}  and ${from} < from and to < ${to}
func (repo *ServiceRepository) MetricValues(orgID, serviceName, metricName string, from, to int64) (*domain.ServiceMetricValues, error) {
	values := make([]domain.ServiceMetricValue, 0)
	if err := repo.Transact(func(tx Tx) error {
		rows, err := tx.Query("select time, value from service_metric_values where org_id=? and service_name=? and name=? and ? < time and time < ?", orgID, serviceName, metricName, from, to)
		if err != nil {
			return fmt.Errorf("select time, value from service_metric_values: %v", err)
		}
		defer rows.Close()

		for rows.Next() {
			var time int64
			var value float64
			if err := rows.Scan(&time, &value); err != nil {
				return fmt.Errorf("scan: %v", err)
			}

			values = append(values, domain.ServiceMetricValue{
				ServiceName: serviceName,
				Name:        metricName,
				Time:        time,
				Value:       value,
			})
		}

		return nil
	}); err != nil {
		return nil, fmt.Errorf("transaction: %v", err)
	}

	return &domain.ServiceMetricValues{Metrics: values}, nil
}

// insert into service_metrics values(${serviceName}, ${name}, ${time}, ${value})
func (repo *ServiceRepository) SaveMetricValues(orgID, serviceName string, values []domain.ServiceMetricValue) (*domain.Success, error) {
	if err := repo.Transact(func(tx Tx) error {
		for i := range values {
			if _, err := tx.Exec(
				"insert into service_metric_values values(?, ?, ?, ?, ?)",
				orgID,
				serviceName,
				values[i].Name,
				values[i].Time,
				values[i].Value,
			); err != nil {
				return fmt.Errorf("insert into service_metric_values: %v", err)
			}
		}

		return nil
	}); err != nil {
		return &domain.Success{Success: false}, nil
	}

	return &domain.Success{Success: true}, nil
}

// select * from service_metadata where service_name=${serviceName} and namespace=${namespace} limit=1
func (repo *ServiceRepository) ExistsMetadata(orgID, serviceName, namespace string) bool {
	rows, err := repo.Query("select 1 from service_meta where org_id=? and service_name=? and namespace=?", orgID, serviceName, namespace)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	if rows.Next() {
		return true
	}

	return false
}

// select namespace from service_metadata where service_name=${serviceName}
func (repo *ServiceRepository) MetadataList(orgID, serviceName string) (*domain.ServiceMetadataList, error) {
	values := make([]domain.ServiceMetadata, 0)
	if err := repo.Transact(func(tx Tx) error {
		rows, err := tx.Query("select namespace from service_meta where org_id=? and service_name=?", orgID, serviceName)
		if err != nil {
			return fmt.Errorf("select namespace from service_meta: %v", err)
		}
		defer rows.Close()

		for rows.Next() {
			var namespace string
			if err := rows.Scan(
				&namespace,
			); err != nil {
				return fmt.Errorf("scan: %v", err)
			}

			values = append(values, domain.ServiceMetadata{Namespace: namespace})
		}

		return nil
	}); err != nil {
		return nil, fmt.Errorf("transaction: %v", err)
	}

	return &domain.ServiceMetadataList{Metadata: values}, nil
}

// select * from service_metadata where service_name=${serviceName} and namespace=${namespace}
func (repo *ServiceRepository) Metadata(orgID, serviceName, namespace string) (interface{}, error) {
	var meta string
	if err := repo.Transact(func(tx Tx) error {
		row := tx.QueryRow("select meta from service_meta where org_id=? and service_name=? and namespace=?", orgID, serviceName, namespace)
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

// insert into service_metadata values(${serviceName}, ${namespace}, ${metadata})
func (repo *ServiceRepository) SaveMetadata(orgID, serviceName, namespace string, metadata interface{}) (*domain.Success, error) {
	meta, err := json.Marshal(metadata)
	if err != nil {
		return &domain.Success{Success: false}, nil
	}

	if err := repo.Transact(func(tx Tx) error {
		if _, err := tx.Exec(
			"insert into service_meta values(?, ?, ?, ?)",
			orgID,
			serviceName,
			namespace,
			string(meta),
		); err != nil {
			return fmt.Errorf("insert into service_meta: %v", err)
		}

		return nil
	}); err != nil {
		return &domain.Success{Success: false}, nil
	}

	return &domain.Success{Success: true}, nil
}

// delete from service_metadata where service_name=${serviceName} and namespace=${namespace}
func (repo *ServiceRepository) DeleteMetadata(orgID, serviceName, namespace string) (*domain.Success, error) {
	if err := repo.Transact(func(tx Tx) error {
		if _, err := tx.Exec(
			"delete from service_meta where org_id=? and service_name=? and namespace=?",
			orgID,
			serviceName,
			namespace,
		); err != nil {
			return fmt.Errorf("delete from service_meta: %v", err)
		}

		return nil
	}); err != nil {
		return &domain.Success{Success: false}, nil
	}

	return &domain.Success{Success: true}, nil
}

// select * from role_metadata where service_name=${serviceName} and role_name=${roleName} and namespace=${namespace} limit=1
func (repo *ServiceRepository) ExistsRoleMetadata(orgID, serviceName, roleName, namespace string) bool {
	rows, err := repo.Query(
		"select 1 from role_meta where org_id=? and service_name=? and role_name=? and namespace=?",
		orgID,
		serviceName,
		roleName,
		namespace,
	)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	if rows.Next() {
		return true
	}

	return false
}

// select namespace from role_metadata where service_name=${serviceName} and role_name=${roleName}
func (repo *ServiceRepository) RoleMetadataList(orgID, serviceName, roleName string) (*domain.RoleMetadataList, error) {
	values := make([]domain.RoleMetadata, 0)
	if err := repo.Transact(func(tx Tx) error {
		rows, err := tx.Query("select namespace from role_meta where org_id=? and service_name=? and role_name=? ", orgID, serviceName, roleName)
		if err != nil {
			return fmt.Errorf("select namespace from role_meta: %v", err)
		}
		defer rows.Close()

		for rows.Next() {
			var namespace string
			if err := rows.Scan(
				&namespace,
			); err != nil {
				return fmt.Errorf("scan: %v", err)
			}

			values = append(values, domain.RoleMetadata{Namespace: namespace})
		}

		return nil
	}); err != nil {
		return nil, fmt.Errorf("transaction: %v", err)
	}

	return &domain.RoleMetadataList{Metadata: values}, nil
}

// select * from role_metadata where service_name=${serviceName} and role_name=${roleName} and namespace=${namespace}
func (repo *ServiceRepository) RoleMetadata(orgID, serviceName, roleName, namespace string) (interface{}, error) {
	var meta string
	if err := repo.Transact(func(tx Tx) error {
		row := tx.QueryRow("select meta from role_meta where org_id=? and service_name=? and role_name=? and namespace=?", orgID, serviceName, roleName, namespace)
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

// insert into role_metadata values(${serviceName}, ${roleName}, ${namespace}, ${metadata})
func (repo *ServiceRepository) SaveRoleMetadata(orgID, serviceName, roleName, namespace string, metadata interface{}) (*domain.Success, error) {
	meta, err := json.Marshal(metadata)
	if err != nil {
		return &domain.Success{Success: false}, nil
	}

	if err := repo.Transact(func(tx Tx) error {
		if _, err := tx.Exec(
			"insert into role_meta values(?, ?, ?, ?, ?)",
			orgID,
			serviceName,
			roleName,
			namespace,
			string(meta),
		); err != nil {
			return fmt.Errorf("insert into role_meta: %v", err)
		}

		return nil
	}); err != nil {
		return &domain.Success{Success: false}, nil
	}

	return &domain.Success{Success: true}, nil
}

// delete from role_metadata where service_name=${serviceName} and role_name=${roleName} and namespace=${namespace}
func (repo *ServiceRepository) DeleteRoleMetadata(orgID, serviceName, roleName, namespace string) (*domain.Success, error) {
	if err := repo.Transact(func(tx Tx) error {
		if _, err := tx.Exec(
			"delete from role_meta where org_id=? and service_name=? and role_name=? and namespace=?",
			orgID,
			serviceName,
			roleName,
			namespace,
		); err != nil {
			return fmt.Errorf("delete from role_meta: %v", err)
		}

		return nil
	}); err != nil {
		return &domain.Success{Success: false}, nil
	}

	return &domain.Success{Success: true}, nil
}
