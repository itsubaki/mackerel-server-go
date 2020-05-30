package usecase

import "github.com/itsubaki/mackerel-server-go/pkg/domain"

type ServiceRepository interface {
	List(orgID string) (*domain.Services, error)
	Exists(orgID, serviceName string) bool
	Service(orgID, serviceName string) (*domain.Service, error)
	Save(orgID string, service *domain.Service) error
	Delete(orgID, serviceName string) error
}
