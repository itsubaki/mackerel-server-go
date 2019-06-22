package domain

type CheckReports struct {
	Reports []CheckReport `json:"reports"`
}

type CheckReport struct {
	OrgID                string `json:"-"`
	Source               Source `json:"source"`
	Name                 string `json:"name"`
	Status               string `json:"status"` // OK, CRITICAL, WARNING, UNKNOWN
	Message              string `json:"message"`
	OccurredAt           int64  `json:"occurredAt"`
	NotificationInterval int64  `json:"notificationInterval,omitempty"`
	MaxCheckAttempts     int64  `json:"maxCheckAttempts,omitempty"`
}

type Source struct {
	Type   string `json:"type"`
	HostID string `json:"hostId"`
}
