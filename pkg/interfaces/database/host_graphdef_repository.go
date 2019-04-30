package database

import "github.com/itsubaki/mackerel-api/pkg/domain"

type CustomGraphDefRepository struct {
	Internal domain.CustomHostMetricDefs
}

func NewCustomHostMetricDefRepository() *CustomGraphDefRepository {
	return &CustomGraphDefRepository{
		Internal: domain.CustomHostMetricDefs{},
	}
}

func (repo *CustomGraphDefRepository) Save(v domain.CustomHostMetricDef) error {
	repo.Internal = append(repo.Internal, v)
	return nil
}
