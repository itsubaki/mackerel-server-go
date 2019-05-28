package usecase

import "github.com/itsubaki/mackerel-api/pkg/domain"

type OrgRepository interface {
	Org(orgID string) (*domain.Org, error)
	XAPIKey(xapikey string) (*domain.XAPIKey, error)
}
