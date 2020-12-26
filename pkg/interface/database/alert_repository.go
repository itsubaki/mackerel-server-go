package database

import (
	"errors"
	"fmt"
	"time"

	"github.com/itsubaki/mackerel-server-go/pkg/domain"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type AlertRepository struct {
	DB *gorm.DB
}

type Alert struct {
	OrgID     string  `gorm:"column:org_id;     type:varchar(16); not null"`
	ID        string  `gorm:"column:id;         type:varchar(16); not null;"`
	Status    string  `gorm:"column:status;     type:enum('OK', 'CRITICAL', 'WARNING', 'UNKNOWN'); not null"`
	MonitorID string  `gorm:"column:monitor_id; type:varchar(16); not null;"`
	Type      string  `gorm:"column:type;       type:enum('connectivity', 'host', 'service', 'external', 'check', 'expression'); not null;"`
	HostID    string  `gorm:"column:host_id;    type:varchar(16);"`
	Value     float64 `gorm:"column:value;      type:double;"`
	Message   string  `gorm:"column:message;    type:text;"`
	Reason    string  `gorm:"column:reason;     type:text;"`
	OpenedAt  int64   `gorm:"column:opened_at;  type:bigint;"`
	ClosedAt  int64   `gorm:"column:closed_at;  type:bigint;"`
}

func (a Alert) Domain() domain.Alert {
	return domain.Alert{
		OrgID:     a.OrgID,
		ID:        a.ID,
		Status:    a.Status,
		MonitorID: a.MonitorID,
		Type:      a.Type,
		HostID:    a.HostID,
		Value:     a.Value,
		Message:   a.Message,
		Reason:    a.Reason,
		OpenedAt:  a.OpenedAt,
		ClosedAt:  a.ClosedAt,
	}
}

type AlertHistory struct {
	OrgID     string `gorm:"column:org_id;     type:varchar(16); not null"`
	AlertID   string `gorm:"column:alert_id;   type:varchar(16); not null;"`
	Status    string `gorm:"column:status;     type:enum('OK', 'CRITICAL', 'WARNING', 'UNKNOWN'); not null"`
	MonitorID string `gorm:"column:monitor_id; type:varchar(16); not null;"`
	HostID    string `gorm:"column:host_id;    type:varchar(16);"`
	Time      int64  `gorm:"column:time;       type:bigint; not null"`
	Message   string `gorm:"column:message;    type:text;"`
}

func (a AlertHistory) TableName() string {
	return "alert_history"
}

func NewAlertRepository(handler SQLHandler) *AlertRepository {
	db, err := gorm.Open(mysql.New(mysql.Config{
		Conn: handler.Raw().(gorm.ConnPool),
	}), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	if handler.IsDebugging() {
		db.Logger.LogMode(4)
	}

	if err := handler.Transact(func(tx Tx) error {
		if _, err := tx.Exec(
			`
			create table if not exists alerts (
				org_id     varchar(16) not null,
				id         varchar(16) not null primary key,
				status     enum('OK', 'CRITICAL', 'WARNING', 'UNKNOWN') not null,
				monitor_id varchar(16) not null,
				type       enum('connectivity', 'host', 'service', 'external', 'check', 'expression') not null,
				host_id    varchar(16),
				value      double,
				message    text,
				reason     text,
				opened_at  bigint,
				closed_at  bigint,
				index(org_id, opened_at desc)
			)
			`,
		); err != nil {
			return fmt.Errorf("create table alerts: %v", err)
		}

		if _, err := tx.Exec(
			`
			create table if not exists alert_history (
				org_id     varchar(16) not null,
				alert_id   varchar(16) not null,
				status     enum('OK', 'CRITICAL', 'WARNING', 'UNKNOWN') not null,
				monitor_id varchar(16) not null,
				host_id    varchar(16),
				time       bigint      not null,
				message    text,
				primary key(alert_id, monitor_id, time desc)
			)
			`,
		); err != nil {
			return fmt.Errorf("create table alert_history: %v", err)
		}

		return nil
	}); err != nil {
		panic(fmt.Errorf("transaction: %v", err))
	}

	return &AlertRepository{
		DB: db,
	}
}

func (repo *AlertRepository) Exists(orgID, alertID string) bool {
	if err := repo.DB.Where(&Alert{OrgID: orgID, ID: alertID}).First(&Alert{}).Error; err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		return false
	}

	return true
}

func (repo *AlertRepository) Save(orgID string, alert *domain.Alert) (*domain.Alert, error) {
	if err := repo.DB.Transaction(func(tx *gorm.DB) error {
		result := AlertHistory{}
		if err := tx.Where(&AlertHistory{OrgID: orgID, HostID: alert.HostID, MonitorID: alert.MonitorID}).Order("time desc").First(&result).Error; err != nil && alert.Status == "OK" {
			// no record and no alert
			return nil
		}

		if result.Status == "OK" && alert.Status == "OK" {
			// have record and alert closed
			return nil
		}

		if (result.Status == "" || result.Status == "OK") && alert.Status != "OK" {
			// new alert
			result.AlertID = alert.ID
		}

		// status != "OK" && reports.Reports[i].Status != "OK"
		// -> continuous alert
		// status != "OK" && reports.Reports[i].Status == "OK"
		// -> close alert

		create := AlertHistory{
			OrgID:     orgID,
			AlertID:   result.AlertID,
			Status:    alert.Status,
			MonitorID: alert.MonitorID,
			HostID:    alert.HostID,
			Time:      alert.OpenedAt,
			Message:   alert.Message,
		}

		if err := tx.Create(&create).Error; err != nil {
			return fmt.Errorf("insert into alert_history: %v", err)
		}

		return nil
	}); err != nil {
		return alert, fmt.Errorf("transaction: %v", err)
	}

	if err := repo.DB.Transaction(func(tx *gorm.DB) error {
		result := AlertHistory{}
		if err := tx.Where(&AlertHistory{OrgID: orgID, HostID: alert.HostID, MonitorID: alert.MonitorID}).Order("time desc").First(&result).Error; err != nil {
			// no record
			return nil
		}

		if result.Status == "OK" && alert.Status == "OK" {
			// have record and alert closed
			return nil
		}

		var closedAt int64
		if alert.Status == "OK" {
			closedAt = result.Time
		}

		update := Alert{
			OrgID:     orgID,
			ID:        result.AlertID,
			Status:    result.Status,
			MonitorID: alert.MonitorID,
			Type:      alert.Type,
			HostID:    result.HostID,
			Value:     alert.Value,
			Message:   result.Message,
			Reason:    alert.Reason,
			OpenedAt:  result.Time,
			ClosedAt:  closedAt,
		}

		if err := tx.Where(&Alert{ID: result.AlertID}).Assign(&update).FirstOrCreate(&Alert{}).Error; err != nil {
			return fmt.Errorf("insert into alerts: %v", err)
		}

		return nil
	}); err != nil {
		return alert, fmt.Errorf("transaction: %v", err)
	}

	return alert, nil
}

func (repo *AlertRepository) List(orgID string, withClosed bool, nextID string, limit int) (*domain.Alerts, error) {
	status := "UNKNOWN"
	if withClosed {
		status = "OK"
	}

	result := make([]Alert, 0)
	if err := repo.DB.Where(&Alert{OrgID: orgID}).Where("status IN ('CRITICAL', 'WARNING', 'UNKNOWN', ?)", status).Order("opened_at desc").Limit(limit).Find(&result).Error; err != nil {
		return nil, fmt.Errorf("select * from alerts: %v", err)
	}

	alerts := make([]domain.Alert, 0)
	for _, r := range result {
		alerts = append(alerts, r.Domain())
	}

	if len(alerts) > limit {
		return &domain.Alerts{
			Alerts: alerts[:len(alerts)-1],
			NextID: alerts[len(alerts)-1].ID,
		}, nil
	}

	return &domain.Alerts{Alerts: alerts}, nil
}

func (repo *AlertRepository) Close(orgID, alertID, reason string) (*domain.Alert, error) {
	var alert domain.Alert
	if err := repo.DB.Transaction(func(tx *gorm.DB) error {
		update := Alert{
			Status:   "OK",
			Reason:   reason,
			ClosedAt: time.Now().Unix(),
		}

		if err := tx.Model(&Alert{}).Where(&Alert{OrgID: orgID, ID: alertID}).Updates(&update).Error; err != nil {
			return fmt.Errorf("update alerts: %v", err)
		}

		result := Alert{}
		if err := tx.Where(&Alert{OrgID: orgID, ID: alertID}).First(&result).Error; err != nil {
			return fmt.Errorf("select * from alerts: %v", err)
		}

		create := AlertHistory{
			OrgID:     result.OrgID,
			AlertID:   result.ID,
			Status:    "OK",
			HostID:    result.HostID,
			MonitorID: result.MonitorID,
			Time:      result.ClosedAt,
			Message:   result.Message,
		}

		if err := tx.Create(&create).Error; err != nil {
			return fmt.Errorf("insert into alert_history: %v", err)
		}

		alert = result.Domain()

		return nil
	}); err != nil {
		return nil, fmt.Errorf("transaction: %v", err)
	}

	return &alert, nil
}
