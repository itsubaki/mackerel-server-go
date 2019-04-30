package database

import "github.com/itsubaki/mackerel-api/pkg/domain"

type UserRepository struct {
	Internal domain.Users
}

func NewUserRepository() *UserRepository {
	return &UserRepository{
		Internal: domain.Users{},
	}
}

func (repo *UserRepository) FindAll() (domain.Users, error) {
	return repo.Internal, nil
}

func (repo *UserRepository) Delete(userId string) error {
	list := domain.Users{}
	for i := range repo.Internal {
		if repo.Internal[i].ID != userId {
			list = append(list, repo.Internal[i])
		}
	}
	repo.Internal = list

	return nil
}
