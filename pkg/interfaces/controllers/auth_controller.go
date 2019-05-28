package controllers

import (
	"github.com/itsubaki/mackerel-api/pkg/domain"
	"github.com/itsubaki/mackerel-api/pkg/interfaces/database"
	"github.com/itsubaki/mackerel-api/pkg/usecase"
)

type AuthController struct {
	Interactor *usecase.OrgInteractor
}

func NewAuthController(handler database.SQLHandler) *AuthController {
	return &AuthController{
		Interactor: &usecase.OrgInteractor{
			OrgRepository: database.NewOrgRepository(handler),
		},
	}
}

func (s *AuthController) Required(c Context) (*domain.XAPIKey, error) {
	return s.Interactor.XAPIKey(c.GetHeader("X-Api-Key"))
}
