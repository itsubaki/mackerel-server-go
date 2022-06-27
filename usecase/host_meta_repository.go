package usecase

import "github.com/itsubaki/mackerel-server-go/domain"

type HostMetaRepository interface {
	Exists(orgID, hostID, namespace string) bool
	List(orgID, hostID string) (*domain.HostMetadataList, error)
	Metadata(orgID, hostID, namespace string) (interface{}, error)
	Save(orgID, hostID, namespace string, metadata interface{}) (*domain.Success, error)
	Delete(orgID, hostID, namespace string) (*domain.Success, error)
}
