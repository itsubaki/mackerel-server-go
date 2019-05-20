package memory

import (
	"fmt"

	"github.com/itsubaki/mackerel-api/pkg/domain"
)

type OrgRepository struct {
	Orgs []domain.Org
}

func NewOrgRepository() *OrgRepository {
	return &OrgRepository{
		Orgs: []domain.Org{
			{
				Name:    "mackerel-api",
				XAPIKey: "secret",
			},
		},
	}
}

func (repo *OrgRepository) Org(apikey string) (*domain.Org, error) {
	for i := range repo.Orgs {
		if repo.Orgs[i].XAPIKey == apikey {
			return &repo.Orgs[i], nil
		}
	}

	return nil, fmt.Errorf("org not found")
}
