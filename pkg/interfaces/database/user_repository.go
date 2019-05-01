package database

import "github.com/itsubaki/mackerel-api/pkg/domain"

type UserRepository struct {
	Users domain.Users
}

func (repo *UserRepository) FindAll() (domain.Users, error) {
	return repo.Users, nil
}

func (repo *UserRepository) Delete(userId string) error {
	list := domain.Users{}
	for i := range repo.Users {
		if repo.Users[i].ID != userId {
			list = append(list, repo.Users[i])
		}
	}
	repo.Users = list

	return nil
}
