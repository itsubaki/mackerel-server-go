package controller

import (
	"net/http"

	"github.com/itsubaki/mackerel-server-go/pkg/domain"
	"github.com/itsubaki/mackerel-server-go/pkg/interface/database"
	"github.com/itsubaki/mackerel-server-go/pkg/usecase"
)

type DashboardController struct {
	Interactor *usecase.DashboardInteractor
}

func NewDashboardController(handler database.SQLHandler) *DashboardController {
	return &DashboardController{
		Interactor: &usecase.DashboardInteractor{
			DashboardRepository: database.NewDashboardRepository(handler),
		},
	}
}

func (s *DashboardController) List(c Context) {
	out, err := s.Interactor.List(
		c.GetString("org_id"),
	)

	doResponse(c, out, err)
}

func (s *DashboardController) Save(c Context) {
	var in domain.Dashboard
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

func (s *DashboardController) Update(c Context) {
	var in domain.Dashboard
	if err := c.BindJSON(&in); err != nil {
		c.Status(http.StatusBadRequest)
		return
	}
	in.ID = c.Param("dashboardId")

	out, err := s.Interactor.Update(
		c.GetString("org_id"),
		&in,
	)

	doResponse(c, out, err)
}

func (s *DashboardController) Dashboard(c Context) {
	out, err := s.Interactor.Dashboard(
		c.GetString("org_id"),
		c.Param("dashboardId"),
	)

	doResponse(c, out, err)
}

func (s *DashboardController) Delete(c Context) {
	out, err := s.Interactor.Delete(
		c.GetString("org_id"),
		c.Param("dashboardId"),
	)

	doResponse(c, out, err)
}
