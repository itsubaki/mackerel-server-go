package usecase

import "github.com/itsubaki/mackerel-server-go/domain"

type RoleMetaRepository interface {
	Exists(orgID, serviceName, roleName, namespace string) bool
	List(orgID, serviceName, roleName string) (*domain.RoleMetadataList, error)
	Metadata(orgID, serviceName, roleName, namespace string) (interface{}, error)
	Save(orgID, serviceName, roleName, namespace string, metadata interface{}) (*domain.Success, error)
	Delete(orgID, serviceName, roleName, namespace string) (*domain.Success, error)
}
