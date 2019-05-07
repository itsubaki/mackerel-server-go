package usecase

import "github.com/itsubaki/mackerel-api/pkg/domain"

type CheckReportInteractor struct {
	CheckReportRepository CheckReportRepository
}

func (s *CheckReportInteractor) Save(reports *domain.CheckReports) (*domain.Success, error) {
	return s.CheckReportRepository.Save(reports)
}