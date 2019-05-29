package controllers

import (
	"github.com/itsubaki/mackerel-api/pkg/domain"
	"github.com/itsubaki/mackerel-api/pkg/interfaces/database"
	"github.com/itsubaki/mackerel-api/pkg/usecase"
)

type AuthController struct {
	Interactor *usecase.AuthInteractor
}

func NewAuthController(handler database.SQLHandler) *AuthController {
	return &AuthController{
		Interactor: &usecase.AuthInteractor{
			AuthRepository: database.NewAuthRepository(handler),
		},
	}
}

func (s *AuthController) XAPIKey(c Context) (*domain.XAPIKey, error) {
	return s.Interactor.XAPIKey(c.GetHeader("X-Api-Key"))
}
