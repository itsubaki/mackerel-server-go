package controllers

import (
	"github.com/itsubaki/mackerel-api/pkg/interfaces/database"
	"github.com/itsubaki/mackerel-api/pkg/usecase"
)

type OrgController struct {
	Interactor *usecase.OrgInteractor
}

func NewOrgController(handler database.SQLHandler) *OrgController {
	return &OrgController{
		Interactor: &usecase.OrgInteractor{
			OrgRepository: database.NewOrgRepository(handler),
		},
	}
}

func (s *OrgController) Org(c Context) {
	out, err := s.Interactor.Org(
		c.GetString("org"),
	)

	doResponse(c, out, err)
}
