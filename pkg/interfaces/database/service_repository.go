package database

import (
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
				name varchar(128) not null primary key,
				memo varchar(128) not null default ''
			)
			`,
		); err != nil {
			return fmt.Errorf("create table services: %v", err)
		}

		if _, err := tx.Exec(
			`
			create table if not exists roles (
				service_name varchar(128) not null,
				name varchar(128) not null,
				memo varchar(128) not null default '',
				primary key(service_name, name),
				foreign key fk_service_name (service_name) references services(name) on delete cascade on update cascade
			)
			`,
		); err != nil {
			return fmt.Errorf("create table roles: %v", err)
		}

		if _, err := tx.Exec(
			`
			create table if not exists service_metric_values (
				service_name varchar(16) not null,
				name varchar(128) not null,
				time bigint not null,
				value double not null,
				primary key(service_name, name, time)
			)
			`,
		); err != nil {
			return fmt.Errorf("create table service_metric_values: %v", err)
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
func (repo *ServiceRepository) List() (*domain.Services, error) {
	services := make([]domain.Service, 0)

	if err := repo.Transact(func(tx Tx) error {
		rows, err := tx.Query("select service_name, name from roles")
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

			row := tx.QueryRow("select memo from services where name=?", svc)
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
func (repo *ServiceRepository) Save(s *domain.Service) error {
	if err := repo.Transact(func(tx Tx) error {
		if _, err := tx.Exec(
			`
			insert into services (
				name,
				memo
			)
			values (?, ?)
			on duplicate key update
				name = values(memo)
			`,
			s.Name,
			s.Memo,
		); err != nil {
			return fmt.Errorf("insert into services: %v", err)
		}

		for i := range s.Roles {
			if _, err := tx.Exec(
				`
				insert into roles (
					service_name,
					name
				)
				select ? ,?
				where not exists (
					select 1 from roles where service_name=? and name=?
				)
				`,
				s.Name,
				s.Roles[i],
				s.Name,
				s.Roles[i],
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
func (repo *ServiceRepository) Service(serviceName string) (*domain.Service, error) {
	service := domain.Service{
		Roles: []string{},
	}

	if err := repo.Transact(func(tx Tx) error {
		row := tx.QueryRow("select * from services where name=?", serviceName)
		if err := row.Scan(
			&service.Name,
			&service.Memo,
		); err != nil {
			return fmt.Errorf("scan: %v", err)
		}

		rows, err := tx.Query("select name from roles where service_name=?", serviceName)
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
func (repo *ServiceRepository) Exists(serviceName string) bool {
	rows, err := repo.Query("select * from services where name=? limit 1", serviceName)
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
func (repo *ServiceRepository) Delete(serviceName string) error {
	if err := repo.Transact(func(tx Tx) error {
		if _, err := tx.Exec("delete from services where name=?", serviceName); err != nil {
			return fmt.Errorf("delete from services: %v", err)
		}

		return nil
	}); err != nil {
		return fmt.Errorf("transaction: %v", err)
	}

	return nil
}

// select * from service_roles where service_name=${serviceName}
func (repo *ServiceRepository) RoleList(serviceName string) (*domain.Roles, error) {
	var roles []domain.Role
	if err := repo.Transact(func(tx Tx) error {
		rows, err := tx.Query("select name, memo from roles where service_name=?", serviceName)
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
func (repo *ServiceRepository) SaveRole(serviceName string, r *domain.Role) error {
	if err := repo.Transact(func(tx Tx) error {
		if _, err := tx.Exec(
			`
			insert into roles (
				service_name,
				name,
				memo
			)
			values (?, ?, ?)
			on duplicate key update
				service_name = values(service_name),
				name = values(name),
				memo = values(memo
			`,
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
func (repo *ServiceRepository) ExistsRole(serviceName, roleName string) bool {
	rows, err := repo.Query("select * from roles where service_name=? and name=? limit 1", serviceName, roleName)
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
func (repo *ServiceRepository) Role(serviceName, roleName string) (*domain.Role, error) {
	role := domain.Role{
		ServiceName: serviceName,
		Name:        roleName,
	}

	if err := repo.Transact(func(tx Tx) error {
		row := repo.QueryRow("select memo from roles where service_name=? and name=?", serviceName, roleName)
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
func (repo *ServiceRepository) DeleteRole(serviceName, roleName string) error {
	if err := repo.Transact(func(tx Tx) error {
		if _, err := tx.Exec("delete from roles where service_name=? and name=?", serviceName, roleName); err != nil {
			return fmt.Errorf("delete from roles: %v", err)
		}

		return nil
	}); err != nil {
		return fmt.Errorf("transaction: %v", err)
	}

	return nil
}

// select * from service_metrics where service_name=${serviceName} and name=${metricName}
func (repo *ServiceRepository) ExistsMetric(serviceName, metricName string) bool {
	rows, err := repo.Query("select * from service_metric_values where service_name=? and name=? limit 1", serviceName, metricName)
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
func (repo *ServiceRepository) MetricNames(serviceName string) (*domain.ServiceMetricValueNames, error) {
	names := make([]string, 0)
	if err := repo.Transact(func(tx Tx) error {
		rows, err := repo.Query("select distinct name from service_metric_values where service_name=?", serviceName)
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
func (repo *ServiceRepository) MetricValues(serviceName, metricName string, from, to int64) (*domain.ServiceMetricValues, error) {
	return &domain.ServiceMetricValues{}, nil
}

// insert into service_metrics values(${serviceName}, ${name}, ${time}, ${value})
func (repo *ServiceRepository) SaveMetricValues(serviceName string, values []domain.ServiceMetricValue) (*domain.Success, error) {
	return &domain.Success{Success: true}, nil
}

// select * from service_metadata where service_name=${serviceName} and namespace=${namespace} limit=1
func (repo *ServiceRepository) ExistsMetadata(serviceName, namespace string) bool {
	return true
}

// select namespacee from service_metadata where service_name=${serviceName}
func (repo *ServiceRepository) MetadataList(serviceName string) (*domain.ServiceMetadataList, error) {
	return nil, nil
}

// select * from service_metadata where service_name=${serviceName} and namespace=${namespace}
func (repo *ServiceRepository) Metadata(serviceName, namespace string) (interface{}, error) {
	return nil, nil
}

// insert into service_metadata values(${serviceName}, ${namespace}, ${metadata})
func (repo *ServiceRepository) SaveMetadata(serviceName, namespace string, metadata interface{}) (*domain.Success, error) {
	return nil, nil
}

// delete from service_metadata where service_name=${serviceName} and namespace=${namespace}
func (repo *ServiceRepository) DeleteMetadata(serviceName, namespace string) (*domain.Success, error) {
	return nil, nil
}

// select * from role_metadata where service_name=${serviceName} and role_name=${roleName} and namespace=${namespace} limit=1
func (repo *ServiceRepository) ExistsRoleMetadata(serviceName, roleName, namespace string) bool {
	return true
}

// select namespace from role_metadata where service_name=${serviceName} and role_name=${roleName}
func (repo *ServiceRepository) RoleMetadataList(serviceName, roleName string) (*domain.RoleMetadataList, error) {
	return nil, nil
}

// select * from role_metadata where service_name=${serviceName} and role_name=${roleName} and namespace=${namespace}
func (repo *ServiceRepository) RoleMetadata(serviceName, roleName, namespace string) (interface{}, error) {
	return nil, nil
}

// insert into role_metadata values(${serviveName}, ${roleName}, ${namespace}, ${metadata})
func (repo *ServiceRepository) SaveRoleMetadata(serviceName, roleName, namespace string, metadata interface{}) (*domain.Success, error) {
	return nil, nil
}

// delete from role_metadata where service_name=${serviceName} and role_name=${roleName} and namespace=${namespace}
func (repo *ServiceRepository) DeleteRoleMetadata(serviceName, roleName, namespace string) (*domain.Success, error) {
	return nil, nil
}
