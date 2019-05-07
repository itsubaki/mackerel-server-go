package memory

import "github.com/itsubaki/mackerel-api/pkg/domain"

type CheckReportRepository struct {
	CheckReports *domain.CheckReports
}

func NewCheckReportRepository() *CheckReportRepository {
	return &CheckReportRepository{
		CheckReports: &domain.CheckReports{
			Reports: []domain.CheckReport{},
		},
	}
}

func (repo *CheckReportRepository) Save(reports *domain.CheckReports) (*domain.Success, error) {
	repo.CheckReports.Reports = append(repo.CheckReports.Reports, reports.Reports...)
	return &domain.Success{Success: true}, nil
}
