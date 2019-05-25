package controllers

import (
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
