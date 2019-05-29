package usecase

import "github.com/itsubaki/mackerel-api/pkg/domain"

type AuthRepository interface {
	XAPIKey(xapikey string) (*domain.XAPIKey, error)
}
