package usecase

import "github.com/itsubaki/mackerel-api/pkg/domain"

type OrgRepository interface {
	Org() (*domain.Org, error)
}
