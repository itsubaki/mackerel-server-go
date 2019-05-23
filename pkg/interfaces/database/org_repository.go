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
			create table if not exists xapikey (
				x_api_key varchar(45)  not null primary key,
				name      varchar(16)  not null,
				xread     boolean not  null default 1,
				xwrite    boolean not  null default 1,
				org       varchar(128) not null
			)
			`,
		); err != nil {
			return fmt.Errorf("create table orgs: %v", err)
		}

		if _, err := tx.Exec(
			`
			insert into xapikey (
				x_api_key,
				name,
				xread,
				xwrite,
				org
			) values (?, ?, ?, ?, ?)
			on duplicate key update
				name   = values(name),
				xread  = values(xread),
				xwrite = values(xwrite),
				org    = values(org)
			`,
			"2684d06cfedbee8499f326037bb6fb7e8c22e73b16bb",
			"default",
			1,
			1,
			"default",
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

func (repo *OrgRepository) Org() (*domain.Org, error) {
	var name string
	if err := repo.Transact(func(tx Tx) error {
		row := tx.QueryRow("select * from orgs")
		if err := row.Scan(
			&name,
		); err != nil {
			return fmt.Errorf("select * from orgs: %v", err)
		}

		return nil
	}); err != nil {
		return nil, fmt.Errorf("transaction: %v", err)
	}

	return &domain.Org{Name: name}, nil
}
