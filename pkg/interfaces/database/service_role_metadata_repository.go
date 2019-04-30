package database

import "github.com/itsubaki/mackerel-api/pkg/domain"

type ServiceRoleMetadataRepository struct {
	Internal domain.ServiceRoleMetadataList
}

func NewServiceRoleMetadataRepository() *ServiceRoleMetadataRepository {
	return &ServiceRoleMetadataRepository{
		Internal: domain.ServiceRoleMetadataList{},
	}
}
