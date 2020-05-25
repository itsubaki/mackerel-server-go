package domain

type Downtimes struct {
	Downtimes []Downtime `json:"downtimes"`
}

type Downtime struct {
	OrgID                string      `json:"-"`
	ID                   string      `json:"id"`
	Name                 string      `json:"name"`
	Memo                 string      `json:"memo"`
	Start                int64       `json:"start"`
	Duration             int64       `json:"duration"`
	Recurrence           *Recurrence `json:"recurrence,omitempty"`
	ServiceScopes        []string    `json:"serviceScopes,omitempty"`
	ServiceExcludeScopes []string    `json:"serviceExcludeScopes,omitempty"`
	RoleScopes           []string    `json:"roleScopes,omitempty"`
	RoleExcludeScopes    []string    `json:"roleExcludeScopes,omitempty"`
	MonitorScopes        []string    `json:"monitorScopes,omitempty"`
	MonitorExcludeScopes []string    `json:"monitorExcludeScopes,omitempty"`
}

type Recurrence struct {
	Type     string   `json:"type"`
	Interval int64    `json:"interval"`
	Weekdays []string `json:"weekdays,omitempty"`
	Until    int64    `json:"until,omitempty"`
}
