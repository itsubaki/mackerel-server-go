package controllers

import (
	"net/http"

	"github.com/itsubaki/mackerel-api/pkg/domain"
	"github.com/itsubaki/mackerel-api/pkg/interfaces/database"
	"github.com/itsubaki/mackerel-api/pkg/interfaces/memory"
	"github.com/itsubaki/mackerel-api/pkg/usecase"
)

type CheckReportController struct {
	Interactor *usecase.CheckReportInteractor
}

func NewCheckReportController(handler database.SQLHandler) *CheckReportController {
	var repo usecase.CheckReportRepository
	repo = memory.NewCheckReportRepository()
	if handler != nil {
		repo = database.NewCheckReportRepository(handler)
	}

	return &CheckReportController{
		Interactor: &usecase.CheckReportInteractor{
			CheckReportRepository: repo,
		},
	}
}

func (s *CheckReportController) Save(c Context) {
	var in domain.CheckReports
	if err := c.BindJSON(&in); err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	out, err := s.Interactor.Save(&in)
	doResponse(c, out, err)
}
