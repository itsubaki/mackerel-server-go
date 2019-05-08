package controllers

import (
	"net/http"

	"github.com/itsubaki/mackerel-api/pkg/domain"
	"github.com/itsubaki/mackerel-api/pkg/interfaces/database"
	"github.com/itsubaki/mackerel-api/pkg/interfaces/memory"
	"github.com/itsubaki/mackerel-api/pkg/usecase"
)

type InvitationController struct {
	Interactor *usecase.InvitationInteractor
}

func NewInvitationController(handler database.DB) *InvitationController {
	return &InvitationController{
		Interactor: &usecase.InvitationInteractor{
			InvitationRepository: memory.NewInvitationRepository(),
		},
	}
}

func (s *InvitationController) List(c Context) {
	out, err := s.Interactor.List()
	doResponse(c, out, err)
}

func (s *InvitationController) Save(c Context) {
	var in domain.Invitation
	if err := c.BindJSON(&in); err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	out, err := s.Interactor.Save(&in)
	doResponse(c, out, err)
}

func (s *InvitationController) Revoke(c Context) {
	var in domain.Revoke
	if err := c.BindJSON(&in); err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	out, err := s.Interactor.Revoke(in.EMail)
	doResponse(c, out, err)
}
