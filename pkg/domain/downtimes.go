package domain

type Downtime struct {
	OrgID                string     `json:"-"`
	ID                   string     `json:"id"`
	Name                 string     `json:"name"`
	Start                int64      `json:"start"`
	Duration             int64      `json:"duration"`
	Recurrence           Recurrence `json:"recurrence"`
	ServiceScopes        []string   `json:"serviceScopes"`
	ServiceExcludeScopes []string   `json:"serviceExcludeScopes"`
	RoleScopes           []string   `json:"roleScopes"`
	RoleExcludeScopes    []string   `json:"roleExcludeScopes"`
	MonitorScopes        []string   `json:"monitorScopes"`
	MonitorExcludeScopes []string   `json:"monitorExcludeScopes"`
}

type Recurrence struct {
	Type     string   `json:"type"`
	Interval int64    `json:"interval"`
	Weekdays []string `json:"weekdays"`
	Until    int64    `json:"until"`
}
