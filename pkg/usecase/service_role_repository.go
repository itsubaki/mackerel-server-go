package usecase

import "github.com/itsubaki/mackerel-api/pkg/domain"

type ServiceRoleRepository interface {
	ExistsByName(serviceName, roleName string) bool
	FindByName(serviceName, roleName string) (domain.ServiceRole, error)
	FindAll(serviceName string) (domain.ServiceRoles, error)
	Save(r domain.ServiceRole) error
	Delete(serviceName, roleName string) error
	DeleteAll(serviceName string) error
}
