package controller

import (
	"github.com/itsubaki/mackerel-server-go/pkg/interface/database"
	"github.com/itsubaki/mackerel-server-go/pkg/usecase"
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

func (cntr *UserController) List(c Context) {
	out, err := cntr.Interactor.List(
		c.GetString("org_id"),
	)

	doResponse(c, out, err)
}

func (cntr *UserController) Delete(c Context) {
	out, err := cntr.Interactor.Delete(
		c.GetString("org_id"),
		c.Param("userId"),
	)

	doResponse(c, out, err)
}
