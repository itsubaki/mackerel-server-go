package usecase

import (
	"github.com/itsubaki/mackerel-api/pkg/domain"
	"github.com/itsubaki/mackerel-api/pkg/interfaces/database"
)

func NewCheckInteractor() *CheckInteractor {
	return &CheckInteractor{
		CheckRepository: database.NewCheckReportRepository(),
	}
}

type CheckInteractor struct {
	CheckRepository *database.CheckReportRepository
}

func (s *CheckInteractor) Save(reports domain.CheckReports) (string, error) {
	return "status", nil
}
