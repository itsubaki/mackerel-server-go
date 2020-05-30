package usecase

import "github.com/itsubaki/mackerel-server-go/pkg/domain"

type UserRepository interface {
	List(orgID string) (*domain.Users, error)
	Exists(orgID, userID string) bool
	Delete(orgID, userID string) (*domain.User, error)
}
