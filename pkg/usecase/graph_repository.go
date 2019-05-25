package usecase

import "github.com/itsubaki/mackerel-api/pkg/domain"

type GraphRepository interface {
	SaveDef(org string, g []domain.GraphDef) (*domain.Success, error)
}
