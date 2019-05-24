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

func NewCheckReportController(handler database.SQLHandler) *CheckReportController {
	return &CheckReportController{
		Interactor: &usecase.CheckReportInteractor{
			CheckReportRepository: database.NewCheckReportRepository(handler),
		},
	}
}

func (s *CheckReportController) Save(c Context) {
	var in domain.CheckReports
	if err := c.BindJSON(&in); err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	out, err := s.Interactor.Save(
		c.GetString("org"),
		&in,
	)

	doResponse(c, out, err)
}
