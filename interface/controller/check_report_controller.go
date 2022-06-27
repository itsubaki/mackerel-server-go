package controller

import (
	"net/http"

	"github.com/itsubaki/mackerel-server-go/domain"
	"github.com/itsubaki/mackerel-server-go/interface/database"
	"github.com/itsubaki/mackerel-server-go/usecase"
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

func (cntr *CheckReportController) Save(c Context) {
	var in domain.CheckReports
	if err := c.BindJSON(&in); err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	out, err := cntr.Interactor.Save(
		c.GetString("org_id"),
		&in,
	)

	doResponse(c, out, err)
}
