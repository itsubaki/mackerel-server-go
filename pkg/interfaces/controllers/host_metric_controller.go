package controllers

import "github.com/itsubaki/mackerel-api/pkg/domain"

type PostHostMetricInput struct {
	MetricValue domain.HostMetricValues `json:"-"`
}

type PostHostMetricOutput struct {
	Success bool `json:"success"`
}

type GetHostMetricInput struct {
	HostID string `json:"-"`
	Name   string `json:"-"`
	From   string `json:"-"`
	To     string `json:"-"`
}

type GetHostMetricOutput struct {
	Metrics domain.HostMetricValues `json:"metrics"`
}

type GetHostMetricNamesInput struct {
	HostID string `json:"-"`
}

type GetHostMetricNamesOutput struct {
	Name []string `json:"names"`
}

type GetHostMetricLatestInput struct {
	HostID string `json:"-"`
	Name   string `json:"-"`
}

type GetHostMetricLatestOutput struct {
	TSDBLatest map[string]map[string]float64 `json:"tsdbLatest"`
}

func (m *Mackerel) PostHostMetric(in *PostHostMetricInput) (*PostHostMetricOutput, error) {
	return &PostHostMetricOutput{}, nil
}

func (m *Mackerel) GetHostMetric(in *GetHostMetricInput) (*GetHostMetricOutput, error) {
	return &GetHostMetricOutput{}, nil
}

func (m *Mackerel) GetHostMetricLatest(in *GetHostMetricLatestInput) (*GetHostMetricLatestOutput, error) {
	return &GetHostMetricLatestOutput{}, nil
}

func (m *Mackerel) GetHostMetricNames(in *GetHostMetricNamesInput) (*GetHostMetricNamesOutput, error) {
	return &GetHostMetricNamesOutput{}, nil
}
