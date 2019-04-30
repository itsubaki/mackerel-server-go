package usecase

import "github.com/itsubaki/mackerel-api/pkg/domain"

type HostMetadataRepository interface {
	Find(hostID, namespace string) (domain.HostMetadataList, error)
	FindByID(hostID string) (domain.HostMetadataList, error)
	Save(metadata domain.HostMetadata) error
	Delete(hostID, namespace string) error
}
