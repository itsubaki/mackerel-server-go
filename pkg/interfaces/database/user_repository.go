package database

import "github.com/itsubaki/mackerel-api/pkg/domain"

type UserRepository struct {
	DB DB
}

func NewUserRepository(db DB) *UserRepository {
	return &UserRepository{
		DB: db,
	}
}

func (repo *UserRepository) List() (*domain.Users, error) {
	return nil, nil
}

func (repo *UserRepository) Delete(userID string) (*domain.User, error) {
	return nil, nil
}
