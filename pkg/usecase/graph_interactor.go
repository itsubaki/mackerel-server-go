package usecase

import "github.com/itsubaki/mackerel-api/pkg/domain"

type GraphInteractor struct {
	GraphRepository GraphRepository
}

func (s *GraphInteractor) SaveDef(orgID string, g []domain.GraphDef) (*domain.Success, error) {
	return s.GraphRepository.SaveDef(orgID, g)
}

func (s *GraphInteractor) List(orgID string) (*domain.GraphAnnotations, error) {
	return s.GraphRepository.List(orgID)
}

func (s *GraphInteractor) Save(orgID string, annotation *domain.GraphAnnotation) (*domain.GraphAnnotation, error) {
	annotation.ID = domain.NewRandomID(11)
	return s.GraphRepository.Save(orgID, annotation)
}

func (s *GraphInteractor) Update(orgID string, annotation *domain.GraphAnnotation) (*domain.GraphAnnotation, error) {
	return s.GraphRepository.Update(orgID, annotation)
}

func (s *GraphInteractor) Delete(orgID, annotationID string) (*domain.GraphAnnotation, error) {
	return s.GraphRepository.Delete(orgID, annotationID)
}
