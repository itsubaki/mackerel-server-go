package database

import (
	"encoding/json"
	"fmt"

	"github.com/itsubaki/mackerel-server-go/pkg/domain"
)

type HostMetaRepository struct {
	SQLHandler
}

func NewHostMetaRepository(handler SQLHandler) *HostMetaRepository {
	if err := handler.Transact(func(tx Tx) error {
		if _, err := tx.Exec(
			`
			create table if not exists host_meta (
				org_id    varchar(16)  not null,
				host_id   varchar(16)  not null,
				namespace varchar(128) not null,
				meta      text,
				primary key(host_id, namespace)
			)
			`,
		); err != nil {
			return fmt.Errorf("create table host_meta: %v", err)
		}

		return nil
	}); err != nil {
		panic(fmt.Errorf("transaction: %v", err))
	}

	return &HostMetaRepository{
		SQLHandler: handler,
	}
}

// select * from host_metadata where host_id=${hostID} and namespace=${namespace} limit=1
func (repo *HostMetaRepository) Exists(orgID, hostID, namespace string) bool {
	rows, err := repo.Query("select 1 from host_meta where org_id=? and host_id=? and namespace=?", orgID, hostID, namespace)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	if rows.Next() {
		return true
	}

	return false
}

// select namespace from host_metadata where host_id=${hostID}
func (repo *HostMetaRepository) List(orgID, hostID string) (*domain.HostMetadataList, error) {
	values := make([]domain.Namespace, 0)
	if err := repo.Transact(func(tx Tx) error {
		rows, err := tx.Query("select namespace from host_meta where org_id=? and host_id=?", orgID, hostID)
		if err != nil {
			return fmt.Errorf("select namespace from host_meta: %v", err)
		}
		defer rows.Close()

		for rows.Next() {
			var namespace string
			if err := rows.Scan(
				&namespace,
			); err != nil {
				return fmt.Errorf("scan: %v", err)
			}

			values = append(values, domain.Namespace{Namespace: namespace})
		}

		return nil
	}); err != nil {
		return nil, fmt.Errorf("transaction: %v", err)
	}

	return &domain.HostMetadataList{Metadata: values}, nil
}

// select * from host_metadata where host_id=${hostID} and namespace=${namespace}
func (repo *HostMetaRepository) Metadata(orgID, hostID, namespace string) (interface{}, error) {
	var meta string
	if err := repo.Transact(func(tx Tx) error {
		row := tx.QueryRow("select meta from host_meta where org_id=? and host_id=? and namespace=?", orgID, hostID, namespace)
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

// insert into host_metadata values(${hostID}, ${namespace}, ${metadata})
func (repo *HostMetaRepository) Save(orgID, hostID, namespace string, metadata interface{}) (*domain.Success, error) {
	meta, err := json.Marshal(metadata)
	if err != nil {
		return &domain.Success{Success: false}, fmt.Errorf("marshal: %v", err)
	}

	if err := repo.Transact(func(tx Tx) error {
		if _, err := tx.Exec(
			"insert into host_meta values(?, ?, ?, ?)",
			orgID,
			hostID,
			namespace,
			string(meta),
		); err != nil {
			return fmt.Errorf("insert into host_meta: %v", err)
		}

		return nil
	}); err != nil {
		return &domain.Success{Success: false}, fmt.Errorf("transaction: %v", err)
	}

	return &domain.Success{Success: true}, nil
}

// delete from host_metadata where host_id=${hostID} and namespace=${namespace}
func (repo *HostMetaRepository) Delete(orgID, hostID, namespace string) (*domain.Success, error) {
	if err := repo.Transact(func(tx Tx) error {
		if _, err := tx.Exec(
			"delete from host_meta where org_id=? and host_id=? and namespace=?",
			orgID,
			hostID,
			namespace,
		); err != nil {
			return fmt.Errorf("delete from host_meta: %v", err)
		}

		return nil
	}); err != nil {
		return &domain.Success{Success: false}, fmt.Errorf("transaction: %v", err)
	}

	return &domain.Success{Success: true}, nil
}
