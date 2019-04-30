package database

import "github.com/itsubaki/mackerel-api/pkg/domain"

type HostMetadataRepository struct {
	Internal domain.HostMetadataList
}

func NewHostMetadataRepository() *HostMetadataRepository {
	return &HostMetadataRepository{
		Internal: domain.HostMetadataList{},
	}
}
