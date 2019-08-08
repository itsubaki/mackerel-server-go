package database

import (
	"fmt"

	"github.com/itsubaki/mackerel-api/pkg/domain"
)

type APIKeyRepository struct {
	SQLHandler
}

func NewAPIKeyRepository(handler SQLHandler) *APIKeyRepository {
	if err := handler.Transact(func(tx Tx) error {
		if _, err := tx.Exec(
			`
			create table if not exists apikeys (
				org_id  varchar(16) not null,
				name    varchar(16) not null,
				api_key varchar(45) not null primary key,
				xread   boolean     not null default 1,
				xwrite  boolean     not null default 1
			)
			`,
		); err != nil {
			return fmt.Errorf("create table apikeys: %v", err)
		}

		if _, err := tx.Exec(
			`
			insert into apikeys (
				org_id,
				name,
				api_key,
				xread,
				xwrite
			) values (?, ?, ?, ?, ?)
			on duplicate key update
				org_id  = values(org_id),
				name    = values(name),
				api_key = values(api_key),
				xread   = values(xread),
				xwrite  = values(xwrite)
			`,
			"4b825dc642c",
			"default",
			"2684d06cfedbee8499f326037bb6fb7e8c22e73b16bb",
			1,
			1,
		); err != nil {
			return fmt.Errorf("insert into apikeys: %v", err)
		}

		return nil
	}); err != nil {
		panic(fmt.Errorf("transaction: %v", err))
	}

	return &APIKeyRepository{
		SQLHandler: handler,
	}
}

func (repo *APIKeyRepository) Save(orgID, name string, write bool) (*domain.APIKey, error) {
	apikey := domain.NewAPIKey()
	if err := repo.Transact(func(tx Tx) error {
		if _, err := tx.Exec(
			`
			insert into apikeys (
				org_id,
				name,
				api_key,
				xread,
				xwrite
			) values (?, ?, ?, ?, ?)
		`,
			orgID,
			name,
			apikey,
			true,
			write,
		); err != nil {
			return fmt.Errorf("insert into apikeys: %v", err)
		}

		return nil
	}); err != nil {
		return nil, fmt.Errorf("transaction: %v", err)
	}

	return &domain.APIKey{
		OrgID:  orgID,
		Name:   name,
		APIKey: apikey,
		Read:   true,
		Write:  write,
	}, nil
}

func (repo *APIKeyRepository) APIKey(xapikey string) (*domain.APIKey, error) {
	var key domain.APIKey
	if err := repo.Transact(func(tx Tx) error {
		row := tx.QueryRow(`select * from apikeys where api_key=?`, xapikey)
		if err := row.Scan(
			&key.OrgID,
			&key.Name,
			&key.APIKey,
			&key.Read,
			&key.Write,
		); err != nil {
			return fmt.Errorf("select * from apikeys: %v", err)
		}

		return nil
	}); err != nil {
		return nil, fmt.Errorf("transaction: %v", err)
	}

	return &key, nil
}
