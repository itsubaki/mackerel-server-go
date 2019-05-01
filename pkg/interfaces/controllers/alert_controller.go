package controllers

import (
	"net/http"

	"github.com/itsubaki/mackerel-api/pkg/domain"
	"github.com/itsubaki/mackerel-api/pkg/interfaces/database"
	"github.com/itsubaki/mackerel-api/pkg/usecase"
)

type AlertController struct {
	Interactor *usecase.AlertInteractor
}

func NewAlertController(sqlHandler database.SQLHandler) *AlertController {
	return &AlertController{}
}

func (s *AlertController) List(c Context) {
	out, err := s.Interactor.List(
		false,
		c.Query("nextId"),
		100,
	)

	doResponse(c, out, err)
}

func (s *AlertController) Close(c Context) {
	var in domain.Reason
	if err := c.BindJSON(in); err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	out, err := s.Interactor.Close(
		c.Param("alertId"),
		in.Reason,
	)
	doResponse(c, out, err)
}
