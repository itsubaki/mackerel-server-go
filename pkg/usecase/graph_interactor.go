package usecase

import "github.com/itsubaki/mackerel-api/pkg/domain"

type GraphInteractor struct {
	GraphRepository GraphRepository
}

func (s *GraphInteractor) SaveDef(orgID string, g []domain.GraphDef) (*domain.Success, error) {
	return s.GraphRepository.SaveDef(orgID, g)
}
