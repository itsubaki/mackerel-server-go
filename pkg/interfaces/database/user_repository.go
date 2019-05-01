package database

import "github.com/itsubaki/mackerel-api/pkg/domain"

type UserRepository struct {
	SQLHandler SQLHandler
}

func (repo *UserRepository) List() (*domain.Users, error) {
	return nil, nil
}

func (repo *UserRepository) Delete(userID string) (*domain.User, error) {
	return nil, nil
}
