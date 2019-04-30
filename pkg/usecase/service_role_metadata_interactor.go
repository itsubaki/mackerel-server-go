package usecase

import (
	"regexp"

	"github.com/itsubaki/mackerel-api/pkg/domain"
	"github.com/itsubaki/mackerel-api/pkg/interfaces/database"
)

func NewServiceRoleMetadataInteractor() *ServiceRoleMetadataInteractor {
	return &ServiceRoleMetadataInteractor{
		ServiceNameRule:               regexp.MustCompile(`^[a-zA-Z0-9]{1,1}[a-zA-Z0-9_-]{1,62}`),
		ServiceRepository:             database.NewServiceRepository(),
		ServiceRoleMetadataRepository: database.NewServiceRoleMetadataRepository(),
	}
}

type ServiceRoleMetadataInteractor struct {
	ServiceNameRule               *regexp.Regexp
	ServiceRepository             *database.ServiceRepository
	ServiceRoleMetadataRepository *database.ServiceRoleMetadataRepository
}

func (s *ServiceRoleMetadataInteractor) Find(serviceName, roleName, spacename string) (interface{}, error) {
	return nil, nil
}

func (s *ServiceRoleMetadataInteractor) Save(serviceName, roleName, spacename string, metadata interface{}) error {
	return nil
}

func (s *ServiceRoleMetadataInteractor) Delete(serviceName, roleName, spacename string) error {
	return nil
}

func (s *ServiceRoleMetadataInteractor) List(serviceName, roleName string) (domain.ServiceRoleMetadataList, error) {
	return domain.ServiceRoleMetadataList{}, nil
}
