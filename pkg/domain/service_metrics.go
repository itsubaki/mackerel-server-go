package domain

type ServiceMetricValues []ServiceMetricValue

type ServiceMetricValue struct {
	ServiceName string  `json:"-"`
	Name        string  `json:"name"`
	Time        int64   `json:"time"`
	Value       float64 `json:"value"`
}
