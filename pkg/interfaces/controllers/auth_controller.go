package controllers

import (
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

func (s *AuthController) Required(c Context) (string, bool, error) {
	key, err := s.Interactor.XAPIKey(c.GetHeader("X-Api-Key"))
	if err != nil {
		return "", false, err
	}

	return key.Org, key.Write, nil
}
