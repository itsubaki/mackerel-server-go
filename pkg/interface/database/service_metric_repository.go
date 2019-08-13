package database

import (
	"fmt"

	"github.com/itsubaki/mackerel-api/pkg/domain"
)

type ServiceMetricRepository struct {
	SQLHandler
}

func NewServiceMetricRepository(handler SQLHandler) *ServiceMetricRepository {
	if err := handler.Transact(func(tx Tx) error {
		if _, err := tx.Exec(
			`
			create table if not exists service_metric_values (
				org_id       varchar(16)  not null,
				service_name varchar(16)  not null,
				name         varchar(128) not null,
				time         bigint not null,
				value        double not null,
				primary key(org_id, service_name, name, time desc)
			)
			`,
		); err != nil {
			return fmt.Errorf("create table service_metric_values: %v", err)
		}

		return nil
	}); err != nil {
		panic(fmt.Errorf("transaction: %v", err))
	}

	return &ServiceMetricRepository{
		SQLHandler: handler,
	}
}

// select * from service_metrics where service_name=${serviceName} and name=${metricName}
func (repo *ServiceMetricRepository) Exists(orgID, serviceName, metricName string) bool {
	rows, err := repo.Query("select 1 from service_metric_values where org_id=? and service_name=? and name=? limit 1", orgID, serviceName, metricName)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	if rows.Next() {
		return true
	}

	return false
}

// select distinct name from service_metrics where service_name=${serviceName}
func (repo *ServiceMetricRepository) Names(orgID, serviceName string) (*domain.ServiceMetricValueNames, error) {
	names := make([]string, 0)
	if err := repo.Transact(func(tx Tx) error {
		rows, err := repo.Query("select distinct name from service_metric_values where org_id=? and service_name=?", orgID, serviceName)
		if err != nil {
			return fmt.Errorf("select distinct name from service_metric_values: %v", err)
		}
		defer rows.Close()

		for rows.Next() {
			var name string
			if err := rows.Scan(&name); err != nil {
				return fmt.Errorf("scan: %v", err)
			}

			names = append(names, name)
		}

		return nil
	}); err != nil {
		return nil, fmt.Errorf("transaction: %v", err)
	}

	return &domain.ServiceMetricValueNames{Names: names}, nil
}

// select * from service_metrics where service_name=${serviceName} and name=${metricName}  and ${from} < from and to < ${to}
func (repo *ServiceMetricRepository) Values(orgID, serviceName, metricName string, from, to int64) (*domain.ServiceMetricValues, error) {
	values := make([]domain.ServiceMetricValue, 0)
	if err := repo.Transact(func(tx Tx) error {
		rows, err := tx.Query("select time, value from service_metric_values where org_id=? and service_name=? and name=? and ? < time and time < ?", orgID, serviceName, metricName, from, to)
		if err != nil {
			return fmt.Errorf("select time, value from service_metric_values: %v", err)
		}
		defer rows.Close()

		for rows.Next() {
			var time int64
			var value float64
			if err := rows.Scan(&time, &value); err != nil {
				return fmt.Errorf("scan: %v", err)
			}

			values = append(values, domain.ServiceMetricValue{
				ServiceName: serviceName,
				Name:        metricName,
				Time:        time,
				Value:       value,
			})
		}

		return nil
	}); err != nil {
		return nil, fmt.Errorf("transaction: %v", err)
	}

	return &domain.ServiceMetricValues{Metrics: values}, nil
}

// insert into service_metrics values(${serviceName}, ${name}, ${time}, ${value})
func (repo *ServiceMetricRepository) Save(orgID, serviceName string, values []domain.ServiceMetricValue) (*domain.Success, error) {
	if err := repo.Transact(func(tx Tx) error {
		for i := range values {
			if _, err := tx.Exec(
				"insert into service_metric_values values(?, ?, ?, ?, ?)",
				orgID,
				serviceName,
				values[i].Name,
				values[i].Time,
				values[i].Value,
			); err != nil {
				return fmt.Errorf("insert into service_metric_values: %v", err)
			}
		}

		return nil
	}); err != nil {
		return &domain.Success{Success: false}, fmt.Errorf("transaction: %v", err)
	}

	return &domain.Success{Success: true}, nil
}
