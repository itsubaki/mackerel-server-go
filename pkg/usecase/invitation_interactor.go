package usecase

import (
	"errors"

	"github.com/itsubaki/mackerel-api/pkg/domain"
)

type InvitationInteractor struct {
	InvitationRepository InvitationRepository
}

func (s *InvitationInteractor) List(org string) (*domain.Invitations, error) {
	return s.InvitationRepository.List(org)
}

func (s *InvitationInteractor) Save(org string, inv *domain.Invitation) (*domain.Invitation, error) {
	return s.InvitationRepository.Save(org, inv)
}

func (s *InvitationInteractor) Revoke(org, email string) (*domain.Success, error) {
	if !s.InvitationRepository.Exists(org, email) {
		return nil, &InvitationNotFound{Err{errors.New("the specified email has not be sent an invitation")}}
	}

	return s.InvitationRepository.Revoke(org, email)
}
