package domain

type Alerts []Alert

type Alert struct {
	ID        string  `json:"id"`
	Status    string  `json:"status"`
	MonitorID string  `json:"monitorId"`
	Type      string  `json:"type"`
	HostID    string  `json:"hostId,omitempty"`
	Value     float64 `json:"value,omitempty"`
	Message   string  `json:"message,omitempty"`
	Reason    string  `json:"reason,omitempty"`
	OpenedAt  string  `json:"openedAt"`
	ClosedAt  string  `json:"closedAt,omitempty"`
}
