package controller

import (
	"net/http"

	"github.com/itsubaki/mackerel-server-go/pkg/domain"
	"github.com/itsubaki/mackerel-server-go/pkg/interface/database"
	"github.com/itsubaki/mackerel-server-go/pkg/usecase"
)

type CheckReportController struct {
	Interactor *usecase.CheckReportInteractor
}

func NewCheckReportController(handler database.SQLHandler) *CheckReportController {
	return &CheckReportController{
		Interactor: &usecase.CheckReportInteractor{
			CheckReportRepository: database.NewCheckReportRepository(handler),
			AlertRepository:       database.NewAlertRepository(handler),
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
		c.GetString("org_id"),
		&in,
	)

	doResponse(c, out, err)
}
