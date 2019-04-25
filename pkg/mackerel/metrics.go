package mackerel

type GetMetricNamesInput struct {
	ServiceName string
}

type GetMetricNamesOutput MetricNames

type GetHostMetricsInput struct{}

type GetHostMetricsOutput struct {
	Metrics []MetricValue `json:"metrics"`
}

type MetricNames struct {
	Name []string `json:"names"`
}

type MetricValue struct {
	HostID string  `json:"hostId,omitempty"`
	Name   string  `json:"name,omitempty"`
	Time   int64   `json:"time"`
	Value  float64 `json:"value"`
}
