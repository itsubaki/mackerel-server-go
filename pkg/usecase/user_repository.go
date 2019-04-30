package usecase

import "github.com/itsubaki/mackerel-api/pkg/domain"

type UserRepository interface {
	FindAll() (domain.Users, error)
	Delete(userId string) error
}
