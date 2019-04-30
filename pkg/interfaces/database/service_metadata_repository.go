package database

import "github.com/itsubaki/mackerel-api/pkg/domain"

type ServiceMetadataRepository struct {
	Internal domain.ServiceMetadataList
}

func NewServiceMetadataRepository() *ServiceMetadataRepository {
	return &ServiceMetadataRepository{
		Internal: domain.ServiceMetadataList{},
	}
}

func (repo *ServiceMetadataRepository) FindAll() (domain.ServiceMetadataList, error) {
	return repo.Internal, nil
}
