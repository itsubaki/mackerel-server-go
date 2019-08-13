package usecase

import "github.com/itsubaki/mackerel-api/pkg/domain"

type RoleRepository interface {
	Exists(orgID, serviceName, roleName string) bool
	List(orgID, serviceName string) (*domain.Roles, error)
	ListWith(orgID string) (map[string][]string, error)
	Role(orgID, serviceName, roleName string) (*domain.Role, error)
	Save(orgID, serviceName string, role *domain.Role) error
	Delete(orgID, serviceName, roleName string) error
}
