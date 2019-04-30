package usecase

import "github.com/itsubaki/mackerel-api/pkg/domain"

type HostRepository interface {
	ExistsByName(hostName string) bool
	FindByID(hostID string) (domain.Host, error)
	FindByName(hostName string) (domain.Host, error)
	FindAll() (domain.Hosts, error)
	Save(host domain.Host) error
	DeleteByID(hostID string) error
	DeleteByName(hostName string) error
}
