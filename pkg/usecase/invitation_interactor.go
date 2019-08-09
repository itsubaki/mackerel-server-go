package usecase

import (
	"errors"
	"fmt"

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
		return &domain.Success{Success: false}, &InvitationNotFound{Err{errors.New("the specified email has not be sent an invitation")}}
	}

	res, err := s.InvitationRepository.Revoke(orgID, email)
	if err != nil {
		return res, fmt.Errorf("revoke: %v", err)
	}

	return res, nil
}
