package usecase

import "github.com/itsubaki/mackerel-api/pkg/domain"

type CheckInteractor struct {
	CheckRepository CheckRepository
}

func (s *CheckInteractor) Save(reports domain.CheckReports) (string, error) {
	return "status", nil
}
