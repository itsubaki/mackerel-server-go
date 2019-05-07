package usecase

import (
	"errors"

	"github.com/itsubaki/mackerel-api/pkg/domain"
)

type OrgInteractor struct {
	OrgRepository OrgRepository
}

func (s *OrgInteractor) Org(apikey string) (*domain.Org, error) {
	org, err := s.OrgRepository.Org(apikey)
	if err != nil {
		return nil, &PermissionDenied{Err{errors.New("permission denied")}}
	}

	return org, nil
}
