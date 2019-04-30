package database

import "github.com/itsubaki/mackerel-api/pkg/domain"

type RoleMetadataRepository struct {
	Internal domain.ServiceRoleMetadataList
}

func NewRoleMetadataRepositoryy() *RoleMetadataRepository {
	return &RoleMetadataRepository{
		Internal: domain.ServiceRoleMetadataList{},
	}
}
