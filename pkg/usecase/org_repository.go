package usecase

import "github.com/itsubaki/mackerel-server-go/pkg/domain"

type OrgRepository interface {
	Org(orgID string) (*domain.Org, error)
}
