package usecase

import (
	"github.com/itsubaki/mackerel-api/pkg/domain"
)

type OrgInteractor struct {
	OrgRepository OrgRepository
}

func (s *OrgInteractor) Org(orgID string) (*domain.Org, error) {
	return s.OrgRepository.Org(orgID)
}
