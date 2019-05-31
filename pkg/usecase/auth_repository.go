package usecase

import "github.com/itsubaki/mackerel-api/pkg/domain"

type AuthRepository interface {
	Save(orgID, name string, write bool) (*domain.XAPIKey, error)
	XAPIKey(xapikey string) (*domain.XAPIKey, error)
}
