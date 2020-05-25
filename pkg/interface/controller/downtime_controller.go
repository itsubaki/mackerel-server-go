package controller

import (
	"net/http"

	"github.com/itsubaki/mackerel-api/pkg/domain"
	"github.com/itsubaki/mackerel-api/pkg/interface/database"
	"github.com/itsubaki/mackerel-api/pkg/usecase"
)

type DowntimeController struct {
	Interactor *usecase.DowntimeInteractor
}

func NewDowntimeController(handler database.SQLHandler) *DowntimeController {
	return &DowntimeController{
		Interactor: &usecase.DowntimeInteractor{
			DowntimeRepository: database.NewDowntimeRepository(handler),
		},
	}
}

func (s *DowntimeController) List(c Context) {
	out, err := s.Interactor.List(
		c.GetString("org_id"),
	)

	doResponse(c, out, err)
}

func (s *DowntimeController) Save(c Context) {
	var in domain.Downtime
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

func (s *DowntimeController) Update(c Context) {
	var in domain.Downtime
	if err := c.BindJSON(&in); err != nil {
		c.Status(http.StatusBadRequest)
		return
	}
	in.ID = c.Param("downtimeId")

	out, err := s.Interactor.Update(
		c.GetString("org_id"),
		&in,
	)

	doResponse(c, out, err)
}

func (s DowntimeController) Downtime(c Context) {
	out, err := s.Interactor.Downtime(
		c.GetString("org_id"),
		c.Param("downtimeId"),
	)

	doResponse(c, out, err)
}

func (s DowntimeController) Delete(c Context) {
	out, err := s.Interactor.Delete(
		c.GetString("org_id"),
		c.Param("downtimeId"),
	)

	doResponse(c, out, err)
}
