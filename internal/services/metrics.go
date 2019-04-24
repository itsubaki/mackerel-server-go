package services

type MetricNames struct {
	Name []string `json:"names"`
}

type GetMetricNamesInput struct {
	ServiceName string
}

type GetMetricNamesOutput struct {
	Status int `json:"-"`
	MetricNames
}
