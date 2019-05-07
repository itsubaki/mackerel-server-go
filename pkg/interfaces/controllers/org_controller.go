package controllers

import (
	"github.com/itsubaki/mackerel-api/pkg/interfaces/database"
	"github.com/itsubaki/mackerel-api/pkg/interfaces/memory"
	"github.com/itsubaki/mackerel-api/pkg/usecase"
)

type OrgController struct {
	Interactor *usecase.OrgInteractor
}

func NewOrgController(handler database.SQLHandler) *OrgController {
	return &OrgController{
		Interactor: &usecase.OrgInteractor{
			OrgRepository: memory.NewOrgRepository(),
		},
	}
}

func (s *OrgController) Org(c Context) {
	out, err := s.Interactor.Org(c.GetHeader("X-Api-Key"))
	doResponse(c, out, err)
}
