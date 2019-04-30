package database

import "github.com/itsubaki/mackerel-api/pkg/domain"

type OrgRepository struct {
	Internal domain.Orgs
}

func NewOrgRepository() *OrgRepository {
	return &OrgRepository{
		Internal: domain.Orgs{},
	}
}
