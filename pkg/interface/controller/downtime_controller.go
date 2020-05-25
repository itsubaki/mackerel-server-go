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

// TODO
//Error
//STATUS CODE	DESCRIPTION
//400	when the input is invalid
//400	when the name or memo is too long
//400	when the downtime duration is invalid
//400	when the recurrence configuration is invalid
//400	when the target service/role/monitor setting of the scope is redundant or does not exist
//400	when the service/role/monitor setting of the scope does not exist
//403	when the API key doesn't have the required permissions / when accessing from outside the permitted IP address range
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

// TODO
//Error
//STATUS CODE	DESCRIPTION
//400	when the input is invalid
//400	when the name or memo is too long
//400	when the downtime duration is invalid
//400	when the recurrence configuration is invalid
//400	when the target service/role/monitor setting of the scope is redundant or does not exist
//400	when the service/role/monitor setting of the scope does not exist
//403	when the API key doesn't have the required permissions / when accessing from outside the permitted IP address range
//404	when the downtime corresponding to the designated ID can't be found
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

// TODO
//Error
//STATUS CODE	DESCRIPTION
//403	when the API key doesn't have the required permissions / when accessing from outside the permitted IP address range
//404	when the downtime corresponding to the designated ID can't be found
func (s DowntimeController) Delete(c Context) {
	out, err := s.Interactor.Delete(
		c.GetString("org_id"),
		c.Param("downtimeId"),
	)

	doResponse(c, out, err)
}
