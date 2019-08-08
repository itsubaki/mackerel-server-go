package usecase

import (
	"errors"
	"fmt"

	"github.com/itsubaki/mackerel-api/pkg/domain"
)

type UserInteractor struct {
	UserRepository UserRepository
}

func (s *UserInteractor) List(orgID string) (*domain.Users, error) {
	return s.UserRepository.List(orgID)
}

func (s *UserInteractor) Delete(orgID, userID string) (*domain.User, error) {
	if !s.UserRepository.Exists(orgID, userID) {
		return nil, &UserNotFound{Err{errors.New(fmt.Sprintf("the <%s> that was designated doesn't belong to the organization<%s>", userID, orgID))}}
	}

	return s.UserRepository.Delete(orgID, userID)
}
