package usecase

import "github.com/itsubaki/mackerel-server-go/domain"

type HostMetaRepository interface {
	Exists(orgID, hostID, namespace string) bool
	List(orgID, hostID string) (*domain.HostMetadataList, error)
	Metadata(orgID, hostID, namespace string) (any, error)
	Save(orgID, hostID, namespace string, metadata any) (*domain.Success, error)
	Delete(orgID, hostID, namespace string) (*domain.Success, error)
}
