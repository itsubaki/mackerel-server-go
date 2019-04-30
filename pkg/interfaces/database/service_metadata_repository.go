package database

import "github.com/itsubaki/mackerel-api/pkg/domain"

type ServiceMetadataRepository struct {
	Internal domain.ServiceMetadataList
}

func NewServiceMetadataRepositoryy() *ServiceMetadataRepository {
	return &ServiceMetadataRepository{
		Internal: domain.ServiceMetadataList{},
	}
}
