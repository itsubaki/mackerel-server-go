package usecase

import (
	"fmt"

	"github.com/itsubaki/mackerel-server-go/domain"
)

type UserInteractor struct {
	UserRepository UserRepository
}

func (intr *UserInteractor) List(orgID string) (*domain.Users, error) {
	return intr.UserRepository.List(orgID)
}

func (intr *UserInteractor) Delete(orgID, userID string) (*domain.User, error) {
	if !intr.UserRepository.Exists(orgID, userID) {
		return nil, &UserNotFound{Err{fmt.Errorf(fmt.Sprintf("the <%s> that was designated doesn't belong to the organization<%s>", userID, orgID))}}
	}

	return intr.UserRepository.Delete(orgID, userID)
}
