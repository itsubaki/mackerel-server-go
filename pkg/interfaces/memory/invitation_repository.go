package memory

import (
	"time"

	"github.com/itsubaki/mackerel-api/pkg/domain"
)

type InvitationRepository struct {
	Invitations *domain.Invitations
}

func NewInvitationRepository() *InvitationRepository {
	return &InvitationRepository{
		Invitations: &domain.Invitations{
			Invitations: []domain.Invitation{},
		},
	}
}

func (repo *InvitationRepository) List(org string) (*domain.Invitations, error) {
	return repo.Invitations, nil
}

func (repo *InvitationRepository) Exists(org, email string) bool {
	for i := range repo.Invitations.Invitations {
		if repo.Invitations.Invitations[i].EMail == email {
			return true
		}
	}

	return false
}

func (repo *InvitationRepository) Save(org string, inv *domain.Invitation) (*domain.Invitation, error) {
	inv.ExpiresAt = time.Now().Unix() + 604800
	repo.Invitations.Invitations = append(repo.Invitations.Invitations, *inv)
	return inv, nil
}

func (repo *InvitationRepository) Revoke(org, email string) (*domain.Success, error) {
	list := make([]domain.Invitation, 0)
	for i := range repo.Invitations.Invitations {
		if repo.Invitations.Invitations[i].EMail == email {
			continue
		}
		list = append(list, repo.Invitations.Invitations[i])
	}

	repo.Invitations.Invitations = list
	return &domain.Success{Success: true}, nil
}
