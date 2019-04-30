package controllers

import "github.com/itsubaki/mackerel-api/pkg/domain"

type GetAlertInput struct {
	WithClosed bool   `json:"withClosed,omitempty"`
	NextID     string `json:"nextId,omitempty"`
	Limit      int64  `json:"limit,omitempty"`
}

type GetAlertOutput struct {
	Alerts domain.Alerts `json:"alerts"`
	NextID string        `json:"nextId"`
}

type PostAlertInput struct {
	AlertID string `json:"-"`
	Reason  string `json:"reason"`
}

type PostAlertOutput domain.Alert

func (m *Mackerel) GetAlert(in *GetAlertInput) (*GetAlertOutput, error) {
	return &GetAlertOutput{}, nil
}

func (m *Mackerel) PostAlert(in *PostAlertInput) (*PostAlertOutput, error) {
	return &PostAlertOutput{}, nil
}
