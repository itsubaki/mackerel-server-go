package usecase

import "github.com/itsubaki/mackerel-api/pkg/domain"

type AuthInteractor struct {
	AuthRepository AuthRepository
}

func (s *AuthInteractor) XAPIKey(xapikey string) (*domain.XAPIKey, error) {
	return s.AuthRepository.XAPIKey(xapikey)
}
