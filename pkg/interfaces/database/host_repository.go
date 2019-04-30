package database

import (
	"fmt"

	"github.com/itsubaki/mackerel-api/pkg/domain"
)

type HostRepository struct {
	Internal domain.Hosts
}

func NewHostRepository() *HostRepository {
	return &HostRepository{
		Internal: domain.Hosts{},
	}
}

func (repo *HostRepository) ExistsByName(hostName string) bool {
	for i := range repo.Internal {
		if repo.Internal[i].Name == hostName {
			return true
		}
	}

	return false
}

func (repo *HostRepository) FindByName(hostName string) (domain.Host, error) {
	for i := range repo.Internal {
		if repo.Internal[i].Name == hostName {
			return repo.Internal[i], nil
		}
	}

	return domain.Host{}, fmt.Errorf("host not found")
}

func (repo *HostRepository) FindAll() (domain.Hosts, error) {
	return repo.Internal, nil
}

func (repo *HostRepository) Save(h domain.Host) error {
	repo.Internal = append(repo.Internal, h)
	return nil
}

func (repo *HostRepository) Delete(hostName string) error {
	list := domain.Hosts{}
	for i := range repo.Internal {
		if repo.Internal[i].Name != hostName {
			list = append(list, repo.Internal[i])
		}
	}
	repo.Internal = list
	return nil
}
