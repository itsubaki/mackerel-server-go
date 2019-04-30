package domain

type HostMetrics []HostMetric

type HostMetric struct {
	Name        string `json:"name"`
	DisplayName string `json:"displayName,omitempty"`
	IsStacked   bool   `json:"isStacked"`
}

type HostMetricValues []HostMetricValue

type HostMetricValue struct {
	HostID string  `json:"hostId,omitempty"`
	Name   string  `json:"name,omitempty"`
	Time   int64   `json:"time"`
	Value  float64 `json:"value"`
}
