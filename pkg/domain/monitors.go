package domain

type Monitors struct {
	Monitors []interface{} `json:"monitors"`
}

type Monitoring struct {
	ID                              string              `json:"id,omitempty"`
	Type                            string              `json:"type"` // host, connectivity, service, external, expression
	Name                            string              `json:"name"`
	Memo                            string              `json:"memo,omitempty"`
	NotificationInterval            int                 `json:"notificationInterval,omitempty"`
	IsMute                          bool                `json:"isMute,omitempty"`
	Duration                        int                 `json:"duration,omitempty"`                        // HostMetric, Service
	Metric                          string              `json:"metric,omitempty"`                          // HostMetric, Service
	Operator                        string              `json:"operator,omitempty"`                        // HostMetric, Service
	Warning                         float64             `json:"warning,omitempty"`                         // HostMetric, Service
	Critical                        float64             `json:"critical,omitempty"`                        // HostMetric, Service
	MaxCheckAttempts                int                 `json:"maxCheckAttempts,omitempty"`                // HostMetric, Service
	Scopes                          []string            `json:"scopes,omitempty"`                          // HostMetric, Connectivity
	ExcludeScopes                   []string            `json:"excludeScopes,omitempty"`                   // HostMetric, Connectivity
	MissingDurationWarning          int                 `json:"missingDurationWarning,omitempty"`          // Service
	MissingDurationCritical         int                 `json:"missingDurationCritical,omitempty"`         // Service
	URL                             string              `json:"url,omitempty"`                             // External
	Method                          string              `json:"method,omitempty"`                          // External
	Service                         string              `json:"service,omitempty"`                         // External
	ResponseTimeWarning             int                 `json:"responseTimeWarning,omitempty"`             // External
	ResponseTimeCritical            int                 `json:"responseTimeCritical,omitempty"`            // External
	ResponseTimeDuration            int                 `json:"responseTimeDuration,omitempty"`            // External
	ContainsString                  string              `json:"containsString,omitempty"`                  // External
	CertificationExpirationWarning  int                 `json:"certificationExpirationWarning,omitempty"`  // External
	CertificationExpirationCritical int                 `json:"certificationExpirationCritical,omitempty"` // External
	SkipCertificateVerification     bool                `json:"skipCertificateVerification,omitempty"`     // External
	Headers                         []map[string]string `json:"notificationInterval,omitempty"`            // External
	RequestBody                     string              `json:"notificationInterval,omitempty"`            // External
	Expression                      string              `json:"expression,omitempty"`                      // Expression
}

func (m *Monitoring) Cast() interface{} {
	if m.Type == "host" {
		return &HostMetricMonitoring{
			ID:                   m.ID,
			Type:                 m.Type,
			Name:                 m.Name,
			Memo:                 m.Memo,
			Duration:             m.Duration,
			Metric:               m.Metric,
			Operator:             m.Operator,
			Warning:              m.Warning,
			Critical:             m.Critical,
			MaxCheckAttempts:     m.MaxCheckAttempts,
			NotificationInterval: m.NotificationInterval,
			Scopes:               m.Scopes,
			ExcludeScopes:        m.ExcludeScopes,
			IsMute:               m.IsMute,
		}
	}

	if m.Type == "connectivity" {
		return &HostConnectivityMonitoring{
			ID:                   m.ID,
			Type:                 m.Type,
			Name:                 m.Name,
			Memo:                 m.Memo,
			NotificationInterval: m.NotificationInterval,
			Scopes:               m.Scopes,
			ExcludeScopes:        m.ExcludeScopes,
			IsMute:               m.IsMute,
		}
	}

	if m.Type == "service" {
		return &ServiceMetricMonitoring{
			ID:                      m.ID,
			Type:                    m.Type,
			Name:                    m.Name,
			Memo:                    m.Memo,
			Duration:                m.Duration,
			Metric:                  m.Metric,
			Operator:                m.Operator,
			Warning:                 m.Warning,
			Critical:                m.Critical,
			MaxCheckAttempts:        m.MaxCheckAttempts,
			MissingDurationWarning:  m.MissingDurationWarning,
			MissingDurationCritical: m.MissingDurationCritical,
			NotificationInterval:    m.NotificationInterval,
			Scopes:                  m.Scopes,
			ExcludeScopes:           m.ExcludeScopes,
			IsMute:                  m.IsMute,
		}
	}

	if m.Type == "external" {
		return &ExternalMonitoring{
			ID:                              m.ID,
			Type:                            m.Type,
			Name:                            m.Name,
			Memo:                            m.Memo,
			URL:                             m.URL,
			Method:                          m.Method,
			Service:                         m.Service,
			NotificationInterval:            m.NotificationInterval,
			ResponseTimeWarning:             m.ResponseTimeWarning,
			ResponseTimeCritical:            m.ResponseTimeCritical,
			ResponseTimeDuration:            m.ResponseTimeDuration,
			ContainsString:                  m.ContainsString,
			MaxCheckAttempts:                m.MaxCheckAttempts,
			CertificationExpirationWarning:  m.CertificationExpirationWarning,
			CertificationExpirationCritical: m.CertificationExpirationCritical,
			IsMute:      m.IsMute,
			Headers:     m.Headers,
			RequestBody: m.RequestBody,
		}
	}

	if m.Type == "expression" {
		return &ExpressionMonitoring{
			ID:                   m.ID,
			Type:                 m.Type,
			Name:                 m.Name,
			Memo:                 m.Memo,
			Expression:           m.Expression,
			Operator:             m.Operator,
			Warning:              m.Warning,
			Critical:             m.Critical,
			NotificationInterval: m.NotificationInterval,
			IsMute:               m.IsMute,
		}
	}

	return m
}

type HostMetricMonitoring struct {
	ID                   string   `json:"id,omitempty"`
	Type                 string   `json:"type"` // host
	Name                 string   `json:"name"`
	Memo                 string   `json:"memo,omitempty"`
	Duration             int      `json:"duration"` // 1~10
	Metric               string   `json:"metric"`
	Operator             string   `json:"operator"` // > or <
	Warning              float64  `json:"warning"`
	Critical             float64  `json:"critical"`
	MaxCheckAttempts     int      `json:"maxCheckAttempts,omitempty"` //1~10
	NotificationInterval int      `json:"notificationInterval,omitempty"`
	Scopes               []string `json:"scopes,omitempty"`
	ExcludeScopes        []string `json:"excludeScopes,omitempty"`
	IsMute               bool     `json:"isMute,omitempty"`
}

type HostConnectivityMonitoring struct {
	ID                   string   `json:"id"`
	Type                 string   `json:"type"` // connectivity
	Name                 string   `json:"name"`
	Memo                 string   `json:"memo,omitempty"`
	NotificationInterval int      `json:"notificationInterval,omitempty"`
	Scopes               []string `json:"scopes,omitempty"`
	ExcludeScopes        []string `json:"excludeScopes,omitempty"`
	IsMute               bool     `json:"isMute,omitempty"`
}

type ServiceMetricMonitoring struct {
	ID                      string   `json:"id,omitempty"`
	Type                    string   `json:"type"` // service
	Name                    string   `json:"name"`
	Memo                    string   `json:"memo,omitempty"`
	Duration                int      `json:"duration"` // 1~10
	Metric                  string   `json:"metric"`
	Operator                string   `json:"operator"` // > or <
	Warning                 float64  `json:"warning"`
	Critical                float64  `json:"critical"`
	MaxCheckAttempts        int      `json:"maxCheckAttempts,omitempty"` //1~10
	MissingDurationWarning  int      `json:"missingDurationWarning,omitempty"`
	MissingDurationCritical int      `json:"missingDurationCritical,omitempty"`
	NotificationInterval    int      `json:"notificationInterval,omitempty"`
	Scopes                  []string `json:"scopes,omitempty"`
	ExcludeScopes           []string `json:"excludeScopes,omitempty"`
	IsMute                  bool     `json:"isMute,omitempty"`
}

type ExternalMonitoring struct {
	ID                              string              `json:"id,omitempty"`
	Type                            string              `json:"type"` // external
	Name                            string              `json:"name"`
	Memo                            string              `json:"memo,omitempty"`
	URL                             string              `json:"url"`
	Method                          string              `json:"method,omitempty"` // GET, PUT, POST, DELETE
	Service                         string              `json:"service,omitempty"`
	NotificationInterval            int                 `json:"notificationInterval,omitempty"`
	ResponseTimeWarning             int                 `json:"responseTimeWarning,omitempty"`
	ResponseTimeCritical            int                 `json:"responseTimeCritical,omitempty"`
	ResponseTimeDuration            int                 `json:"responseTimeDuration,omitempty"`
	ContainsString                  string              `json:"containsString,omitempty"`
	MaxCheckAttempts                int                 `json:"maxCheckAttempts,omitempty"`
	CertificationExpirationWarning  int                 `json:"certificationExpirationWarning,omitempty"`
	CertificationExpirationCritical int                 `json:"certificationExpirationCritical,omitempty"`
	SkipCertificateVerification     bool                `json:"skipCertificateVerification,omitempty"`
	IsMute                          bool                `json:"isMute,omitempty"`
	Headers                         []map[string]string `json:"notificationInterval,omitempty"`
	RequestBody                     string              `json:"notificationInterval,omitempty"`
}

type ExpressionMonitoring struct {
	ID                   string  `json:"id,omitempty"`
	Type                 string  `json:"type"` // expression
	Name                 string  `json:"name"`
	Memo                 string  `json:"memo,omitempty"`
	Expression           string  `json:"expression"`
	Operator             string  `json:"operator"` // > or <
	Warning              float64 `json:"warning"`
	Critical             float64 `json:"critical"`
	NotificationInterval int     `json:"notificationInterval,omitempty"`
	IsMute               bool    `json:"isMute,omitempty"`
}
