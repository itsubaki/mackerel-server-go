package mackerel

import "fmt"

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

func (repo *ServiceRepository) Find(s Service) (Service, error) {
	for i := range repo.Internal {
		if repo.Internal[i].Name == s.Name {
			return repo.Internal[i], nil
		}
	}

	return Service{}, fmt.Errorf("service not found")
}

func (repo *ServiceRepository) FindAll() ([]Service, error) {
	return repo.Internal, nil
}

func (repo *ServiceRepository) Insert(s Service) error {
	repo.Internal = append(repo.Internal, s)
	return nil
}

func (repo *ServiceRepository) Delete(s Service) error {
	list := []Service{}
	for i := range repo.Internal {
		if repo.Internal[i].Name != s.Name {
			list = append(list, repo.Internal[i])
		}
	}
	repo.Internal = list
	return nil
}

type GetServicesOutput struct {
	Services []Service `json:"services"`
}

type PostServiceInput struct {
	Name string `json:"name"`
	Memo string `json:"memo"`
}

type PostServiceOutput Service

type DeleteServiceInput struct {
	ServiceName string
}

type DeleteServiceOutput Service
