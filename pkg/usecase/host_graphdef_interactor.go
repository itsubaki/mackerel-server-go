package usecase

import "github.com/itsubaki/mackerel-api/pkg/domain"

type CustomGraphDefInteractor struct {
	CustomHostGraphDefRepository CustomHostGraphDefRepository
}

func (s *CustomGraphDefInteractor) Save(defs domain.CustomHostMetricDefs) error {
	return s.CustomHostGraphDefRepository.Save(defs)
}
