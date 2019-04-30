package usecase

import (
	"regexp"

	"github.com/itsubaki/mackerel-api/pkg/domain"
)

type ServiceRoleMetadataInteractor struct {
	ServiceNameRule               *regexp.Regexp
	ServiceRepository             ServiceRepository
	ServiceRoleMetadataRepository ServiceRoleMetadataRepository
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
