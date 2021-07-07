package controller

import (
	"github.com/itsubaki/mackerel-server-go/pkg/interface/database"
	"github.com/itsubaki/mackerel-server-go/pkg/usecase"
)

type CheckMonitorController struct {
	Interactor *usecase.CheckMonitorInteractor
}

func NewCheckMonitorController(handler database.SQLHandler) *CheckMonitorController {
	return &CheckMonitorController{
		Interactor: &usecase.CheckMonitorInteractor{
			MonitorRepository:    database.NewMonitorRepository(handler),
			HostRepository:       database.NewHostRepository(handler),
			HostMetricRepository: database.NewHostMetricRepository(handler),
			AlertRepository:      database.NewAlertRepository(handler),
		},
	}
}

func (cntr *CheckMonitorController) HostMetric(c Context) {
	out, err := cntr.Interactor.HostMetric(
		c.GetString("org_id"),
	)

	doResponse(c, out, err)
}

func (cntr *CheckMonitorController) Connectivity(c Context) {
	out, err := cntr.Interactor.Connectivity(
		c.GetString("org_id"),
	)

	doResponse(c, out, err)
}

func (cntr *CheckMonitorController) ServiceMetric(c Context) {
	out, err := cntr.Interactor.ServiceMetric(
		c.GetString("org_id"),
	)

	doResponse(c, out, err)
}

func (cntr *CheckMonitorController) External(c Context) {
	out, err := cntr.Interactor.External(
		c.GetString("org_id"),
	)

	doResponse(c, out, err)
}

func (cntr *CheckMonitorController) Expression(c Context) {
	out, err := cntr.Interactor.Expression(
		c.GetString("org_id"),
	)

	doResponse(c, out, err)
}
