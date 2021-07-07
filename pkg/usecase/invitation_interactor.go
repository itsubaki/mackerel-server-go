package usecase

import (
	"errors"
	"fmt"

	"github.com/itsubaki/mackerel-server-go/pkg/domain"
)

type InvitationInteractor struct {
	InvitationRepository InvitationRepository
}

func (intr *InvitationInteractor) List(orgID string) (*domain.Invitations, error) {
	return intr.InvitationRepository.List(orgID)
}

func (intr *InvitationInteractor) Save(orgID string, inv *domain.Invitation) (*domain.Invitation, error) {
	return intr.InvitationRepository.Save(orgID, inv)
}

func (intr *InvitationInteractor) Revoke(orgID, email string) (*domain.Success, error) {
	if !intr.InvitationRepository.Exists(orgID, email) {
		return &domain.Success{Success: false}, &InvitationNotFound{Err{errors.New("the specified email has not be sent an invitation")}}
	}

	res, err := intr.InvitationRepository.Revoke(orgID, email)
	if err != nil {
		return res, fmt.Errorf("revoke: %v", err)
	}

	return res, nil
}
