package mackerel

type GetAlertInput struct {
	WithClosed bool   `json:"withClosed,omitempty"`
	NextID     string `json:"nextId,omitempty"`
	Limit      int64  `json:"limit,omitempty"`
}

type GetAlertOutput struct {
	Alerts []Alert `json:"alerts"`
	NextID string  `json:"nextId"`
}

type PostAlertInput struct {
	Reason string `json:"reason"`
}

type PostAlertOutput Alert

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

type AlertRepository struct {
	Internal []Alert
}

func NewAlertRepository() *AlertRepository {
	return &AlertRepository{
		Internal: []Alert{},
	}
}
