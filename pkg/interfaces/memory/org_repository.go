package memory

import "github.com/itsubaki/mackerel-api/pkg/domain"

type OrgRepository struct {
	Orgs []domain.Org
}

func NewOrgRepository() *OrgRepository {
	return &OrgRepository{
		Orgs: []domain.Org{},
	}
}

func (repo *OrgRepository) Org(apikey string) (*domain.Org, error) {
	return &domain.Org{}, nil
}
