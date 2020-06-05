package usecase

import "github.com/itsubaki/mackerel-server-go/pkg/domain"

type OrgRepository interface {
	Save(orgID, name string) (*domain.Org, error)
	Org(orgID string) (*domain.Org, error)
}
