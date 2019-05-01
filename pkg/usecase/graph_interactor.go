package usecase

import "github.com/itsubaki/mackerel-api/pkg/domain"

type GraphInteractor struct {
	GraphRepository GraphRepository
}

func (s *GraphInteractor) Save(g []domain.GraphDef) (*domain.Success, error) {
	return s.GraphRepository.Save(g)
}
