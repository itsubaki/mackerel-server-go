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
				id   varchar(64) not null primary key,
				name varchar(16) not null
			)
			`,
		); err != nil {
			return fmt.Errorf("create table orgs: %v", err)
		}

		if _, err := tx.Exec(
			`
			create table if not exists xapikey (
				org_id    varchar(64) not null,
				name      varchar(16) not null,
				x_api_key varchar(45) not null primary key,
				xread     boolean     not null default 1,
				xwrite    boolean     not null default 1
			)
			`,
		); err != nil {
			return fmt.Errorf("create table xapikey: %v", err)
		}

		id := domain.NewOrgID()
		if _, err := tx.Exec(
			`
			insert into orgs (
				id,
				name
			) values (?, ?)
			on duplicate key update
				name = values(name)
			`,
			id,
			"default",
		); err != nil {
			return fmt.Errorf("insert into orgs: %v", err)
		}

		if _, err := tx.Exec(
			`
			insert into xapikey (
				org_id,
				name,
				x_api_key,
				xread,
				xwrite
			) values (?, ?, ?, ?, ?)
			on duplicate key update
				org_id    = values(org_id),
				name      = values(name),
				x_api_key = values(x_api_key),
				xread     = values(xread),
				xwrite    = values(xwrite)
			`,
			id,
			"default",
			"2684d06cfedbee8499f326037bb6fb7e8c22e73b16bb",
			1,
			1,
		); err != nil {
			return fmt.Errorf("insert into xapikey: %v", err)
		}

		return nil
	}); err != nil {
		panic(fmt.Errorf("transaction: %v", err))
	}

	return &OrgRepository{
		handler,
	}
}

func (repo *OrgRepository) Org(orgID string) (*domain.Org, error) {
	var name string
	if err := repo.Transact(func(tx Tx) error {
		row := tx.QueryRow(`select name from orgs where id=?`, orgID)
		if err := row.Scan(
			&name,
		); err != nil {
			return fmt.Errorf("select * from xapikey: %v", err)
		}

		return nil
	}); err != nil {
		return nil, fmt.Errorf("transaction: %v", err)
	}
	return &domain.Org{Name: name}, nil
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
		row := tx.QueryRow(`select * from xapikey where x_api_key=?`, xapikey)
		if err := row.Scan(
			&key.OrgID,
			&key.XAPIKey,
			&key.Name,
			&key.Read,
			&key.Write,
		); err != nil {
			return fmt.Errorf("select * from xapikey: %v", err)
		}

		return nil
	}); err != nil {
		return nil, fmt.Errorf("transaction: %v", err)
	}

	return &key, nil
}
