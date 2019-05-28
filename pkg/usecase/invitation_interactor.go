package usecase

import (
	"errors"

	"github.com/itsubaki/mackerel-api/pkg/domain"
)

type InvitationInteractor struct {
	InvitationRepository InvitationRepository
}

func (s *InvitationInteractor) List(orgID string) (*domain.Invitations, error) {
	return s.InvitationRepository.List(orgID)
}

func (s *InvitationInteractor) Save(orgID string, inv *domain.Invitation) (*domain.Invitation, error) {
	return s.InvitationRepository.Save(orgID, inv)
}

func (s *InvitationInteractor) Revoke(orgID, email string) (*domain.Success, error) {
	if !s.InvitationRepository.Exists(orgID, email) {
		return nil, &InvitationNotFound{Err{errors.New("the specified email has not be sent an invitation")}}
	}

	return s.InvitationRepository.Revoke(orgID, email)
}
