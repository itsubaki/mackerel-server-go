package usecase

import (
	"github.com/itsubaki/mackerel-api/pkg/domain"
	"github.com/itsubaki/mackerel-api/pkg/interfaces/database"
)

func NewCustomGraphDefInteractor() *CustomGraphDefInteractor {
	return &CustomGraphDefInteractor{
		CustomHostGraphDefRepository: database.NewCustomHostMetricDefRepository(),
	}
}

type CustomGraphDefInteractor struct {
	CustomHostGraphDefRepository *database.CustomHostGraphDefRepository
}

func (s *CustomGraphDefInteractor) Save(defs domain.CustomHostMetricDefs) error {
	return s.CustomHostGraphDefRepository.Save(defs)
}
