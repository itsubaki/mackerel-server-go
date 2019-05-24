package memory

import "github.com/itsubaki/mackerel-api/pkg/domain"

type GraphRepository struct {
	GraphDef []domain.GraphDef
}

func NewGraphRepository() *GraphRepository {
	return &GraphRepository{
		GraphDef: []domain.GraphDef{},
	}
}

func (repo *GraphRepository) Save(org string, g []domain.GraphDef) (*domain.Success, error) {
	repo.GraphDef = append(repo.GraphDef, g...)
	return &domain.Success{Success: true}, nil
}
