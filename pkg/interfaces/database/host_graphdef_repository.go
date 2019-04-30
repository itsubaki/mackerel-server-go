package database

import "github.com/itsubaki/mackerel-api/pkg/domain"

type CustomHostGraphDefRepository struct {
	Internal domain.CustomHostMetricDefs
}

func NewCustomHostMetricDefRepository() *CustomHostGraphDefRepository {
	return &CustomHostGraphDefRepository{
		Internal: domain.CustomHostMetricDefs{},
	}
}

func (repo *CustomHostGraphDefRepository) Save(v domain.CustomHostMetricDefs) error {
	repo.Internal = append(repo.Internal, v...)
	return nil
}
