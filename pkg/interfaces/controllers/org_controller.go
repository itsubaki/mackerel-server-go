package controllers

import (
	"net/http"

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

func (s *OrgController) AuthRequired(c Context) {
	xapikey, err := s.Interactor.XAPIKey(c.GetHeader("X-Api-Key"))
	if err != nil {
		c.Status(http.StatusForbidden)
		c.Abort()
		return
	}

	if c.GetString("Method") != http.MethodGet && !xapikey.Write {
		c.Status(http.StatusForbidden)
		c.Abort()
		return
	}

	c.Set("org", xapikey.Org)
	c.Next()
}

func (s *OrgController) Org(c Context) {
	out, err := s.Interactor.Org(c.GetString("org"))
	doResponse(c, out, err)
}
