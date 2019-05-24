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
	var repo usecase.OrgRepository
	repo = memory.NewOrgRepository()
	if handler != nil {
		repo = database.NewOrgRepository(handler)
	}

	return &OrgController{
		Interactor: &usecase.OrgInteractor{
			OrgRepository: repo,
		},
	}
}

func (s *OrgController) Org(c Context) {
	out, err := s.Interactor.Org(c.GetString("org"))
	doResponse(c, out, err)
}
