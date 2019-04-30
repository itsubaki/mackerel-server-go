package controllers

import "github.com/itsubaki/mackerel-api/pkg/domain"

type PostCheckReportInput struct {
	Reports domain.CheckReports `json:"reports"`
}

type PostCheckReportOutput struct {
	Status string `json:"status"`
}

func (m *Mackerel) PostCheckReport(in *PostCheckReportInput) (*PostCheckReportOutput, error) {
	return &PostCheckReportOutput{}, nil
}
