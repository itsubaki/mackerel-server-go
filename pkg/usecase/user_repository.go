package usecase

import "github.com/itsubaki/mackerel-api/pkg/domain"

type UserRepository interface {
	List(org string) (*domain.Users, error)
	Exists(org, userID string) bool
	Delete(org, userID string) (*domain.User, error)
}
