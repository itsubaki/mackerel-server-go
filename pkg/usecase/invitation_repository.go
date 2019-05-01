package usecase

import "github.com/itsubaki/mackerel-api/pkg/domain"

type InvitationRepository interface {
	List() (*domain.Invitations, error)
	Exists(email string) bool
	Save(inv *domain.Invitation) (*domain.Invitation, error)
	Revoke(email string) (*domain.Success, error)
}
