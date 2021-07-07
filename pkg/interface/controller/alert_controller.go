package controller

import (
	"net/http"
	"strconv"

	"github.com/itsubaki/mackerel-server-go/pkg/domain"
	"github.com/itsubaki/mackerel-server-go/pkg/interface/database"
	"github.com/itsubaki/mackerel-server-go/pkg/usecase"
)

type AlertController struct {
	Interactor *usecase.AlertInteractor
}

func NewAlertController(handler database.SQLHandler) *AlertController {
	return &AlertController{
		Interactor: &usecase.AlertInteractor{
			AlertRepository: database.NewAlertRepository(handler),
		},
	}
}

func (cntr *AlertController) List(c Context) {
	withClosed, err := strconv.ParseBool(c.DefaultQuery("withClosed", "false"))
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	limit, err := strconv.Atoi(c.DefaultQuery("limit", "100"))
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	out, err := cntr.Interactor.List(
		c.GetString("org_id"),
		withClosed,
		c.Query("nextId"),
		limit,
	)

	doResponse(c, out, err)
}

func (cntr *AlertController) Close(c Context) {
	var in domain.Reason
	if err := c.BindJSON(&in); err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	out, err := cntr.Interactor.Close(
		c.GetString("org_id"),
		c.Param("alertId"),
		in.Reason,
	)

	doResponse(c, out, err)
}
