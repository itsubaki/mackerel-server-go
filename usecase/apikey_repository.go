package usecase

import "github.com/itsubaki/mackerel-server-go/domain"

type APIKeyRepository interface {
	Save(orgID, name, apikey string, write bool) (*domain.APIKey, error)
	APIKey(apikey string) (*domain.APIKey, error)
}
