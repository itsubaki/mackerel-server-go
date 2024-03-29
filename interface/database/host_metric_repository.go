package database

import (
	"errors"
	"fmt"

	"github.com/itsubaki/mackerel-server-go/domain"
	"github.com/itsubaki/mackerel-server-go/usecase"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var _ usecase.HostMetricRepository = (*HostMetricRepository)(nil)

type HostMetricRepository struct {
	DB *gorm.DB
}

type HostMetricValue struct {
	OrgID  string  `gorm:"column:org_id;  type:varchar(16);  not null;"`
	HostID string  `gorm:"column:host_id; type:varchar(16);  not null; primary_key"`
	Name   string  `gorm:"column:name;    type:varchar(128); not null; primary_key"`
	Time   int64   `gorm:"column:time;    type:bigint;       not null; primary_key"`
	Value  float64 `gorm:"column:value;   type:double;       not null;"`
}

type HostMetricValuesLatest struct {
	OrgID  string  `gorm:"column:org_id;  type:varchar(16);  not null;"`
	HostID string  `gorm:"column:host_id; type:varchar(16);  not null; primary_key"`
	Name   string  `gorm:"column:name;    type:varchar(128); not null; primary_key"`
	Value  float64 `gorm:"column:value;   type:double;       not null;"`
}

func (v HostMetricValuesLatest) TableName() string {
	return "host_metric_values_latest"
}

func NewHostMetricRepository(handler SQLHandler) *HostMetricRepository {
	db, err := gorm.Open(mysql.New(mysql.Config{
		Conn: handler.Raw().(gorm.ConnPool),
	}), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	if handler.IsDebugMode() {
		db.Logger.LogMode(4)
	}

	if err := db.AutoMigrate(&HostMetricValuesLatest{}); err != nil {
		panic(fmt.Errorf("auto migrate host_metric_values_latest: %v", err))
	}

	if err := db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Exec(
			`
			create table if not exists host_metric_values (
				org_id  varchar(16)  not null,
				host_id varchar(16)  not null,
				name    varchar(128) not null,
				time    bigint not null,
				value   double not null,
				primary key(host_id, name, time desc)
			)
			`,
		).Error; err != nil {
			return fmt.Errorf("create table host_metric_values: %v", err)
		}

		return nil
	}); err != nil {
		panic(fmt.Errorf("transaction: %v", err))
	}

	return &HostMetricRepository{
		DB: db,
	}
}

func (r *HostMetricRepository) Exists(orgID, hostID, name string) bool {
	if err := r.DB.Where(&HostMetricValue{OrgID: orgID, HostID: hostID, Name: name}).First(&HostMetricValue{}).Error; err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		return false
	}

	return true
}

func (r *HostMetricRepository) Names(orgID, hostID string) (*domain.MetricNames, error) {
	result := make([]HostMetricValue, 0)
	if err := r.DB.Model(&HostMetricValue{}).Where(&HostMetricValue{OrgID: orgID, HostID: hostID}).Select("distinct(name)").Find(&result).Error; err != nil {
		return nil, fmt.Errorf("select distinct name from host_metric_values: %v", err)
	}

	names := make([]string, 0)
	for _, r := range result {
		names = append(names, r.Name)
	}

	return &domain.MetricNames{Names: names}, nil
}

func (r *HostMetricRepository) Values(orgID, hostID, name string, from, to int64) (*domain.MetricValues, error) {
	result := make([]HostMetricValue, 0)
	if err := r.DB.Where(&HostMetricValue{OrgID: orgID, HostID: hostID, Name: name}).Where("? < time and time < ?", from, to).Find(&result).Error; err != nil {
		return nil, fmt.Errorf("select time, value from host_metric_values: %v", err)
	}

	values := make([]domain.MetricValue, 0)
	for _, r := range result {
		values = append(values, domain.MetricValue{
			OrgID:  r.OrgID,
			HostID: r.HostID,
			Name:   r.Name,
			Time:   r.Time,
			Value:  r.Value,
		})
	}

	return &domain.MetricValues{Metrics: values}, nil
}

func (r *HostMetricRepository) ValuesLimit(orgID, hostID, name string, limit int) (*domain.MetricValues, error) {
	result := make([]HostMetricValue, 0)
	if err := r.DB.Where(&HostMetricValue{OrgID: orgID, HostID: hostID, Name: name}).Order("time desc").Limit(limit).Find(&result).Error; err != nil {
		return nil, fmt.Errorf("select time, value from host_metric_values: %v", err)
	}

	values := make([]domain.MetricValue, 0)
	for _, r := range result {
		values = append(values, domain.MetricValue{
			OrgID:  r.OrgID,
			HostID: r.HostID,
			Name:   r.Name,
			Time:   r.Time,
			Value:  r.Value,
		})
	}

	return &domain.MetricValues{Metrics: values}, nil
}

func (r *HostMetricRepository) ValuesLatest(orgID string, hostID, name []string) (*domain.TSDBLatest, error) {
	result := make([]HostMetricValuesLatest, 0)
	if len(hostID) > 0 && len(name) > 0 {
		if err := r.DB.Where("host_id IN (?)", hostID).Where("name IN (?)", name).Find(&result).Error; err != nil {
			return nil, fmt.Errorf("select * from host_metric_value_latest: %v", err)
		}
	} else {
		if err := r.DB.Where(&HostMetricValuesLatest{OrgID: orgID}).Find(&result).Error; err != nil {
			return nil, fmt.Errorf("select * from host_metric_value_latest: %v", err)
		}
	}

	latest := make(map[string]map[string]domain.MetricValue)
	for _, r := range result {
		if _, ok := latest[r.HostID]; !ok {
			latest[r.HostID] = make(map[string]domain.MetricValue)
		}

		latest[r.HostID][r.Name] = domain.MetricValue{Name: r.Name, Value: r.Value}
	}

	return &domain.TSDBLatest{TSDBLatest: latest}, nil
}

func (r *HostMetricRepository) Save(orgID string, values []domain.MetricValue) (*domain.Success, error) {
	if err := r.DB.Transaction(func(tx *gorm.DB) error {
		for i := range values {
			if err := tx.Create(&HostMetricValue{
				OrgID:  orgID,
				HostID: values[i].HostID,
				Name:   values[i].Name,
				Time:   values[i].Time,
				Value:  values[i].Value,
			}).Error; err != nil {
				return fmt.Errorf("insert into host_metric_values: %v", err)
			}

			where := HostMetricValuesLatest{OrgID: orgID, HostID: values[i].HostID, Name: values[i].Name}
			update := HostMetricValuesLatest{Value: values[i].Value}
			if err := tx.Where(&where).Assign(&update).FirstOrCreate(&HostMetricValuesLatest{}).Error; err != nil {
				return fmt.Errorf("insert into host_metric_values_latest: %v", err)
			}
		}

		return nil
	}); err != nil {
		return &domain.Success{Success: false}, fmt.Errorf("transaction: %v", err)
	}

	return &domain.Success{Success: true}, nil
}
