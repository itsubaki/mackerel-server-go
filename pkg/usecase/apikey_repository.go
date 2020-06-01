package usecase

import "github.com/itsubaki/mackerel-server-go/pkg/domain"

type APIKeyRepository interface {
	APIKey(apikey string) (*domain.APIKey, error)
}
