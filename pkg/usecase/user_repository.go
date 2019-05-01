package usecase

import "github.com/itsubaki/mackerel-api/pkg/domain"

type UserRepository interface {
	List() (*domain.Users, error)
	Exists(userID string) bool
	Delete(userID string) (*domain.User, error)
}
