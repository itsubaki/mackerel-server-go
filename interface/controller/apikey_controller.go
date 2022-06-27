package controller

import (
	"github.com/itsubaki/mackerel-server-go/domain"
	"github.com/itsubaki/mackerel-server-go/interface/database"
	"github.com/itsubaki/mackerel-server-go/usecase"
)

const XAPIKEY string = "X-Api-Key"

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

func (cntr *APIKeyController) APIKey(c Context) (*domain.APIKey, error) {
	return cntr.Interactor.APIKey(c.GetHeader(XAPIKEY))
}
