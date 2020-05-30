package controller

import (
	"net/http"

	"github.com/itsubaki/mackerel-server-go/pkg/domain"
	"github.com/itsubaki/mackerel-server-go/pkg/interface/database"
	"github.com/itsubaki/mackerel-server-go/pkg/usecase"
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

func (s *InvitationController) List(c Context) {
	out, err := s.Interactor.List(
		c.GetString("org_id"),
	)

	doResponse(c, out, err)
}

func (s *InvitationController) Save(c Context) {
	var in domain.Invitation
	if err := c.BindJSON(&in); err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	out, err := s.Interactor.Save(
		c.GetString("org_id"),
		&in,
	)

	doResponse(c, out, err)
}

func (s *InvitationController) Revoke(c Context) {
	var in domain.Revoke
	if err := c.BindJSON(&in); err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	out, err := s.Interactor.Revoke(
		c.GetString("org_id"),
		in.EMail,
	)

	doResponse(c, out, err)
}
