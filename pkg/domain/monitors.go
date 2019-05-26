package domain

type Monitoring struct {
	Type                 string `json:"type"` // host, connectivity, service, external, expression
	Name                 string `json:"name"`
	NotificationInterval int    `json:"notificationInterval,omitempty"`
	IsMute               bool   `json:"isMute,omitempty"`
}

type HostMetricMonitoring struct {
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
	Type                 string   `json:"type"` // connectivity
	Name                 string   `json:"name"`
	Memo                 string   `json:"memo,omitempty"`
	NotificationInterval int      `json:"notificationInterval,omitempty"`
	Scopes               []string `json:"scopes,omitempty"`
	ExcludeScopes        []string `json:"excludeScopes,omitempty"`
	IsMute               bool     `json:"isMute,omitempty"`
}

type ServiceMetricMonitoring struct {
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
	Type                 string  `json:"type"` // expression
	Name                 string  `json:"name"`
	Memo                 string  `json:"memo,omitempty"`
	Expression           int     `json:"expression"`
	Operator             string  `json:"operator"` // > or <
	Warning              float64 `json:"warning"`
	Critical             float64 `json:"critical"`
	NotificationInterval int     `json:"notificationInterval,omitempty"`
	IsMute               bool    `json:"isMute,omitempty"`
}
