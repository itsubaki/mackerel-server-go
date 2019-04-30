package controllers

import "github.com/itsubaki/mackerel-api/pkg/domain"

type PostCustomHostMetricDefInput struct {
	CustomHostMetricDef domain.CustomHostMetricDefs `json:"-"`
}

type PostCustomHostMetricDefOutput struct {
	Success bool `json:"success"`
}

func (m *Mackerel) PostCustomHostMetricDef(in *PostCustomHostMetricDefInput) (*PostCustomHostMetricDefOutput, error) {
	return &PostCustomHostMetricDefOutput{}, nil
}
