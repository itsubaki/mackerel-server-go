package controllers

import (
	"github.com/itsubaki/mackerel-api/pkg/interfaces/database"
	"github.com/itsubaki/mackerel-api/pkg/usecase"
)

type UserController struct {
	Interactor *usecase.UserInteractor
}

func NewUserController(sqlHandler database.SQLHandler) *UserController {
	return &UserController{}
}

func (s *UserController) List(c Context) {
	out, err := s.Interactor.List()
	doResponse(c, out, err)
}

func (s *UserController) Delete(c Context) {
	out, err := s.Interactor.Delete(c.Param("userId"))
	doResponse(c, out, err)
}
