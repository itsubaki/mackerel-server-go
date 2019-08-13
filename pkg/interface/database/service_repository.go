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
				org_id varchar(16)  not null,
				name   varchar(128) not null,
				memo   varchar(128) not null default '',
				primary key(org_id, name)
			)
			`,
		); err != nil {
			return fmt.Errorf("create table services: %v", err)
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
func (repo *ServiceRepository) List(orgID string, roles map[string][]string) (*domain.Services, error) {
	services := make([]domain.Service, 0)

	if err := repo.Transact(func(tx Tx) error {
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
