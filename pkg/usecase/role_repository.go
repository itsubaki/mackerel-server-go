package usecase

import "github.com/itsubaki/mackerel-server-go/pkg/domain"

type RoleRepository interface {
	Exists(orgID, serviceName, roleName string) bool
	List(orgID string) (map[string][]string, error)
	ListWith(orgID, serviceName string) (*domain.Roles, error)
	Role(orgID, serviceName, roleName string) (*domain.Role, error)
	Save(orgID, serviceName string, role *domain.Role) error
	Delete(orgID, serviceName, roleName string) error
}
