package database

import (
	"encoding/json"
	"fmt"

	"github.com/itsubaki/mackerel-server-go/pkg/domain"
)

type ServiceMetaRepository struct {
	SQLHandler
}

func NewServiceMetaRepository(handler SQLHandler) *ServiceMetaRepository {
	if err := handler.Transact(func(tx Tx) error {
		if _, err := tx.Exec(
			`
			create table if not exists service_meta (
				org_id       varchar(16)  not null,
				service_name varchar(16)  not null,
				namespace    varchar(128) not null,
				meta         text,
				primary key(org_id, service_name, namespace)
			)
			`,
		); err != nil {
			return fmt.Errorf("create table service_meta: %v", err)
		}

		return nil
	}); err != nil {
		panic(fmt.Errorf("transaction: %v", err))
	}

	return &ServiceMetaRepository{
		SQLHandler: handler,
	}
}

// select * from service_metadata where service_name=${serviceName} and namespace=${namespace} limit=1
func (repo *ServiceMetaRepository) Exists(orgID, serviceName, namespace string) bool {
	rows, err := repo.Query("select 1 from service_meta where org_id=? and service_name=? and namespace=?", orgID, serviceName, namespace)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	if rows.Next() {
		return true
	}

	return false
}

// select namespace from service_metadata where service_name=${serviceName}
func (repo *ServiceMetaRepository) List(orgID, serviceName string) (*domain.ServiceMetadataList, error) {
	values := make([]domain.ServiceMetadata, 0)
	if err := repo.Transact(func(tx Tx) error {
		rows, err := tx.Query("select namespace from service_meta where org_id=? and service_name=?", orgID, serviceName)
		if err != nil {
			return fmt.Errorf("select namespace from service_meta: %v", err)
		}
		defer rows.Close()

		for rows.Next() {
			var namespace string
			if err := rows.Scan(
				&namespace,
			); err != nil {
				return fmt.Errorf("scan: %v", err)
			}

			values = append(values, domain.ServiceMetadata{Namespace: namespace})
		}

		return nil
	}); err != nil {
		return nil, fmt.Errorf("transaction: %v", err)
	}

	return &domain.ServiceMetadataList{Metadata: values}, nil
}

// select * from service_metadata where service_name=${serviceName} and namespace=${namespace}
func (repo *ServiceMetaRepository) Metadata(orgID, serviceName, namespace string) (interface{}, error) {
	var meta string
	if err := repo.Transact(func(tx Tx) error {
		row := tx.QueryRow("select meta from service_meta where org_id=? and service_name=? and namespace=?", orgID, serviceName, namespace)
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

// insert into service_metadata values(${serviceName}, ${namespace}, ${metadata})
func (repo *ServiceMetaRepository) Save(orgID, serviceName, namespace string, metadata interface{}) (*domain.Success, error) {
	meta, err := json.Marshal(metadata)
	if err != nil {
		return &domain.Success{Success: false}, fmt.Errorf("marshal: %v", err)
	}

	if err := repo.Transact(func(tx Tx) error {
		if _, err := tx.Exec(
			"insert into service_meta values(?, ?, ?, ?)",
			orgID,
			serviceName,
			namespace,
			string(meta),
		); err != nil {
			return fmt.Errorf("insert into service_meta: %v", err)
		}

		return nil
	}); err != nil {
		return &domain.Success{Success: false}, fmt.Errorf("transaction: %v", err)
	}

	return &domain.Success{Success: true}, nil
}

// delete from service_metadata where service_name=${serviceName} and namespace=${namespace}
func (repo *ServiceMetaRepository) Delete(orgID, serviceName, namespace string) (*domain.Success, error) {
	if err := repo.Transact(func(tx Tx) error {
		if _, err := tx.Exec(
			"delete from service_meta where org_id=? and service_name=? and namespace=?",
			orgID,
			serviceName,
			namespace,
		); err != nil {
			return fmt.Errorf("delete from service_meta: %v", err)
		}

		return nil
	}); err != nil {
		return &domain.Success{Success: false}, fmt.Errorf("transaction: %v", err)
	}

	return &domain.Success{Success: true}, nil
}
