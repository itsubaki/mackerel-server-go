package usecase

import "github.com/itsubaki/mackerel-api/pkg/domain"

type APIKeyRepository interface {
	Save(orgID, name string, write bool) (*domain.APIKey, error)
	APIKey(apikey string) (*domain.APIKey, error)
}
