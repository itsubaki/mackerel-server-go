package usecase

import "github.com/itsubaki/mackerel-server-go/pkg/domain"

type CheckReportRepository interface {
	Save(orgID string, reports *domain.CheckReports) (*domain.Success, error)
}
