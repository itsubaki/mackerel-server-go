package usecase

import (
	"errors"
	"fmt"

	"github.com/itsubaki/mackerel-server-go/pkg/domain"
)

type UserInteractor struct {
	UserRepository UserRepository
}

func (intr *UserInteractor) List(orgID string) (*domain.Users, error) {
	return intr.UserRepository.List(orgID)
}

func (intr *UserInteractor) Delete(orgID, userID string) (*domain.User, error) {
	if !intr.UserRepository.Exists(orgID, userID) {
		return nil, &UserNotFound{Err{errors.New(fmt.Sprintf("the <%s> that was designated doesn't belong to the organization<%s>", userID, orgID))}}
	}

	return intr.UserRepository.Delete(orgID, userID)
}
