package usecase

import (
	"github.com/itsubaki/mackerel-api/pkg/domain"
)

type OrgInteractor struct {
	OrgRepository OrgRepository
}

func (s *OrgInteractor) Org(org string) (*domain.Org, error) {
	return s.OrgRepository.Org(org)
}

func (s *OrgInteractor) XAPIKey(xapikey string) (*domain.XAPIKey, error) {
	return s.OrgRepository.XAPIKey(xapikey)
}
