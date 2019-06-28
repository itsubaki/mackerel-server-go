package usecase

import "github.com/itsubaki/mackerel-api/pkg/domain"

type APIKeyInteractor struct {
	APIKeyRepository APIKeyRepository
}

func (s *APIKeyInteractor) APIKey(xapikey string) (*domain.APIKey, error) {
	return s.APIKeyRepository.APIKey(xapikey)
}
