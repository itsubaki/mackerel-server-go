package controllers

import (
	"net/http"

	"github.com/itsubaki/mackerel-api/pkg/domain"
	"github.com/itsubaki/mackerel-api/pkg/interfaces/database"
	"github.com/itsubaki/mackerel-api/pkg/usecase"
)

type CheckReportController struct {
	Interactor *usecase.CheckReportInteractor
}

func NewCheckReportController(sqlHandler database.SQLHandler) *CheckReportController {
	return &CheckReportController{}
}

func (s *CheckReportController) Save(c Context) {
	var in domain.CheckReports
	if err := c.BindJSON(in); err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	out, err := s.Interactor.Save(&in)
	doResponse(c, out, err)
}
