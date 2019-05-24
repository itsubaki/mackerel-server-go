package controllers

import (
	"net/http"
	"strconv"

	"github.com/itsubaki/mackerel-api/pkg/domain"
	"github.com/itsubaki/mackerel-api/pkg/interfaces/database"
	"github.com/itsubaki/mackerel-api/pkg/usecase"
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

func (s *AlertController) List(c Context) {
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

	out, err := s.Interactor.List(
		c.GetString("org"),
		withClosed,
		c.Query("nextId"),
		limit,
	)

	doResponse(c, out, err)
}

func (s *AlertController) Close(c Context) {
	var in domain.Reason
	if err := c.BindJSON(&in); err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	out, err := s.Interactor.Close(
		c.GetString("org"),
		c.Param("alertId"),
		in.Reason,
	)

	doResponse(c, out, err)
}
