package memory

import "github.com/itsubaki/mackerel-api/pkg/domain"

type UserRepository struct {
	Users *domain.Users
}

func NewUserRepository() *UserRepository {
	return &UserRepository{
		Users: &domain.Users{},
	}
}

func (repo *UserRepository) List() (*domain.Users, error) {
	return repo.Users, nil
}

func (repo *UserRepository) Exists(userID string) bool {
	for i := range repo.Users.Users {
		if repo.Users.Users[i].ID == userID {
			return true
		}
	}

	return false
}

func (repo *UserRepository) Delete(userID string) (*domain.User, error) {
	users := make([]domain.User, 0)

	var user domain.User
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
