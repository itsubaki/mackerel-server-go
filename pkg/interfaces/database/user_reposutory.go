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
