package usecase

import "github.com/itsubaki/mackerel-api/pkg/domain"

type ServiceRepository interface {
	FindAll() (domain.Services, error)
	FindByName(serviceName string) (domain.Service, error)
	ExistsByName(serviceName string) bool
	Save(s domain.Service) error
	Delete(serviceName string) error
}
