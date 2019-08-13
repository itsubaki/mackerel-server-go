package usecase

import "github.com/itsubaki/mackerel-api/pkg/domain"

type HostRepository interface {
	List(orgID string) (*domain.Hosts, error)
	ActiveList(orgID string) (*domain.Hosts, error)
	Save(orgID string, host *domain.Host) (*domain.HostID, error)
	Host(orgID, hostID string) (*domain.Host, error)
	Exists(orgID, hostID string) bool
	Status(orgID, hostID, status string) (*domain.Success, error)
	Retire(orgID, hostID string, retire *domain.HostRetire) (*domain.Success, error)
	SaveRoleFullNames(orgID, hostID string, names *domain.RoleFullNames) (*domain.Success, error)
}
