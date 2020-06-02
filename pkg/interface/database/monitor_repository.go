package database

import (
	"encoding/json"
	"fmt"

	"github.com/jinzhu/gorm"

	"github.com/itsubaki/mackerel-server-go/pkg/domain"
)

type MonitorRepository struct {
	SQLHandler
	DB *gorm.DB
}

type Monitor struct {
	OrgID                           string  `gorm:"column:org_id;                            type:varchar(16);  not null;"`
	ID                              string  `gorm:"column:id;                                type:varchar(16);  not null; primary_key"`
	Type                            string  `gorm:"column:type;                              type:enum('host', 'connectivity', 'service', 'external', 'expression')"`
	Name                            string  `gorm:"column:name;                              type:varchar(128); not null;"`
	Memo                            string  `gorm:"column:memo;                              type:varchar(128); not null; default:''"`
	NotificationInterval            int     `gorm:"column:notification_interval;             type:int;          not null; default: '1'"`
	IsMute                          bool    `gorm:"column:is_mute;                           type:bool;         not null; default: '0'"`
	Duration                        int     `gorm:"column:duration;                          type:int;"`
	Metric                          string  `gorm:"column:metric;                            type:varchar(128);"`
	Operator                        string  `gorm:"column:operator;                          type:enum('>', '<'); not null; default: '<'"`
	Warning                         float64 `gorm:"column:warning;                           type:double;"`
	Critical                        float64 `gorm:"column:critical;                          type:double;"`
	MaxCheckAttempts                int     `gorm:"column:max_check_attempts;                type:int;"`
	Scopes                          string  `gorm:"column:scopes;                            type:text;"`
	ExcludeScopes                   string  `gorm:"column:exclude_scopes;                    type:text;"`
	MissingDurationWarning          int     `gorm:"column:missing_duration_warning;          type:int;"`
	MissingDurationCritical         int     `gorm:"column:missing_duration_critical;         type:int;"`
	URL                             string  `gorm:"column:url;                               type:text;"`
	Method                          string  `gorm:"column:method;                            type:enum('GET', 'PUT', 'POST', 'DELETE', '');"`
	Service                         string  `gorm:"column:service;                           type:text;"`
	ResponseTimeWarning             int     `gorm:"column:response_time_warning;             type:int;"`
	ResponseTimeCritical            int     `gorm:"column:response_time_critical;            type:int;"`
	ResponseTimeDuration            int     `gorm:"column:response_time_duration;            type:int;"`
	ContainsString                  string  `gorm:"column:contains_string;                   type:text;"`
	CertificationExpirationWarning  int     `gorm:"column:certification_expiration_warning;  type:int;"`
	CertificationExpirationCritical int     `gorm:"column:certification_expiration_critical; type:int;"`
	SkipCertificateVerification     bool    `gorm:"column:skip_certificate_verification;     type:bool;"`
	Headers                         string  `gorm:"column:headers;                           type:text;"`
	RequestBody                     string  `gorm:"column:request_body;                      type:text;"`
	Expression                      string  `gorm:"column:expression;                        type:text;"`
}

func (m Monitor) Domain() (domain.Monitoring, error) {
	monitoring := domain.Monitoring{
		OrgID:                           m.OrgID,
		ID:                              m.ID,
		Type:                            m.Type,
		Name:                            m.Name,
		Memo:                            m.Memo,
		NotificationInterval:            m.NotificationInterval,
		IsMute:                          m.IsMute,
		Duration:                        m.Duration,
		Metric:                          m.Metric,
		Operator:                        m.Operator,
		Warning:                         m.Warning,
		Critical:                        m.Critical,
		MaxCheckAttempts:                m.MaxCheckAttempts,
		MissingDurationWarning:          m.MissingDurationWarning,
		MissingDurationCritical:         m.MissingDurationCritical,
		URL:                             m.URL,
		Method:                          m.Method,
		ResponseTimeWarning:             m.ResponseTimeWarning,
		ResponseTimeCritical:            m.ResponseTimeCritical,
		ResponseTimeDuration:            m.ResponseTimeDuration,
		ContainsString:                  m.ContainsString,
		CertificationExpirationWarning:  m.CertificationExpirationWarning,
		CertificationExpirationCritical: m.CertificationExpirationCritical,
		SkipCertificateVerification:     m.SkipCertificateVerification,
		Expression:                      m.Expression,
	}

	if err := json.Unmarshal([]byte(m.Scopes), &monitoring.Scopes); err != nil {
		return monitoring, fmt.Errorf("unmarshal monitoring.Scopes: %v", err)
	}

	if err := json.Unmarshal([]byte(m.ExcludeScopes), &monitoring.ExcludeScopes); err != nil {
		return monitoring, fmt.Errorf("unmarshal monitoring.ExcludeScopes: %v", err)
	}

	if err := json.Unmarshal([]byte(m.Service), &monitoring.Service); err != nil {
		return monitoring, fmt.Errorf("unmarshal monitoring.Service: %v", err)
	}

	if err := json.Unmarshal([]byte(m.Headers), &monitoring.Headers); err != nil {
		return monitoring, fmt.Errorf("unmarshal monitoring.Headers: %v", err)
	}

	if err := json.Unmarshal([]byte(m.RequestBody), &monitoring.RequestBody); err != nil {
		return monitoring, fmt.Errorf("unmarshal monitoring.RequestBody: %v", err)
	}

	return monitoring, nil
}

func NewMonitorRepository(handler SQLHandler) *MonitorRepository {
	db, err := gorm.Open(handler.Dialect(), handler.Raw())
	if err != nil {
		panic(err)
	}
	db.LogMode(handler.IsDebugging())

	if err := db.AutoMigrate(&Monitor{}).Error; err != nil {
		panic(fmt.Errorf("auto migrate monitoring: %v", err))
	}

	return &MonitorRepository{
		SQLHandler: handler,
		DB:         db,
	}
}

func (repo *MonitorRepository) ListHostMetric(orgID string) ([]domain.HostMetricMonitoring, error) {
	list := make([]domain.HostMetricMonitoring, 0)
	monitors, err := repo.List(orgID)
	if err != nil {
		return nil, fmt.Errorf("list monitors: %v", err)
	}

	for i := range monitors.Monitors {
		m, ok := monitors.Monitors[i].(*domain.HostMetricMonitoring)
		if !ok {
			continue
		}

		list = append(list, *m)
	}

	return list, nil
}

func (repo *MonitorRepository) List(orgID string) (*domain.Monitors, error) {
	result := make([]Monitor, 0)
	if err := repo.DB.Where(&Monitor{OrgID: orgID}).Find(&result).Error; err != nil {
		return nil, fmt.Errorf("select * from monitors: %v", err)
	}

	out := make([]interface{}, 0)
	for _, r := range result {
		m, err := r.Domain()
		if err != nil {
			return nil, fmt.Errorf("domain: %v", err)
		}

		out = append(out, m.Cast())
	}

	return &domain.Monitors{Monitors: out}, nil
}

func (repo *MonitorRepository) Save(orgID string, monitor *domain.Monitoring) (interface{}, error) {
	if err := repo.Transact(func(tx Tx) error {
		scopes, err := json.Marshal(monitor.Scopes)
		if err != nil {
			return fmt.Errorf("marshal monitor.Scopes: %v", err)
		}

		exclude, err := json.Marshal(monitor.ExcludeScopes)
		if err != nil {
			return fmt.Errorf("marshal monitor.ExcludeScopes: %v", err)
		}

		service, err := json.Marshal(monitor.Service)
		if err != nil {
			return fmt.Errorf("marshal monitor.Service: %v", err)
		}

		headers, err := json.Marshal(monitor.Headers)
		if err != nil {
			return fmt.Errorf("marshal monitor.Headers: %v", err)
		}

		body, err := json.Marshal(monitor.RequestBody)
		if err != nil {
			return fmt.Errorf("marshal monitor.RequestBody: %v", err)
		}

		if _, err := tx.Exec(
			`
			insert into monitors (
				org_id,
				id,
				type,
				name,
				memo,
				notification_interval,
				is_mute,
				duration,
				metric,
				operator,
				warning,
				critical,
				max_check_attempts,
				scopes,
				exclude_scopes,
				missing_duration_warning,
				missing_duration_critical,
				url,
				method,
				service,
				response_time_warning,
				response_time_critical,
				response_time_duration,
				contains_string,
				certification_expiration_warning,
				certification_expiration_critical,
				skip_certificate_verification,
				headers,
				request_body,
				expression
			)
			values (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
			`,
			orgID,
			monitor.ID,
			monitor.Type,
			monitor.Name,
			monitor.Memo,
			monitor.NotificationInterval,
			monitor.IsMute,
			monitor.Duration,
			monitor.Metric,
			monitor.Operator,
			monitor.Warning,
			monitor.Critical,
			monitor.MaxCheckAttempts,
			scopes,
			exclude,
			monitor.MissingDurationWarning,
			monitor.MissingDurationCritical,
			monitor.URL,
			monitor.Method,
			service,
			monitor.ResponseTimeWarning,
			monitor.ResponseTimeCritical,
			monitor.ResponseTimeDuration,
			monitor.ContainsString,
			monitor.CertificationExpirationWarning,
			monitor.CertificationExpirationCritical,
			monitor.SkipCertificateVerification,
			headers,
			body,
			monitor.Expression,
		); err != nil {
			return fmt.Errorf("insert into monitors: %v", err)

		}

		return nil
	}); err != nil {
		return nil, fmt.Errorf("transaction: %v", err)
	}

	return monitor, nil
}

func (repo *MonitorRepository) Update(orgID string, monitor *domain.Monitoring) (interface{}, error) {
	if err := repo.Transact(func(tx Tx) error {
		scopes, err := json.Marshal(monitor.Scopes)
		if err != nil {
			return fmt.Errorf("marshal monitor.Scopes: %v", err)
		}

		exclude, err := json.Marshal(monitor.ExcludeScopes)
		if err != nil {
			return fmt.Errorf("marshal monitor.ExcludeScopes: %v", err)
		}

		service, err := json.Marshal(monitor.Service)
		if err != nil {
			return fmt.Errorf("marshal monitor.Service: %v", err)
		}

		headers, err := json.Marshal(monitor.Headers)
		if err != nil {
			return fmt.Errorf("marshal monitor.Headers: %v", err)
		}

		body, err := json.Marshal(monitor.RequestBody)
		if err != nil {
			return fmt.Errorf("marshal monitor.RequestBody: %v", err)
		}

		if _, err := tx.Exec(
			`
			update monitors set
				type=?,
				name=?,
				memo=?,
				notification_interval=?,
				is_mute=?,
				duration=?,
				metric=?,
				operator=?,
				warning=?,
				critical=?,
				max_check_attempts=?,
				scopes=?,
				exclude_scopes=?,
				missing_duration_warning=?,
				missing_duration_critical=?,
				url=?,
				method=?,
				service=?,
				response_time_warning=?,
				response_time_critical=?,
				response_time_duration=?,
				contains_string=?,
				certification_expiration_warning=?,
				certification_expiration_critical=?,
				skip_certificate_verification=?,
				headers=?,
				request_body=?,
				expression=?
			where org_id=? and id=?
			`,
			monitor.Type,
			monitor.Name,
			monitor.Memo,
			monitor.NotificationInterval,
			monitor.IsMute,
			monitor.Duration,
			monitor.Metric,
			monitor.Operator,
			monitor.Warning,
			monitor.Critical,
			monitor.MaxCheckAttempts,
			scopes,
			exclude,
			monitor.MissingDurationWarning,
			monitor.MissingDurationCritical,
			monitor.URL,
			monitor.Method,
			service,
			monitor.ResponseTimeWarning,
			monitor.ResponseTimeCritical,
			monitor.ResponseTimeDuration,
			monitor.ContainsString,
			monitor.CertificationExpirationWarning,
			monitor.CertificationExpirationCritical,
			monitor.SkipCertificateVerification,
			headers,
			body,
			monitor.Expression,
			orgID,
			monitor.ID,
		); err != nil {
			return fmt.Errorf("update monitors: %v", err)

		}

		return nil
	}); err != nil {
		return nil, fmt.Errorf("transaction: %v", err)
	}

	return monitor.Cast(), nil
}

func (repo *MonitorRepository) Monitor(orgID, monitorID string) (interface{}, error) {
	result := Monitor{}
	if err := repo.DB.Where(&Monitor{OrgID: orgID, ID: monitorID}).First(&result).Error; err != nil {
		return nil, fmt.Errorf("select * from monitors: %v", err)
	}

	m, err := result.Domain()
	if err != nil {
		return nil, fmt.Errorf("domain: %v", err)
	}

	return m.Cast(), nil
}

func (repo *MonitorRepository) Delete(orgID, monitorID string) (interface{}, error) {
	result := Monitor{}
	if err := repo.DB.Where(&Monitor{OrgID: orgID, ID: monitorID}).First(&result).Error; err != nil {
		return nil, fmt.Errorf("select * from monitors: %v", err)
	}

	if err := repo.DB.Delete(&Monitor{OrgID: orgID, ID: monitorID}).Error; err != nil {
		return nil, fmt.Errorf("delete from monitors: %v", err)
	}

	m, err := result.Domain()
	if err != nil {
		return nil, fmt.Errorf("domain: %v", err)
	}

	return m.Cast(), nil
}
