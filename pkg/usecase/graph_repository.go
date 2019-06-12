package usecase

import "github.com/itsubaki/mackerel-api/pkg/domain"

type GraphRepository interface {
	SaveDef(orgID string, g []domain.GraphDef) (*domain.Success, error)

	List(orgID string) (*domain.GraphAnnotations, error)
	Save(orgID string, annotation *domain.GraphAnnotation) (*domain.GraphAnnotation, error)
	Update(orgID string, annotation *domain.GraphAnnotation) (*domain.GraphAnnotation, error)
	Delete(orgID, annotationID string) (*domain.GraphAnnotation, error)
}
