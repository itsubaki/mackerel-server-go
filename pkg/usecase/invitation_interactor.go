package usecase

import (
	"errors"

	"github.com/itsubaki/mackerel-api/pkg/domain"
)

type InvitationInteractor struct {
	InvitationRepository InvitationRepository
}

func (s *InvitationInteractor) List() (*domain.Invitations, error) {
	return s.InvitationRepository.List()
}

func (s *InvitationInteractor) Save(inv *domain.Invitation) (*domain.Invitation, error) {
	return s.InvitationRepository.Save(inv)
}

func (s *InvitationInteractor) Revoke(email string) (*domain.Success, error) {
	if !s.InvitationRepository.Exists(email) {
		return nil, &InvitationNotFound{Err{errors.New("the specified email has not be sent an invitation")}}
	}

	return s.InvitationRepository.Revoke(email)
}
