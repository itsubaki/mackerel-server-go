package usecase

import (
	"fmt"

	"github.com/itsubaki/mackerel-server-go/domain"
)

type GraphInteractor struct {
	GraphRepository GraphRepository
}

func (intr *GraphInteractor) SaveDef(orgID string, g []domain.GraphDef) (*domain.Success, error) {
	res, err := intr.GraphRepository.SaveDef(orgID, g)
	if err != nil {
		return res, fmt.Errorf("save graph definition: %v", err)
	}

	return res, nil
}

func (intr *GraphInteractor) List(orgID string) (*domain.GraphAnnotations, error) {
	return intr.GraphRepository.List(orgID)
}

func (intr *GraphInteractor) Save(orgID string, annotation *domain.GraphAnnotation) (*domain.GraphAnnotation, error) {
	annotation.ID = domain.NewRandomID()
	return intr.GraphRepository.Save(orgID, annotation)
}

func (intr *GraphInteractor) Update(orgID string, annotation *domain.GraphAnnotation) (*domain.GraphAnnotation, error) {
	return intr.GraphRepository.Update(orgID, annotation)
}

func (intr *GraphInteractor) Delete(orgID, annotationID string) (*domain.GraphAnnotation, error) {
	return intr.GraphRepository.Delete(orgID, annotationID)
}
