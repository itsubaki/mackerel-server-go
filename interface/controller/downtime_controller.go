package controller

import (
	"net/http"

	"github.com/itsubaki/mackerel-server-go/domain"
	"github.com/itsubaki/mackerel-server-go/interface/database"
	"github.com/itsubaki/mackerel-server-go/usecase"
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

func (cntr *DowntimeController) List(c Context) {
	out, err := cntr.Interactor.List(
		c.GetString("org_id"),
	)

	doResponse(c, out, err)
}

func (cntr *DowntimeController) Save(c Context) {
	var in domain.Downtime
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

func (cntr *DowntimeController) Update(c Context) {
	var in domain.Downtime
	if err := c.BindJSON(&in); err != nil {
		c.Status(http.StatusBadRequest)
		return
	}
	in.ID = c.Param("downtimeId")

	out, err := cntr.Interactor.Update(
		c.GetString("org_id"),
		&in,
	)

	doResponse(c, out, err)
}

func (cntr *DowntimeController) Downtime(c Context) {
	out, err := cntr.Interactor.Downtime(
		c.GetString("org_id"),
		c.Param("downtimeId"),
	)

	doResponse(c, out, err)
}

func (cntr *DowntimeController) Delete(c Context) {
	out, err := cntr.Interactor.Delete(
		c.GetString("org_id"),
		c.Param("downtimeId"),
	)

	doResponse(c, out, err)
}
