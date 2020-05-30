package usecase

import "github.com/itsubaki/mackerel-server-go/pkg/domain"

type ServiceMetaRepository interface {
	Exists(orgID, serviceName, namespace string) bool
	List(orgID, serviceName string) (*domain.ServiceMetadataList, error)
	Metadata(orgID, serviceName, namespace string) (interface{}, error)
	Save(orgID, serviceName, namespace string, metadata interface{}) (*domain.Success, error)
	Delete(orgID, serviceName, namespace string) (*domain.Success, error)
}
