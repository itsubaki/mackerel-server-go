package database

import (
	"encoding/json"
	"fmt"

	"github.com/itsubaki/mackerel-api/pkg/domain"
)

type GraphRepository struct {
	SQLHandler
}

func NewGraphRepository(handler SQLHandler) *GraphRepository {
	if err := handler.Transact(func(tx Tx) error {
		if _, err := tx.Exec(
			`
			create table if not exists graph_defs (
				org          varchar(64) not null,
				name         varchar(64) not null,
				display_name varchar(64),
				unit         varchar(64),
				metrics      text,
				primary key(org, name)
			)
			`,
		); err != nil {
			return fmt.Errorf("create table graph_defs: %v", err)
		}

		if _, err := tx.Exec(
			`
			create table if not exists graph_annotations (
				org         varchar(64) not null,
				id          varchar(16) not null primary key,
				title       varchar(64) not null,
				description varchar(64),
				time_from   bigint,
				time_to     bigint,
				service     varchar(128) not null,
				roles       text
			)
			`,
		); err != nil {
			return fmt.Errorf("create table graph_annotations: %v", err)
		}

		return nil
	}); err != nil {
		panic(fmt.Errorf("transaction: %v", err))
	}

	return &GraphRepository{
		SQLHandler: handler,
	}
}

func (repo *GraphRepository) SaveDef(org string, g []domain.GraphDef) (*domain.Success, error) {
	if err := repo.Transact(func(tx Tx) error {
		for i := range g {
			metrics, err := json.Marshal(g[i].Metrics)
			if err != nil {
				return fmt.Errorf("marshal host.Roles: %v", err)
			}

			if _, err := tx.Exec(
				`
				insert into graph_defs (
					org,
					name,
					display_name,
					unit,
					metrics
				)
				values (?, ?, ?, ?, ?)
				on duplicate key update
					display_name = values(display_name),
					unit = values(unit),
					metrics = values(metrics)
				`,
				org,
				g[i].Name,
				g[i].DisplayName,
				g[i].Unit,
				string(metrics),
			); err != nil {
				return fmt.Errorf("insert into graph_defs: %v", err)
			}
		}

		return nil
	}); err != nil {
		return &domain.Success{Success: false}, nil
	}

	return &domain.Success{Success: true}, nil
}
