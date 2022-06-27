package controller

import (
	"net/http"

	"github.com/itsubaki/mackerel-server-go/domain"
	"github.com/itsubaki/mackerel-server-go/interface/database"
	"github.com/itsubaki/mackerel-server-go/usecase"
)

type InvitationController struct {
	Interactor *usecase.InvitationInteractor
}

func NewInvitationController(handler database.SQLHandler) *InvitationController {
	return &InvitationController{
		Interactor: &usecase.InvitationInteractor{
			InvitationRepository: database.NewInvitationRepository(handler),
		},
	}
}

func (cntr *InvitationController) List(c Context) {
	out, err := cntr.Interactor.List(
		c.GetString("org_id"),
	)

	doResponse(c, out, err)
}

func (cntr *InvitationController) Save(c Context) {
	var in domain.Invitation
	if err := c.BindJSON(&in); err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	out, err := cntr.Interactor.Save(
		c.GetString("org_id"),
		&in,
	)

	doResponse(c, out, err)
}

func (cntr *InvitationController) Revoke(c Context) {
	var in domain.Revoke
	if err := c.BindJSON(&in); err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	out, err := cntr.Interactor.Revoke(
		c.GetString("org_id"),
		in.EMail,
	)

	doResponse(c, out, err)
}
