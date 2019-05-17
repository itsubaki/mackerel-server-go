package controllers

import (
	"github.com/itsubaki/mackerel-api/pkg/interfaces/database"
	"github.com/itsubaki/mackerel-api/pkg/interfaces/memory"
	"github.com/itsubaki/mackerel-api/pkg/usecase"
)

type UserController struct {
	Interactor *usecase.UserInteractor
}

func NewUserController(handler database.SQLHandler) *UserController {
	var repo usecase.UserRepository
	repo = memory.NewUserRepository()
	if handler != nil {
		repo = database.NewUserRepository(handler)
	}

	return &UserController{
		Interactor: &usecase.UserInteractor{
			UserRepository: repo,
		},
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
