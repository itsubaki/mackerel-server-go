package mackerel

import "fmt"

type GetServicesInput struct {
}

type GetServicesOutput struct {
	Services []Service `json:"services"`
}

type PostServiceInput struct {
	Name string `json:"name"`
	Memo string `json:"memo"`
}

type PostServiceOutput struct {
	Service
}

type DeleteServiceInput struct {
	ServiceName string
}

type DeleteServiceOutput struct {
	Service
}

type Service struct {
	Name  string   `json:"name"`
	Memo  string   `json:"memo"`
	Roles []string `json:"roles"`
}

type ServiceRepository struct {
	Internal []Service
}

func NewServiceRepository() *ServiceRepository {
	return &ServiceRepository{
		Internal: []Service{},
	}
}

func (repo *ServiceRepository) ExistsByName(serviceName string) bool {
	for i := range repo.Internal {
		if repo.Internal[i].Name == serviceName {
			return true
		}
	}

	return false
}

func (repo *ServiceRepository) FindByName(serviceName string) (Service, error) {
	for i := range repo.Internal {
		if repo.Internal[i].Name == serviceName {
			return repo.Internal[i], nil
		}
	}

	return Service{}, fmt.Errorf("service not found")
}

func (repo *ServiceRepository) FindAll() ([]Service, error) {
	return repo.Internal, nil
}

func (repo *ServiceRepository) Save(s Service) error {
	repo.Internal = append(repo.Internal, s)
	return nil
}

func (repo *ServiceRepository) Delete(serviceName string) error {
	list := []Service{}
	for i := range repo.Internal {
		if repo.Internal[i].Name != serviceName {
			list = append(list, repo.Internal[i])
		}
	}
	repo.Internal = list
	return nil
}
