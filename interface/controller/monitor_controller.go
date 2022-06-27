package controller

import (
	"net/http"

	"github.com/itsubaki/mackerel-server-go/domain"
	"github.com/itsubaki/mackerel-server-go/interface/database"
	"github.com/itsubaki/mackerel-server-go/usecase"
)

type MonitorController struct {
	Interactor *usecase.MonitorInteractor
}

func NewMonitorController(handler database.SQLHandler) *MonitorController {
	return &MonitorController{
		Interactor: &usecase.MonitorInteractor{
			MonitorRepository: database.NewMonitorRepository(handler),
		},
	}
}

func (cntr *MonitorController) List(c Context) {
	out, err := cntr.Interactor.List(
		c.GetString("org_id"),
	)

	doResponse(c, out, err)
}

func (cntr *MonitorController) Save(c Context) {
	var in domain.Monitoring
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

func (cntr *MonitorController) Update(c Context) {
	var in domain.Monitoring
	if err := c.BindJSON(&in); err != nil {
		c.Status(http.StatusBadRequest)
		return
	}
	in.ID = c.Param("monitorId")

	out, err := cntr.Interactor.Update(
		c.GetString("org_id"),
		&in,
	)

	doResponse(c, out, err)
}

func (cntr *MonitorController) Monitor(c Context) {
	out, err := cntr.Interactor.Monitor(
		c.GetString("org_id"),
		c.Param("monitorId"),
	)

	doResponse(c, out, err)
}

func (cntr *MonitorController) Delete(c Context) {
	out, err := cntr.Interactor.Delete(
		c.GetString("org_id"),
		c.Param("monitorId"),
	)

	doResponse(c, out, err)
}
