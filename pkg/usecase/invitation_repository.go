package usecase

import "github.com/itsubaki/mackerel-api/pkg/domain"

type InvitationRepository interface {
	List(org string) (*domain.Invitations, error)
	Exists(org, email string) bool
	Save(org string, inv *domain.Invitation) (*domain.Invitation, error)
	Revoke(org, email string) (*domain.Success, error)
}
