package database

import "github.com/itsubaki/mackerel-api/pkg/domain"

type CheckReportRepository struct {
	Internal domain.CheckReports
}

func (repo *CheckReportRepository) Save(v domain.CheckReport) error {
	repo.Internal = append(repo.Internal, v)
	return nil
}
