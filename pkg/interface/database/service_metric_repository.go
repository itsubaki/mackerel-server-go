package database

import (
	"fmt"

	"github.com/itsubaki/mackerel-server-go/pkg/domain"
	"github.com/jinzhu/gorm"
)

type ServiceMetricRepository struct {
	DB *gorm.DB
}

type ServiceMetricValue struct {
	OrgID       string  `gorm:"column:org_id;       type:varchar(16);  not null; primary_key"`
	ServiceName string  `gorm:"column:service_name; type:varchar(16);  not null; primary_key"`
	Name        string  `gorm:"column:name;         type:varchar(128); not null; primary_key"`
	Time        int64   `gorm:"column:time;         type:bigint;       not null; primary_key"`
	Value       float64 `gorm:"column:value;        type:double;       not null;"`
}

func NewServiceMetricRepository(handler SQLHandler) *ServiceMetricRepository {
	db, err := gorm.Open(handler.Dialect(), handler.Raw())
	if err != nil {
		panic(err)
	}
	db.LogMode(handler.IsDebugging())

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
		DB: db,
	}
}

func (repo *ServiceMetricRepository) Exists(orgID, serviceName, metricName string) bool {
	if repo.DB.Where(&ServiceMetricValue{OrgID: orgID, ServiceName: serviceName, Name: metricName}).First(&ServiceMetricValue{}).RecordNotFound() {
		return false
	}

	return true
}

func (repo *ServiceMetricRepository) Names(orgID, serviceName string) (*domain.ServiceMetricValueNames, error) {
	result := make([]ServiceMetricValue, 0)
	if err := repo.DB.Model(&ServiceMetricValue{}).Where(&ServiceMetricValue{OrgID: orgID, ServiceName: serviceName}).Select("distinct(name)").Find(&result).Error; err != nil {
		return nil, fmt.Errorf("select distinct name from service_metric_values: %v", err)
	}

	names := make([]string, 0)
	for _, r := range result {
		names = append(names, r.Name)
	}

	return &domain.ServiceMetricValueNames{Names: names}, nil
}

func (repo *ServiceMetricRepository) Values(orgID, serviceName, metricName string, from, to int64) (*domain.ServiceMetricValues, error) {
	result := make([]ServiceMetricValue, 0)
	if err := repo.DB.Where(&ServiceMetricValue{OrgID: orgID, ServiceName: serviceName, Name: metricName}).Where("? < time and time < ?", from, to).Find(&result).Error; err != nil {
		return nil, fmt.Errorf("select time, value from service_metric_values: %v", err)
	}

	out := make([]domain.ServiceMetricValue, 0)
	for _, r := range result {
		out = append(out, domain.ServiceMetricValue{
			OrgID:       orgID,
			ServiceName: serviceName,
			Name:        metricName,
			Time:        r.Time,
			Value:       r.Value,
		})
	}

	return &domain.ServiceMetricValues{Metrics: out}, nil
}

func (repo *ServiceMetricRepository) Save(orgID, serviceName string, values []domain.ServiceMetricValue) (*domain.Success, error) {
	if err := repo.DB.Transaction(func(tx *gorm.DB) error {
		for i := range values {
			if err := tx.Create(&ServiceMetricValue{
				OrgID:       orgID,
				ServiceName: serviceName,
				Name:        values[i].Name,
				Time:        values[i].Time,
				Value:       values[i].Value,
			}).Error; err != nil {
				return fmt.Errorf("insert into service_metric_values: %v", err)
			}
		}

		return nil
	}); err != nil {
		return &domain.Success{Success: false}, fmt.Errorf("transaction: %v", err)
	}

	return &domain.Success{Success: true}, nil
}
