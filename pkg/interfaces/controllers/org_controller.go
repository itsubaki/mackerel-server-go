package controllers

import "github.com/itsubaki/mackerel-api/pkg/interfaces/database"

type OrgController struct {
}

func NewOrgController(sqlHandler database.SQLHandler) *OrgController {
	return &OrgController{}
}

func (s *OrgController) Org(c Context) {
}
