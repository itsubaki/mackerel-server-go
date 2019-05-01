package usecase

import (
	"errors"

	"github.com/itsubaki/mackerel-api/pkg/domain"
)

type UserInteractor struct {
	UserRepository UserRepository
}

func (s *UserInteractor) List() (*domain.Users, error) {
	return s.UserRepository.List()
}

func (s *UserInteractor) Delete(userID string) (*domain.User, error) {
	if !s.UserRepository.Exists(userID) {
		return nil, &UserNotFound{Err{errors.New("the <userId> that was designated doesn't belong to the organization")}}
	}

	return s.UserRepository.Delete(userID)
}
