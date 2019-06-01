package usecase

import "github.com/itsubaki/mackerel-api/pkg/domain"

type CheckReportInteractor struct {
	CheckReportRepository CheckReportRepository
}

func (s *CheckReportInteractor) Save(orgID string, reports *domain.CheckReports) (*domain.Success, error) {
	for i := range reports.Reports {
		reports.Reports[i].Message = reports.Reports[i].Message[:len(reports.Reports[i].Message)-1] // remove \n
	}

	return s.CheckReportRepository.Save(orgID, reports)
}
