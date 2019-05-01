package database

import "github.com/itsubaki/mackerel-api/pkg/domain"

type CheckReportRepository struct {
	CheckReports domain.CheckReports
}

func (repo *CheckReportRepository) Save(v domain.CheckReport) error {
	repo.CheckReports = append(repo.CheckReports, v)
	return nil
}
