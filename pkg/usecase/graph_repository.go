package usecase

import "github.com/itsubaki/mackerel-api/pkg/domain"

type GraphRepository interface {
	Save(org string, g []domain.GraphDef) (*domain.Success, error)
}
