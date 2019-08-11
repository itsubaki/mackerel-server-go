package controller

import (
	"github.com/itsubaki/mackerel-api/pkg/interface/database"
	"github.com/itsubaki/mackerel-api/pkg/usecase"
)

type CheckMonitorController struct {
	Interactor *usecase.CheckMonitorInteractor
}

func NewCheckMonitorController(handler database.SQLHandler) *CheckMonitorController {
	return &CheckMonitorController{
		Interactor: &usecase.CheckMonitorInteractor{
			MonitorRepository: database.NewMonitorRepository(handler),
			HostRepository:    database.NewHostRepository(handler),
			AlertRepository:   database.NewAlertRepository(handler),
		},
	}
}

func (s *CheckMonitorController) HostMetric(c Context) {
	out, err := s.Interactor.HostMetric(
		c.GetString("org_id"),
	)

	doResponse(c, out, err)
}

func (s *CheckMonitorController) Connectivity(c Context) {
	out, err := s.Interactor.Connectivity(
		c.GetString("org_id"),
	)

	doResponse(c, out, err)
}

func (s *CheckMonitorController) ServiceMetric(c Context) {
	out, err := s.Interactor.ServiceMetric(
		c.GetString("org_id"),
	)

	doResponse(c, out, err)
}

func (s *CheckMonitorController) External(c Context) {
	out, err := s.Interactor.External(
		c.GetString("org_id"),
	)

	doResponse(c, out, err)
}

func (s *CheckMonitorController) Expression(c Context) {
	out, err := s.Interactor.Expression(
		c.GetString("org_id"),
	)

	doResponse(c, out, err)
}
