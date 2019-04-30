package controllers

import "github.com/itsubaki/mackerel-api/pkg/domain"

type PostServiceMetricInput struct {
	ServiceName        string                     `json:"-"`
	ServiceMetricValue domain.ServiceMetricValues `json:"-"`
}

type PostServiceMetricOutput struct {
	Success bool `json:"success"`
}

type GetServiceMetricInput struct {
	ServiceName string `json:"-"`
	Name        string `json:"-"`
	From        string `json:"-"`
	To          string `json:"-"`
}

type GetServiceMetricOutput struct {
	Metrics domain.ServiceMetricValues `json:"metrics"`
}

type GetServiceMetricNamesInput struct {
	ServiceName string `json:"-"`
}

type GetServiceMetricNamesOutput struct {
	Name []string `json:"names"`
}

type ServiceMetricValue struct {
	ServiceName string  `json:"-"`
	Name        string  `json:"name"`
	Time        int64   `json:"time"`
	Value       float64 `json:"value"`
}

func (m *Mackerel) GetServiceMetricNames(in *GetServiceMetricNamesInput) (*GetServiceMetricNamesOutput, error) {
	return &GetServiceMetricNamesOutput{}, nil
}

func (m *Mackerel) PostServiceMetric(in *PostServiceMetricInput) (*PostServiceMetricOutput, error) {
	return &PostServiceMetricOutput{}, nil
}

func (m *Mackerel) GetServiceMetric(in *GetServiceMetricInput) (*GetServiceMetricOutput, error) {
	return &GetServiceMetricOutput{}, nil
}

func (m *Mackerel) GetServiceMetadata(in *GetServiceMetadataInput) (GetServiceMetadataOutput, error) {
	return "", nil
}
