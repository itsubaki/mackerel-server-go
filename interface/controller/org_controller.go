package controller

import (
	"github.com/itsubaki/mackerel-server-go/interface/database"
	"github.com/itsubaki/mackerel-server-go/usecase"
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

func (cntr *OrgController) Org(c Context) {
	out, err := cntr.Interactor.Org(
		c.GetString("org_id"),
	)

	doResponse(c, out, err)
}
