package controllers

import (
	"github.com/itsubaki/mackerel-api/pkg/interfaces/database"
	"github.com/itsubaki/mackerel-api/pkg/usecase"
)

type UserController struct {
	Interactor *usecase.UserInteractor
}

func NewUserController(handler database.SQLHandler) *UserController {
	return &UserController{
		Interactor: &usecase.UserInteractor{},
	}
}

func (s *UserController) List(c Context) {
	out, err := s.Interactor.List()
	doResponse(c, out, err)
}

func (s *UserController) Delete(c Context) {
	out, err := s.Interactor.Delete(
		c.Param("userId"),
	)

	doResponse(c, out, err)
}
