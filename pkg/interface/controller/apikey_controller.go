package controller

import (
	"github.com/itsubaki/mackerel-server-go/pkg/domain"
	"github.com/itsubaki/mackerel-server-go/pkg/interface/database"
	"github.com/itsubaki/mackerel-server-go/pkg/usecase"
)

type APIKeyController struct {
	Interactor *usecase.APIKeyInteractor
}

func NewAPIKeyController(handler database.SQLHandler) *APIKeyController {
	return &APIKeyController{
		Interactor: &usecase.APIKeyInteractor{
			APIKeyRepository: database.NewAPIKeyRepository(handler),
		},
	}
}

func (s *APIKeyController) APIKey(c Context) (*domain.APIKey, error) {
	return s.Interactor.APIKey(c.GetHeader("X-Api-Key"))
}
