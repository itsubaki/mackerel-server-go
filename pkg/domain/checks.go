package domain

type CheckReports []CheckReport

type CheckReport struct {
	Source               Source `json:"source"`
	Name                 string `json:"name"`
	Status               string `json:"status"` // OK, CRITICAL, WARNING, UNKNOWN
	Message              string `json:"message"`
	OccurredAt           string `json:"occurredAt"`
	NotificationInterval string `json:"notificationInterval,omitempty"`
	MaxCheckAttempts     string `json:"maxCheckAttempts,omitempty"`
}

type Source struct {
	Type   string `json:"type"`
	HostID string `json:"hostId"`
}

type CheckReportStatus struct {
	Status string `json:"status"`
}
