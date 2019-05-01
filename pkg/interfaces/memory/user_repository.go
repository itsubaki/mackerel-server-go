package memory

import "github.com/itsubaki/mackerel-api/pkg/domain"

type UserRepository struct {
	Users *domain.Users
}

func (repo *UserRepository) List() (*domain.Users, error) {
	return repo.Users, nil
}

func (repo *UserRepository) Delete(userID string) (*domain.User, error) {
	var user domain.User

	users := []domain.User{}
	for i := range repo.Users.Users {
		if repo.Users.Users[i].ID == userID {
			user = repo.Users.Users[i]
			continue
		}
		users = append(users, repo.Users.Users[i])
	}

	repo.Users.Users = users

	return &user, nil
}
