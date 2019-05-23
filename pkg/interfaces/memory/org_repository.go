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
				Name: "mackerel-api",
			},
		},
	}
}

func (repo *OrgRepository) Org() (*domain.Org, error) {
	return &repo.Orgs[0], nil
}
