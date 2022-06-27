package usecase

import "github.com/itsubaki/mackerel-server-go/domain"

type InvitationRepository interface {
	List(orgID string) (*domain.Invitations, error)
	Exists(orgID, email string) bool
	Save(orgID string, inv *domain.Invitation) (*domain.Invitation, error)
	Revoke(orgID, email string) (*domain.Success, error)
}
