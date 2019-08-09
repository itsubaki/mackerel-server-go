package database

import (
	"encoding/json"
	"fmt"
	"strings"

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
				org_id       varchar(16) not null,
				name         varchar(64) not null,
				display_name varchar(64),
				unit         varchar(64),
				metrics      text,
				primary key(org_id, name)
			)
			`,
		); err != nil {
			return fmt.Errorf("create table graph_defs: %v", err)
		}

		if _, err := tx.Exec(
			`
			create table if not exists graph_annotations (
				org_id      varchar(16) not null,
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

func (repo *GraphRepository) SaveDef(orgID string, g []domain.GraphDef) (*domain.Success, error) {
	if err := repo.Transact(func(tx Tx) error {
		for i := range g {
			metrics, err := json.Marshal(g[i].Metrics)
			if err != nil {
				return fmt.Errorf("marshal host.Roles: %v", err)
			}

			if _, err := tx.Exec(
				`
				insert into graph_defs (
					org_id,
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
				orgID,
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
		return &domain.Success{Success: false}, fmt.Errorf("transaction: %v", err)
	}

	return &domain.Success{Success: true}, nil
}

func (repo *GraphRepository) List(orgID string) (*domain.GraphAnnotations, error) {
	annotations := make([]domain.GraphAnnotation, 0)
	if err := repo.Transact(func(tx Tx) error {
		rows, err := tx.Query("select * from graph_annotations where org_id=?", orgID)
		if err != nil {
			return fmt.Errorf("select * from graph_annotations: %v", err)
		}
		defer rows.Close()

		for rows.Next() {
			var annotation domain.GraphAnnotation
			var roles string
			if err := rows.Scan(
				&annotation.OrgID,
				&annotation.ID,
				&annotation.Title,
				&annotation.Description,
				&annotation.From,
				&annotation.To,
				&annotation.Service,
				&roles,
			); err != nil {
				return fmt.Errorf("scan: %v", err)
			}
			annotation.Roles = strings.Split(roles, ",")

			annotations = append(annotations, annotation)
		}

		return nil
	}); err != nil {
		return nil, fmt.Errorf("transaction: %v", err)
	}

	return &domain.GraphAnnotations{GraphAnnotations: annotations}, nil
}

func (repo *GraphRepository) Save(orgID string, annotation *domain.GraphAnnotation) (*domain.GraphAnnotation, error) {
	if err := repo.Transact(func(tx Tx) error {
		if _, err := tx.Exec(
			`
			insert into graph_annotations (
				org_id,
				id,
				title,
				description,
				time_from,
				time_to,
				service,
				roles
			) values (?, ?, ?, ?, ?, ?, ?, ?)
			`,
			orgID,
			annotation.ID,
			annotation.Title,
			annotation.Description,
			annotation.From,
			annotation.To,
			annotation.Service,
			strings.Join(annotation.Roles, ","),
		); err != nil {
			return fmt.Errorf("insert into graph_annotations: %v", err)
		}

		return nil
	}); err != nil {
		return nil, fmt.Errorf("transaction: %v", err)
	}

	return annotation, nil
}

func (repo *GraphRepository) Update(orgID string, annotation *domain.GraphAnnotation) (*domain.GraphAnnotation, error) {
	if err := repo.Transact(func(tx Tx) error {
		if _, err := tx.Exec(
			`
			update graph_annotations set
				title=?,
				description=?,
				time_from=?,
				time_to=?,
				service=?,
				roles=?
			where org_id=? and id=?
			`,
			annotation.Title,
			annotation.Description,
			annotation.From,
			annotation.To,
			annotation.Service,
			strings.Join(annotation.Roles, ","),
			orgID,
			annotation.ID,
		); err != nil {
			return fmt.Errorf("update graph_annotations: %v", err)
		}

		return nil
	}); err != nil {
		return nil, fmt.Errorf("transaction: %v", err)
	}

	return annotation, nil
}

func (repo *GraphRepository) Delete(orgID, annotationID string) (*domain.GraphAnnotation, error) {
	var annotation domain.GraphAnnotation
	if err := repo.Transact(func(tx Tx) error {
		row := tx.QueryRow("select * from graph_annotations where org_id=? and id=?", orgID, annotationID)

		var roles string
		if err := row.Scan(
			&annotation.OrgID,
			&annotation.ID,
			&annotation.Title,
			&annotation.Description,
			&annotation.From,
			&annotation.To,
			&annotation.Service,
			&roles,
		); err != nil {
			return fmt.Errorf("scan: %v", err)
		}
		annotation.Roles = strings.Split(roles, ",")

		if _, err := tx.Exec("delete from graph_annotations where org_id=? and id=?", orgID, annotationID); err != nil {
			return fmt.Errorf("delete from graph_annotations: %v", err)
		}

		return nil
	}); err != nil {
		return nil, fmt.Errorf("transaction: %v", err)
	}

	return &annotation, nil
}
