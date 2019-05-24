package controllers

import (
	"net/http"

	"github.com/itsubaki/mackerel-api/pkg/interfaces/database"
	"github.com/itsubaki/mackerel-api/pkg/usecase"
)

type AuthController struct {
	Interactor *usecase.OrgInteractor
}

func NewAuthController(handler database.SQLHandler) *AuthController {
	return &AuthController{
		Interactor: &usecase.OrgInteractor{
			OrgRepository: database.NewOrgRepository(handler),
		},
	}
}

func (s *AuthController) Required(c Context) {
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
