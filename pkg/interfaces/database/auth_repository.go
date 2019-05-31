package database

import (
	"fmt"

	"github.com/itsubaki/mackerel-api/pkg/domain"
)

type AuthRepository struct {
	SQLHandler
}

func NewAuthRepository(handler SQLHandler) *AuthRepository {
	if err := handler.Transact(func(tx Tx) error {
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
			"4b825dc642c",
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

	return &AuthRepository{
		SQLHandler: handler,
	}
}

func (repo *AuthRepository) Save(orgID, name string, write bool) (*domain.XAPIKey, error) {
	xapikey := domain.NewXAPIKey()
	if err := repo.Transact(func(tx Tx) error {
		if _, err := tx.Exec(
			`
			insert into xapikey (
				org_id,
				name,
				x_api_key,
				xread,
				xwrite
			) values (?, ?, ?, ?, ?)
		`,
			orgID,
			name,
			xapikey,
			true,
			write,
		); err != nil {
			return fmt.Errorf("insert into xapikey: %v", err)
		}

		return nil
	}); err != nil {
		return nil, fmt.Errorf("transaction: %v", err)
	}

	return &domain.XAPIKey{
		OrgID:   orgID,
		Name:    name,
		XAPIKey: xapikey,
		Read:    true,
		Write:   write,
	}, nil
}

func (repo *AuthRepository) XAPIKey(xapikey string) (*domain.XAPIKey, error) {
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
