package controllers

import (
	"net/http"

	"github.com/itsubaki/mackerel-api/pkg/domain"
	"github.com/itsubaki/mackerel-api/pkg/interfaces/database"
	"github.com/itsubaki/mackerel-api/pkg/usecase"
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

func (s *MonitorController) List(c Context) {
	out, err := s.Interactor.List(
		c.GetString("org_id"),
	)

	doResponse(c, out, err)
}

func (s *MonitorController) Save(c Context) {
	var in domain.Monitoring
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

func (s *MonitorController) Update(c Context) {
	var in domain.Monitoring
	if err := c.BindJSON(&in); err != nil {
		c.Status(http.StatusBadRequest)
		return
	}
	in.ID = c.Param("monitorId")

	out, err := s.Interactor.Update(
		c.GetString("org_id"),
		&in,
	)

	doResponse(c, out, err)
}

func (s *MonitorController) Monitor(c Context) {
	out, err := s.Interactor.Monitor(
		c.GetString("org_id"),
		c.Param("monitorId"),
	)

	doResponse(c, out, err)
}

func (s *MonitorController) Delete(c Context) {
	out, err := s.Interactor.Delete(
		c.GetString("org_id"),
		c.Param("monitorId"),
	)

	doResponse(c, out, err)
}
