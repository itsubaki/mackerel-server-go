package usecase

import (
	"github.com/itsubaki/mackerel-server-go/domain"
)

type OrgInteractor struct {
	OrgRepository OrgRepository
}

func (intr *OrgInteractor) Org(orgID string) (*domain.Org, error) {
	return intr.OrgRepository.Org(orgID)
}
