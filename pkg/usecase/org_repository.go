package usecase

import "github.com/itsubaki/mackerel-api/pkg/domain"

type OrgRepository interface {
	Org(apikey string) (*domain.Org, error)
}
