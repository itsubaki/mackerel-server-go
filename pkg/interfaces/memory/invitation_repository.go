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

func (repo *InvitationRepository) List() (*domain.Invitations, error) {
	return repo.Invitations, nil
}

func (repo *InvitationRepository) Exists(email string) bool {
	for i := range repo.Invitations.Invitations {
		if repo.Invitations.Invitations[i].EMail == email {
			return true
		}
	}

	return false
}

func (repo *InvitationRepository) Save(inv *domain.Invitation) (*domain.Invitation, error) {
	inv.ExpiresAt = time.Now().Unix() + 604800
	repo.Invitations.Invitations = append(repo.Invitations.Invitations, *inv)
	return inv, nil
}

func (repo *InvitationRepository) Revoke(email string) (*domain.Success, error) {
	list := []domain.Invitation{}
	for i := range repo.Invitations.Invitations {
		if repo.Invitations.Invitations[i].EMail == email {
			continue
		}
		list = append(list, repo.Invitations.Invitations[i])
	}

	repo.Invitations.Invitations = list
	return &domain.Success{Success: true}, nil
}
