package database

import (
	"encoding/json"
	"fmt"

	"github.com/itsubaki/mackerel-api/pkg/domain"
)

type ServiceRepository struct {
	SQLHandler
}

func NewServiceRepository(handler SQLHandler) *ServiceRepository {
	if err := handler.Transact(func(tx Tx) error {
		if _, err := tx.Exec(
			`
			create table if not exists services (
				name varchar(128) not null primary key,
				memo varchar(128),
				roles text
			)
			`,
		); err != nil {
			return err
		}

		return nil
	}); err != nil {
		panic(err)
	}

	return &ServiceRepository{
		SQLHandler: handler,
	}
}

// select * from services
func (repo *ServiceRepository) List() (*domain.Services, error) {
	services := make([]domain.Service, 0)

	if err := repo.Transact(func(tx Tx) error {
		rows, err := tx.Query("select * from services")
		if err != nil {
			return err
		}
		defer rows.Close()

		for rows.Next() {
			var service domain.Service
			var roles string
			if err := rows.Scan(
				&service.Name,
				&service.Memo,
				&roles,
			); err != nil {
				return err
			}

			if err := json.Unmarshal([]byte(roles), &service.Roles); err != nil {
				return err
			}

			services = append(services, service)
		}

		return nil
	}); err != nil {
		return nil, err
	}

	return &domain.Services{Services: services}, nil
}

// insert into services values()
func (repo *ServiceRepository) Save(s *domain.Service) error {
	roles, err := json.Marshal(s.Roles)
	if err != nil {
		return err
	}

	return repo.Transact(func(tx Tx) error {
		if _, err := tx.Exec(
			"insert into services values (?, ?, ?)",
			s.Name,
			s.Memo,
			roles,
		); err != nil {
			return err
		}

		return nil
	})
}

// select * from services where service_name=${serviceName}
func (repo *ServiceRepository) Service(serviceName string) (*domain.Service, error) {
	var service domain.Service

	if err := repo.Transact(func(tx Tx) error {
		row := tx.QueryRow("select * from services where name=?", serviceName)
		var roles string
		if err := row.Scan(
			&service.Name,
			&service.Memo,
			&roles,
		); err != nil {
			return err
		}

		if err := json.Unmarshal([]byte(roles), &service.Roles); err != nil {
			return err
		}

		return nil
	}); err != nil {
		return nil, err
	}

	return &service, nil
}

// select * from services where service_name=${serviceName} limit=1
func (repo *ServiceRepository) Exists(serviceName string) bool {
	rows, err := repo.Query("select * from services where name=? limit 1", serviceName)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	if rows.Next() {
		return true
	}

	return false
}

// delete from services where service_name=${serviceName}
func (repo *ServiceRepository) Delete(serviceName string) error {
	if err := repo.Transact(func(tx Tx) error {
		if _, err := tx.Exec("delete from services where name=?", serviceName); err != nil {
			return err
		}

		return nil
	}); err != nil {
		return err
	}

	return nil
}

// select * from service_roles where service_name=${serviceName}
func (repo *ServiceRepository) RoleList(serviceName string) (*domain.Roles, error) {
	return nil, nil
}

// insert into service_roles values(${serviceName}, ${roleName}, ${Memo})
func (repo *ServiceRepository) SaveRole(serviceName string, r *domain.Role) error {
	return nil
}

// select * from service_roles where service_name=${serviceName} and role_name=${roleName}
func (repo *ServiceRepository) ExistsRole(serviceName, roleName string) bool {
	return false
}

// select * from service_roles where service_name=${serviceName} and role_name=${roleName}
func (repo *ServiceRepository) Role(serviceName, roleName string) (*domain.Role, error) {
	return nil, fmt.Errorf("role not found")
}

// delete from service_roles where service_name=${serviceName} and role_name=${roleName}
func (repo *ServiceRepository) DeleteRole(serviceName, roleName string) error {
	return nil
}

// select * from service_metrics where service_name=${serviceName} and name=${metricName}
func (repo *ServiceRepository) ExistsMetric(serviceName, metricName string) bool {
	return true
}

// select distinct name from service_metrics where service_name=${serviceName}
func (repo *ServiceRepository) MetricNames(serviceName string) (*domain.ServiceMetricValueNames, error) {
	return &domain.ServiceMetricValueNames{}, nil
}

// select * from service_metrics where service_name=${serviceName} and name=${metricName}  and ${from} < from and to < ${to}
func (repo *ServiceRepository) MetricValues(serviceName, metricName string, from, to int64) (*domain.ServiceMetricValues, error) {
	return &domain.ServiceMetricValues{}, nil
}

// insert into service_metrics values(${serviceName}, ${name}, ${time}, ${value})
func (repo *ServiceRepository) SaveMetricValues(serviceName string, values []domain.ServiceMetricValue) (*domain.Success, error) {
	return &domain.Success{Success: true}, nil
}

// select * from service_metadata where service_name=${serviceName} and namespace=${namespace} limit=1
func (repo *ServiceRepository) ExistsMetadata(serviceName, namespace string) bool {
	return true
}

// select namespacee from service_metadata where service_name=${serviceName}
func (repo *ServiceRepository) MetadataList(serviceName string) (*domain.ServiceMetadataList, error) {
	return nil, nil
}

// select * from service_metadata where service_name=${serviceName} and namespace=${namespace}
func (repo *ServiceRepository) Metadata(serviceName, namespace string) (interface{}, error) {
	return nil, nil
}

// insert into service_metadata values(${serviceName}, ${namespace}, ${metadata})
func (repo *ServiceRepository) SaveMetadata(serviceName, namespace string, metadata interface{}) (*domain.Success, error) {
	return nil, nil
}

// delete from service_metadata where service_name=${serviceName} and namespace=${namespace}
func (repo *ServiceRepository) DeleteMetadata(serviceName, namespace string) (*domain.Success, error) {
	return nil, nil
}

// select * from role_metadata where service_name=${serviceName} and role_name=${roleName} and namespace=${namespace} limit=1
func (repo *ServiceRepository) ExistsRoleMetadata(serviceName, roleName, namespace string) bool {
	return true
}

// select namespace from role_metadata where service_name=${serviceName} and role_name=${roleName}
func (repo *ServiceRepository) RoleMetadataList(serviceName, roleName string) (*domain.RoleMetadataList, error) {
	return nil, nil
}

// select * from role_metadata where service_name=${serviceName} and role_name=${roleName} and namespace=${namespace}
func (repo *ServiceRepository) RoleMetadata(serviceName, roleName, namespace string) (interface{}, error) {
	return nil, nil
}

// insert into role_metadata values(${serviveName}, ${roleName}, ${namespace}, ${metadata})
func (repo *ServiceRepository) SaveRoleMetadata(serviceName, roleName, namespace string, metadata interface{}) (*domain.Success, error) {
	return nil, nil
}

// delete from role_metadata where service_name=${serviceName} and role_name=${roleName} and namespace=${namespace}
func (repo *ServiceRepository) DeleteRoleMetadata(serviceName, roleName, namespace string) (*domain.Success, error) {
	return nil, nil
}
