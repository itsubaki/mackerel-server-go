package usecase

import "github.com/itsubaki/mackerel-server-go/pkg/domain"

type APIKeyRepository interface {
	Save(orgID, name string, write bool) (*domain.APIKey, error)
	List(orgID string) ([]domain.APIKey, error)
	APIKey(apikey string) (*domain.APIKey, error)
}
