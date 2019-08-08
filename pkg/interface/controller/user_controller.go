package controller

import (
	"github.com/itsubaki/mackerel-api/pkg/interface/database"
	"github.com/itsubaki/mackerel-api/pkg/usecase"
)

type UserController struct {
	Interactor *usecase.UserInteractor
}

func NewUserController(handler database.SQLHandler) *UserController {
	return &UserController{
		Interactor: &usecase.UserInteractor{
			UserRepository: database.NewUserRepository(handler),
		},
	}
}

func (s *UserController) List(c Context) {
	out, err := s.Interactor.List(
		c.GetString("org_id"),
	)

	doResponse(c, out, err)
}

func (s *UserController) Delete(c Context) {
	out, err := s.Interactor.Delete(
		c.GetString("org_id"),
		c.Param("userId"),
	)

	doResponse(c, out, err)
}
