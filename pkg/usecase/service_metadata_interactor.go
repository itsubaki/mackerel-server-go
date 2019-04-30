package usecase

import (
	"github.com/itsubaki/mackerel-api/pkg/domain"
	"github.com/itsubaki/mackerel-api/pkg/interfaces/database"
)

func NewServiceMetadataInteractor() *ServiceMetadataInteractor {
	return &ServiceMetadataInteractor{
		ServiceMetadataRepository: database.NewServiceMetadataRepository(),
	}
}

type ServiceMetadataInteractor struct {
	ServiceMetadataRepository *database.ServiceMetadataRepository
}

func (s *ServiceMetadataInteractor) MetadataList() (domain.ServiceMetadataList, error) {
	return domain.ServiceMetadataList{}, nil
}

func (s *ServiceMetadataInteractor) Find(serviceName, namespace string) (interface{}, error) {
	return nil, nil
}

func (s *ServiceMetadataInteractor) Save(serviceName, namespace string, metadata interface{}) error {
	return nil
}

func (s *ServiceMetadataInteractor) Delete(serviceName, namespace string) error {
	return nil
}
