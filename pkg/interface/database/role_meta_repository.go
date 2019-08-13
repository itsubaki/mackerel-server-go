package database

import (
	"encoding/json"
	"fmt"

	"github.com/itsubaki/mackerel-api/pkg/domain"
)

type RoleMetaRepository struct {
	SQLHandler
}

func NewRoleMetaRepository(handler SQLHandler) *RoleMetaRepository {
	if err := handler.Transact(func(tx Tx) error {
		if _, err := tx.Exec(
			`
			create table if not exists role_meta (
				org_id       varchar(16)  not null,
				service_name varchar(16)  not null,
				role_name    varchar(16)  not null,
				namespace    varchar(128) not null,
				meta         text,
				primary key(org_id, role_name, namespace)
			)
			`,
		); err != nil {
			return fmt.Errorf("create table role_meta: %v", err)
		}

		return nil
	}); err != nil {
		panic(fmt.Errorf("transaction: %v", err))
	}

	return &RoleMetaRepository{
		SQLHandler: handler,
	}
}

// select * from role_metadata where service_name=${serviceName} and role_name=${roleName} and namespace=${namespace} limit=1
func (repo *RoleMetaRepository) Exists(orgID, serviceName, roleName, namespace string) bool {
	rows, err := repo.Query(
		"select 1 from role_meta where org_id=? and service_name=? and role_name=? and namespace=?",
		orgID,
		serviceName,
		roleName,
		namespace,
	)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	if rows.Next() {
		return true
	}

	return false
}

// select namespace from role_metadata where service_name=${serviceName} and role_name=${roleName}
func (repo *RoleMetaRepository) List(orgID, serviceName, roleName string) (*domain.RoleMetadataList, error) {
	values := make([]domain.RoleMetadata, 0)
	if err := repo.Transact(func(tx Tx) error {
		rows, err := tx.Query("select namespace from role_meta where org_id=? and service_name=? and role_name=? ", orgID, serviceName, roleName)
		if err != nil {
			return fmt.Errorf("select namespace from role_meta: %v", err)
		}
		defer rows.Close()

		for rows.Next() {
			var namespace string
			if err := rows.Scan(
				&namespace,
			); err != nil {
				return fmt.Errorf("scan: %v", err)
			}

			values = append(values, domain.RoleMetadata{Namespace: namespace})
		}

		return nil
	}); err != nil {
		return nil, fmt.Errorf("transaction: %v", err)
	}

	return &domain.RoleMetadataList{Metadata: values}, nil
}

// select * from role_metadata where service_name=${serviceName} and role_name=${roleName} and namespace=${namespace}
func (repo *RoleMetaRepository) Metadata(orgID, serviceName, roleName, namespace string) (interface{}, error) {
	var meta string
	if err := repo.Transact(func(tx Tx) error {
		row := tx.QueryRow("select meta from role_meta where org_id=? and service_name=? and role_name=? and namespace=?", orgID, serviceName, roleName, namespace)
		if err := row.Scan(&meta); err != nil {
			return fmt.Errorf("scan: %v", err)
		}

		return nil
	}); err != nil {
		return nil, fmt.Errorf("transaction: %v", err)
	}

	var out interface{}
	if err := json.Unmarshal([]byte(meta), &out); err != nil {
		return nil, fmt.Errorf("unmarshal: %v", err)
	}

	return out, nil
}

// insert into role_metadata values(${serviceName}, ${roleName}, ${namespace}, ${metadata})
func (repo *RoleMetaRepository) Save(orgID, serviceName, roleName, namespace string, metadata interface{}) (*domain.Success, error) {
	meta, err := json.Marshal(metadata)
	if err != nil {
		return &domain.Success{Success: false}, fmt.Errorf("marshal: %v", err)
	}

	if err := repo.Transact(func(tx Tx) error {
		if _, err := tx.Exec(
			"insert into role_meta values(?, ?, ?, ?, ?)",
			orgID,
			serviceName,
			roleName,
			namespace,
			string(meta),
		); err != nil {
			return fmt.Errorf("insert into role_meta: %v", err)
		}

		return nil
	}); err != nil {
		return &domain.Success{Success: false}, fmt.Errorf("transaction: %v", err)
	}

	return &domain.Success{Success: true}, nil
}

// delete from role_metadata where service_name=${serviceName} and role_name=${roleName} and namespace=${namespace}
func (repo *RoleMetaRepository) Delete(orgID, serviceName, roleName, namespace string) (*domain.Success, error) {
	if err := repo.Transact(func(tx Tx) error {
		if _, err := tx.Exec(
			"delete from role_meta where org_id=? and service_name=? and role_name=? and namespace=?",
			orgID,
			serviceName,
			roleName,
			namespace,
		); err != nil {
			return fmt.Errorf("delete from role_meta: %v", err)
		}

		return nil
	}); err != nil {
		return &domain.Success{Success: false}, fmt.Errorf("transaction: %v", err)
	}

	return &domain.Success{Success: true}, nil
}
