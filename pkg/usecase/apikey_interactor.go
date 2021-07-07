package usecase

import "github.com/itsubaki/mackerel-server-go/pkg/domain"

type APIKeyInteractor struct {
	APIKeyRepository APIKeyRepository
}

func (intr *APIKeyInteractor) APIKey(xapikey string) (*domain.APIKey, error) {
	return intr.APIKeyRepository.APIKey(xapikey)
}
