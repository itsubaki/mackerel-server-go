package usecase

import "github.com/itsubaki/mackerel-api/pkg/domain"

type CheckReportRepository interface {
	Save(orgID string, reports *domain.CheckReports) (*domain.Success, error)
}
