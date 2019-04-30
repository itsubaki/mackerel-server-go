package usecase

import "github.com/itsubaki/mackerel-api/pkg/domain"

type CustomHostGraphDefRepository interface {
	Save(v domain.CustomHostMetricDefs) error
}
