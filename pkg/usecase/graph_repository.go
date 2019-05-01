package usecase

import "github.com/itsubaki/mackerel-api/pkg/domain"

type GraphRepository interface {
	Save(g []domain.GraphDef) (*domain.Success, error)
}
