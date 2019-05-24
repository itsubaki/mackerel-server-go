package database

import (
	"fmt"

	"github.com/itsubaki/mackerel-api/pkg/domain"
)

type OrgRepository struct {
	SQLHandler
}

func NewOrgRepository(handler SQLHandler) *OrgRepository {
	if err := handler.Transact(func(tx Tx) error {
		if _, err := tx.Exec(
			`
			create table if not exists orgs (
				org       varchar(64) not null,
				x_api_key varchar(45) not null primary key,
				name      varchar(16) not null,
				xread     boolean     not null default 1,
				xwrite    boolean     not null default 1
			)
			`,
		); err != nil {
			return fmt.Errorf("create table orgs: %v", err)
		}

		if _, err := tx.Exec(
			`
			insert into orgs (
				org,
				x_api_key,
				name,
				xread,
				xwrite
			) values (?, ?, ?, ?, ?)
			on duplicate key update
				org       = values(org),
				x_api_key = values(x_api_key),
				name      = values(name),
				xread     = values(xread),
				xwrite    = values(xwrite)
			`,
			"default",
			"2684d06cfedbee8499f326037bb6fb7e8c22e73b16bb",
			"default",
			1,
			1,
		); err != nil {
			return fmt.Errorf("insert into orgs: %v", err)
		}

		return nil
	}); err != nil {
		panic(fmt.Errorf("transaction: %v", err))
	}

	return &OrgRepository{
		handler,
	}
}

// mysql> explain select * from orgs where x_api_key='2684d06cfedbee8499f326037bb6fb7e8c22e73b16bb';
// +----+-------------+---------+------------+-------+---------------+---------+---------+-------+------+----------+-------+
// | id | select_type | table   | partitions | type  | possible_keys | key     | key_len | ref   | rows | filtered | Extra |
// +----+-------------+---------+------------+-------+---------------+---------+---------+-------+------+----------+-------+
// |  1 | SIMPLE      | xapikey | NULL       | const | PRIMARY       | PRIMARY | 182     | const |    1 |   100.00 | NULL  |
// +----+-------------+---------+------------+-------+---------------+---------+---------+-------+------+----------+-------+
// 1 row in set, 1 warning (0.00 sec)
func (repo *OrgRepository) XAPIKey(xapikey string) (*domain.XAPIKey, error) {
	var key domain.XAPIKey
	if err := repo.Transact(func(tx Tx) error {
		row := tx.QueryRow(`select * from orgs where x_api_key=?`, xapikey)

		if err := row.Scan(
			&key.Org,
			&key.XAPIKey,
			&key.Name,
			&key.Read,
			&key.Write,
		); err != nil {
			return fmt.Errorf("select * from orgs: %v", err)
		}

		return nil
	}); err != nil {
		return nil, fmt.Errorf("transaction: %v", err)
	}

	return &key, nil
}

func (repo *OrgRepository) Org(org string) (*domain.Org, error) {
	return &domain.Org{Name: org}, nil
}
