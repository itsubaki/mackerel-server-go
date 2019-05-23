package database

import (
	"fmt"
	"log"

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

		key := domain.NewXAPIKey("default", "default", true)
		if _, err := tx.Exec(`
			insert into xapikey values (?, ?, ?, ?, ?) 
		`,
			key.XAPIKey,
			key.Name,
			key.Read,
			key.Write,
			key.Org,
		); err != nil {
			return fmt.Errorf("insert into xapikey: %v", err)
		}
		log.Printf("%#v\n", key)

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
