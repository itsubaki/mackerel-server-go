package database

import (
	"fmt"

	"github.com/itsubaki/mackerel-api/pkg/domain"
)

type RoleRepository struct {
	SQLHandler
}

func NewRoleRepository(handler SQLHandler) *RoleRepository {
	if err := handler.Transact(func(tx Tx) error {
		if _, err := tx.Exec(
			`
			create table if not exists roles (
				org_id       varchar(16)  not null,
				service_name varchar(128) not null,
				name         varchar(128) not null,
				memo         varchar(128) not null default '',
				primary key(org_id, service_name, name),
				foreign key fk_service_name (org_id, service_name) references services(org_id, name) on delete cascade on update cascade
			)
			`,
		); err != nil {
			return fmt.Errorf("create table roles: %v", err)
		}

		return nil
	}); err != nil {
		panic(fmt.Errorf("transaction: %v", err))
	}

	return &RoleRepository{
		SQLHandler: handler,
	}
}

func (repo *RoleRepository) ListWith(orgID string) (map[string][]string, error) {
	roles := make(map[string][]string)

	if err := repo.Transact(func(tx Tx) error {
		rows, err := tx.Query("select service_name, name from roles where org_id=?", orgID)
		if err != nil {
			return fmt.Errorf("select service_name, name from roles: %v", err)
		}
		defer rows.Close()

		for rows.Next() {
			var service, role string
			if err := rows.Scan(
				&service,
				&role,
			); err != nil {
				return fmt.Errorf("scan roles: %v", err)
			}

			if _, ok := roles[service]; !ok {
				roles[service] = make([]string, 0)
			}

			roles[service] = append(roles[service], role)
		}
		return nil
	}); err != nil {
		return nil, fmt.Errorf("transaction: %v", err)
	}

	return roles, nil
}

// select * from service_roles where service_name=${serviceName}
func (repo *RoleRepository) List(orgID, serviceName string) (*domain.Roles, error) {
	var roles []domain.Role
	if err := repo.Transact(func(tx Tx) error {
		rows, err := tx.Query("select name, memo from roles where org_id=? and service_name=?", orgID, serviceName)
		if err != nil {
			return fmt.Errorf("select name, memo from roles: %v", err)
		}
		defer rows.Close()

		for rows.Next() {
			var role domain.Role
			if err := rows.Scan(
				&role.Name,
				&role.Memo,
			); err != nil {
				return fmt.Errorf("scan: %v", err)
			}

			roles = append(roles, role)
		}

		return nil
	}); err != nil {
		return nil, fmt.Errorf("transaction: %v", err)
	}

	return &domain.Roles{Roles: roles}, nil
}

// select * from service_roles where service_name=${serviceName} and role_name=${roleName}
func (repo *RoleRepository) Role(orgID, serviceName, roleName string) (*domain.Role, error) {
	role := domain.Role{
		ServiceName: serviceName,
		Name:        roleName,
	}

	if err := repo.Transact(func(tx Tx) error {
		row := repo.QueryRow("select memo from roles where org_id=? and service_name=? and name=?", orgID, serviceName, roleName)
		if err := row.Scan(
			&role.Memo,
		); err != nil {
			return fmt.Errorf("scan: %v", err)
		}

		return nil
	}); err != nil {
		return nil, fmt.Errorf("transaction: %v", err)
	}

	return &role, nil
}

// insert into service_roles values(${serviceName}, ${roleName}, ${Memo})
func (repo *RoleRepository) Save(orgID, serviceName string, r *domain.Role) error {
	if err := repo.Transact(func(tx Tx) error {
		if _, err := tx.Exec(
			`
			insert into roles (
				org_id,
				service_name,
				name,
				memo
			)
			values (?, ?, ?, ?)
			on duplicate key update
				service_name = values(service_name),
				name = values(name),
				memo = values(memo)
			`,
			orgID,
			serviceName,
			r.Name,
			r.Memo,
		); err != nil {
			return fmt.Errorf("insert into roles: %v", err)
		}

		return nil
	}); err != nil {
		return fmt.Errorf("transaction: %v", err)
	}

	return nil
}

// select * from service_roles where service_name=${serviceName} and role_name=${roleName}
func (repo *RoleRepository) Exists(orgID, serviceName, roleName string) bool {
	rows, err := repo.Query("select 1 from roles where org_id=? and service_name=? and name=? limit 1", orgID, serviceName, roleName)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	if rows.Next() {
		return true
	}

	return false
}

// delete from service_roles where service_name=${serviceName} and role_name=${roleName}
func (repo *RoleRepository) Delete(orgID, serviceName, roleName string) error {
	if err := repo.Transact(func(tx Tx) error {
		if _, err := tx.Exec("delete from roles where org_id=? and service_name=? and name=?", orgID, serviceName, roleName); err != nil {
			return fmt.Errorf("delete from roles: %v", err)
		}

		return nil
	}); err != nil {
		return fmt.Errorf("transaction: %v", err)
	}

	return nil
}
