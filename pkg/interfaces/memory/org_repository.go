package memory

import (
	"github.com/itsubaki/mackerel-api/pkg/domain"
)

type OrgRepository struct {
	Orgs []domain.Org
}

func NewOrgRepository() *OrgRepository {
	return &OrgRepository{
		Orgs: []domain.Org{
			{
				Name: "default",
			},
		},
	}
}

func (repo *OrgRepository) Org(org string) (*domain.Org, error) {
	return &repo.Orgs[0], nil
}

func (repo *OrgRepository) XAPIKey(xapikey string) (*domain.XAPIKey, error) {
	return &domain.XAPIKey{}, nil
}
