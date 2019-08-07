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
				id   varchar(16) not null primary key,
				name varchar(16) not null unique
			)
			`,
		); err != nil {
			return fmt.Errorf("create table orgs: %v", err)
		}

		if _, err := tx.Exec(
			`
			insert into orgs (
				id,
				name
			) values (?, ?)
			on duplicate key update
				name = values(name)
			`,
			"4b825dc642c",
			"mackerel",
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

func (repo *OrgRepository) Org(orgID string) (*domain.Org, error) {
	var name string
	if err := repo.Transact(func(tx Tx) error {
		row := tx.QueryRow(`select name from orgs where id=?`, orgID)
		if err := row.Scan(
			&name,
		); err != nil {
			return fmt.Errorf("select name from orgs where id=?: %v", err)
		}

		return nil
	}); err != nil {
		return nil, fmt.Errorf("transaction: %v", err)
	}
	return &domain.Org{ID: orgID, Name: name}, nil
}
